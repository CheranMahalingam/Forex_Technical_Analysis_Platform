package controllers

import (
	"context"
	"log"
	"math/rand"
	"os"
	"time"

	finnhub "github.com/Finnhub-Stock-API/finnhub-go"
)

// Store previous open price when generating fake data
var testOpen float32 = 1.0

// Toggle to switch between getting real data and generating fake data
// When false there is no risk of being rate limited by finnhub api
const isRealData bool = false

type ExchangeRate struct {
	Timestamp time.Time `json:"timestamp"`
	Open      float32   `json:"open"`
	High      float32   `json:"high"`
	Low       float32   `json:"low"`
	Close     float32   `json:"close"`
	Volume    float32   `json:"volume"`
}

// Used to generate data for weekends to simulate live trading
// Useful to avoid being rate limited by finnhub api
func generateFakeData(seconds int64) *[]ExchangeRate {
	var newExchangeRateData []ExchangeRate
	dataLength := int(seconds / 60)
	var timestamp time.Time
	var open float32
	var close float32
	var high float32
	var low float32
	var volume float32
	for i := 0; i < dataLength; i++ {
		if i == 0 {
			open = testOpen //(rand.Float32()*2-1)*0.001 + 1
		} else {
			open = (rand.Float32()*2-1)*0.001 + close
			testOpen = open
		}
		close = (rand.Float32()*2-1)*0.001 + open
		if open > close {
			high = rand.Float32()*0.001 + open
			low = close - rand.Float32()*0.001
		} else {
			high = rand.Float32()*0.001 + close
			low = open - rand.Float32()*0.001
		}
		timestamp = time.Unix(time.Now().Unix()-seconds+int64((i+1)*60), 0).UTC()
		volume = rand.Float32() * 100
		newExchangeRateData = append(newExchangeRateData, ExchangeRate{Timestamp: timestamp, Open: open, High: high, Low: low, Close: close, Volume: volume})
	}
	return &newExchangeRateData
}

func GetLatestRate(symbol string, seconds int64, period string) *[]ExchangeRate {
	if isRealData {
		finnhubClient := finnhub.NewAPIClient(finnhub.NewConfiguration()).DefaultApi
		auth := context.WithValue(context.Background(), finnhub.ContextAPIKey, finnhub.APIKey{
			Key: os.Getenv("PROJECT_API_KEY"),
		})
		forexCandles, _, err := finnhubClient.ForexCandles(auth, "OANDA:"+symbol, period, time.Now().Unix()-seconds, time.Now().Unix())
		if err != nil {
			log.Println("Issue connecting to forex client")
			log.Println(err)
			return nil
		}
		log.Printf("%+v\n", forexCandles)

		if forexCandles.S != "ok" {
			log.Println("No data received")
			return nil
		}

		rateLength := len(forexCandles.O)
		pastRate := make([]ExchangeRate, 0, rateLength)
		for i := 0; i < rateLength; i++ {
			open := forexCandles.O[i]
			high := forexCandles.H[i]
			low := forexCandles.L[i]
			close := forexCandles.C[i]
			volume := forexCandles.V[i]
			timestamp := time.Unix(int64(forexCandles.T[i]), 0).UTC()
			pastRate = append(pastRate, ExchangeRate{Timestamp: timestamp, Open: open, High: high, Low: low, Close: close, Volume: volume})
		}
		log.Printf("%+v\n", pastRate)

		return &pastRate
	} else {
		return generateFakeData(seconds)
	}
}
