package exchangerate

import (
	"context"
	"log"
	"os"
	"time"

	finnhub "github.com/Finnhub-Stock-API/finnhub-go"
)

type SymbolData struct {
	Open   float32
	High   float32
	Low    float32
	Close  float32
	Volume float32
}

type SymbolRateItem struct {
	Timestamp time.Time
	EurUsd    SymbolData
	GbpUsd    SymbolData
}

// Toggle to switch between getting real data and generating fake data
// When false there is no risk of being rate limited by finnhub api
const isRealData bool = true

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
		currentTime := time.Now().Unix()
		log.Println("START", currentTime-startSeconds)
		log.Println("END", currentTime-endSeconds)
		forexCandles, _, err := finnhubClient.ForexCandles(
			auth, "OANDA:"+symbolFinnhub, period, currentTime-startSeconds,
			currentTime-endSeconds,
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

		log.Println(forexCandles)

		rateLength := len(forexCandles.O)
		pastRate := make([]ExchangeRate, 0, rateLength)
		for i := 0; i < rateLength; i++ {
			intTimestamp := int64(forexCandles.T[i])
			if (intTimestamp > currentTime-startSeconds) && (intTimestamp < currentTime-endSeconds) {
				if i != rateLength-1 && intTimestamp == int64(forexCandles.T[i+1]) {
					log.Println(intTimestamp)
					continue
				}
				open := forexCandles.O[i]
				high := forexCandles.H[i]
				low := forexCandles.L[i]
				close := forexCandles.C[i]
				volume := forexCandles.V[i]
				timestamp := time.Unix(intTimestamp, 0).UTC()
				pastRate = append(pastRate, ExchangeRate{Timestamp: timestamp, Open: open, High: high, Low: low, Close: close, Volume: volume})
			}
			log.Println(intTimestamp)
		}
		// filteredRate := rateReduce(&pastRate, prevRate)
		// if filteredRate == nil {
		// 	log.Println("No new data received")
		// 	return nil
		// }
		log.Printf("%+v New Data\n", pastRate)
		symbolData := map[string][]ExchangeRate{
			symbol: pastRate,
		}
		if len(symbolData) == 0 {
			return nil
		}
		return symbolData
	} else {
		return GenerateFakeData(prevRate, symbol, startSeconds, endSeconds)
	}
}

func NewCurrencyRate(symbols *[2]string, startSeconds int64, endSeconds int64, period string) map[string][]ExchangeRate {
	finnhubClient := finnhub.NewAPIClient(finnhub.NewConfiguration()).DefaultApi
	auth := context.WithValue(context.Background(), finnhub.ContextAPIKey, finnhub.APIKey{
		Key: os.Getenv("PROJECT_API_KEY"),
	})
	symbolData := map[string][]ExchangeRate{}
	for _, symbol := range *symbols {
		symbolFinnhub := symbol[:3] + "_" + symbol[3:]
		currentTime := time.Now().Unix()
		log.Println("START", currentTime-startSeconds)
		log.Println("END", currentTime-endSeconds)
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
			intTimestamp := int64(forexCandles.T[i])
			if (intTimestamp > currentTime-startSeconds) && (intTimestamp < currentTime-endSeconds) {
				if i != rateLength-1 && intTimestamp == int64(forexCandles.T[i+1]) {
					log.Println(intTimestamp)
					continue
				}
				open := forexCandles.O[i]
				high := forexCandles.H[i]
				low := forexCandles.L[i]
				close := forexCandles.C[i]
				volume := forexCandles.V[i]
				timestamp := time.Unix(intTimestamp, 0).UTC()
				pastRate = append(pastRate, ExchangeRate{Timestamp: timestamp, Open: open, High: high, Low: low, Close: close, Volume: volume})
			}
		}
		symbolData[symbol] = pastRate
	}
	return symbolData
}

func CreateNewSymbolRate(symbols *[2]string, startSeconds int64, endSeconds int64, period string) *[]SymbolRateItem {
	finnhubClient := finnhub.NewAPIClient(finnhub.NewConfiguration()).DefaultApi
	auth := context.WithValue(context.Background(), finnhub.ContextAPIKey, finnhub.APIKey{
		Key: os.Getenv("PROJECT_API_KEY"),
	})
	symbolRateItems := []SymbolRateItem{}
	for _, symbol := range *symbols {
		symbolFinnhub := symbol[:3] + "_" + symbol[3:]
		currentTime := time.Now().Unix()
		log.Println("START", currentTime-startSeconds)
		log.Println("END", currentTime-endSeconds)
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
		// pastRate := make([]ExchangeRate, 0, rateLength)
		for i := 0; i < rateLength; i++ {
			intTimestamp := int64(forexCandles.T[i])
			if (intTimestamp > currentTime-startSeconds) && (intTimestamp < currentTime-endSeconds) {
				if i != rateLength-1 && intTimestamp == int64(forexCandles.T[i+1]) {
					log.Println(intTimestamp)
					continue
				}
				open := forexCandles.O[i]
				high := forexCandles.H[i]
				low := forexCandles.L[i]
				close := forexCandles.C[i]
				volume := forexCandles.V[i]
				timestamp := time.Unix(intTimestamp, 0).UTC()
				// pastRate = append(pastRate, ExchangeRate{Timestamp: timestamp, Open: open, High: high, Low: low, Close: close, Volume: volume})

				newSymbolData := SymbolData{Open: open, High: high, Low: low, Close: close, Volume: volume}
				if searchTimestampIndex(&symbolRateItems, timestamp) == -1 {
					switch symbol {
					case "EURUSD":
						newSymbolRateItem := SymbolRateItem{Timestamp: timestamp, EurUsd: newSymbolData}
						symbolRateItems = append(symbolRateItems, newSymbolRateItem)
					case "GBPUSD":
						newSymbolRateItem := SymbolRateItem{Timestamp: timestamp, GbpUsd: newSymbolData}
						symbolRateItems = append(symbolRateItems, newSymbolRateItem)
					}
				} else {
					index := searchTimestampIndex(&symbolRateItems, timestamp)
					switch symbol {
					case "EURUSD":
						symbolRateItems[index].EurUsd = newSymbolData
					case "GBPUSD":
						symbolRateItems[index].GbpUsd = newSymbolData
					}
				}
			}
		}
	}

	return &symbolRateItems
}

func searchTimestampIndex(symbolData *[]SymbolRateItem, timestamp time.Time) int {
	for index, data := range *symbolData {
		if data.Timestamp == timestamp {
			return index
		}
	}
	return -1
}
