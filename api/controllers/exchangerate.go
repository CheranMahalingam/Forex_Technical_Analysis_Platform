package controllers

import (
	"context"
	"log"
	"os"
	"time"

	finnhub "github.com/Finnhub-Stock-API/finnhub-go"
)

type ExchangeRate struct {
	Timestamp time.Time `json:"timestamp"`
	Open      float32   `json:"open"`
	High      float32   `json:"high"`
	Low       float32   `json:"low"`
	Close     float32   `json:"close"`
	Volume    float32   `json:"volume"`
}

func GetLatestRate(symbol string, seconds int64, period string) *[]ExchangeRate {
	finnhubClient := finnhub.NewAPIClient(finnhub.NewConfiguration()).DefaultApi
	auth := context.WithValue(context.Background(), finnhub.ContextAPIKey, finnhub.APIKey{
		Key: os.Getenv("PROJECT_API_KEY"),
	})
	forexCandles, _, err := finnhubClient.ForexCandles(auth, "OANDA:"+symbol, period, time.Now().Unix()-seconds, time.Now().Unix())
	if err != nil {
		log.Println(err)
		return nil
	}

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
}
