package exchangerate

import (
	"context"
	"log"
	"os"
	"time"

	finnhub "github.com/Finnhub-Stock-API/finnhub-go"
)

// Toggle to switch between getting real data and generating fake data
// When false there is no risk of being rate limited by finnhub api
const isRealData bool = false

func rateReduce(currencyData *[]ExchangeRate, prevRate *ExchangeRate) *[]ExchangeRate {
	filteredCurrencyData := []ExchangeRate{}

	for i := 0; i < len(*currencyData); i++ {
		currentData := (*currencyData)[i]
		if currentData.Timestamp.After(prevRate.Timestamp) {
			prevRate.Timestamp = currentData.Timestamp
			prevRate.Open = currentData.Open
			prevRate.High = currentData.High
			prevRate.Low = currentData.Low
			prevRate.Close = currentData.Close
			prevRate.Volume = currentData.Volume

			filteredCurrencyData = append(filteredCurrencyData, *prevRate)
		}
	}
	log.Println(filteredCurrencyData, "Filtered")

	return &filteredCurrencyData
}

func GetLatestRate(symbol string, prevRate *ExchangeRate, startSeconds int64, endSeconds int64, period string) map[string][]ExchangeRate {
	if isRealData {
		symbolFinnhub := symbol[:3] + "_" + symbol[3:]
		finnhubClient := finnhub.NewAPIClient(finnhub.NewConfiguration()).DefaultApi
		auth := context.WithValue(context.Background(), finnhub.ContextAPIKey, finnhub.APIKey{
			Key: os.Getenv("PROJECT_API_KEY"),
		})
		forexCandles, _, err := finnhubClient.ForexCandles(
			auth, "OANDA:"+symbolFinnhub, period, time.Now().Unix()-startSeconds,
			time.Now().Unix()-endSeconds,
		)
		if err != nil {
			log.Println("Issue connecting to forex client")
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
		filteredRate := rateReduce(&pastRate, prevRate)
		if filteredRate == nil {
			log.Println("No new data received")
			return nil
		}
		log.Printf("%+v New Data\n", pastRate)
		symbolData := map[string][]ExchangeRate{
			symbol: *filteredRate,
		}

		return symbolData
	} else {
		return GenerateFakeData(prevRate, symbol, startSeconds, endSeconds)
	}
}
