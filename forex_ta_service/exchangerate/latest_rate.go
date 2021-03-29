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
const isRealData bool = true

func GetLatestRate(symbol string, seconds int64, period string) *map[string][]ExchangeRate {
	if isRealData {
		symbolFinnhub := symbol[:3] + "_" + symbol[3:]
		finnhubClient := finnhub.NewAPIClient(finnhub.NewConfiguration()).DefaultApi
		auth := context.WithValue(context.Background(), finnhub.ContextAPIKey, finnhub.APIKey{
			Key: os.Getenv("PROJECT_API_KEY"),
		})
		forexCandles, _, err := finnhubClient.ForexCandles(auth, "OANDA:"+symbolFinnhub, period, time.Now().Unix()-seconds, time.Now().Unix())
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
		symbolData := map[string][]ExchangeRate{
			symbol: pastRate,
		}

		return &symbolData
	} else {
		return GenerateFakeData(symbol, seconds)
	}
}
