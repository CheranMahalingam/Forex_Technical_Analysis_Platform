package exchangerate

import (
	"context"
	"errors"
	"log"
	"os"
	"time"

	"technical-analysis-lambda/dynamosymbol"

	finnhub "github.com/Finnhub-Stock-API/finnhub-go"
)

func CreateNewSymbolRate(symbols *[2]string, startSeconds int64, endSeconds int64, period string) (*[]dynamosymbol.SymbolRateItem, error) {
	finnhubClient := finnhub.NewAPIClient(finnhub.NewConfiguration()).DefaultApi
	auth := context.WithValue(context.Background(), finnhub.ContextAPIKey, finnhub.APIKey{
		Key: os.Getenv("ProjectApiKey"),
	})
	symbolRateItems := []dynamosymbol.SymbolRateItem{}
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
			return nil, errors.New("forex client connection error")
		}
		if forexCandles.S != "ok" {
			log.Println("No data received")
			return nil, nil
		}

		log.Println(forexCandles)

		rateLength := len(forexCandles.O)
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

				newSymbolData := dynamosymbol.SymbolData{Open: open, High: high, Low: low, Close: close, Volume: volume}
				formattedDate := timestamp.Format("2006-01-02")
				formattedTimestamp := timestamp.Format("15:04:05")
				log.Println(formattedDate, formattedTimestamp)
				if searchTimeIndex(&symbolRateItems, formattedDate, formattedTimestamp) == -1 {
					switch symbol {
					case "EURUSD":
						newSymbolRateItem := dynamosymbol.SymbolRateItem{Date: formattedDate, Timestamp: formattedTimestamp, EurUsd: newSymbolData}
						symbolRateItems = append(symbolRateItems, newSymbolRateItem)
					case "GBPUSD":
						newSymbolRateItem := dynamosymbol.SymbolRateItem{Date: formattedDate, Timestamp: formattedTimestamp, GbpUsd: newSymbolData}
						symbolRateItems = append(symbolRateItems, newSymbolRateItem)
					}
				} else {
					index := searchTimeIndex(&symbolRateItems, formattedDate, formattedTimestamp)
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
	return &symbolRateItems, nil
}

func searchTimeIndex(symbolData *[]dynamosymbol.SymbolRateItem, date string, timestamp string) int {
	for index, data := range *symbolData {
		if data.Date == date && data.Timestamp == timestamp {
			return index
		}
	}
	return -1
}
