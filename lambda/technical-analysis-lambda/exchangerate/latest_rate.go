package exchangerate

import (
	"context"
	"errors"
	"log"
	"os"
	"time"

	"technical-analysis-lambda/finance"

	finnhub "github.com/Finnhub-Stock-API/finnhub-go"
)

func CreateNewSymbolRate(symbols *[4]string, startSeconds int64, endSeconds int64, period string, headline *string) (*[]finance.FinancialDataItem, *[]finance.NewsItem, error) {
	finnhubClient := finnhub.NewAPIClient(finnhub.NewConfiguration()).DefaultApi
	auth := context.WithValue(context.Background(), finnhub.ContextAPIKey, finnhub.APIKey{
		Key: os.Getenv("PROJECT_API_KEY"),
	})
	symbolRateItems := []finance.FinancialDataItem{}
	currentTime := time.Now().Unix()
	for _, symbol := range *symbols {
		symbolFinnhub := symbol[:3] + "_" + symbol[3:]
		forexCandles, _, err := finnhubClient.ForexCandles(
			auth, "OANDA:"+symbolFinnhub, period, time.Now().Unix()-startSeconds,
			time.Now().Unix()-endSeconds,
		)
		if err != nil {
			log.Println("Issue connecting to forex client")
			log.Println(err)
			return nil, nil, errors.New("forex client connection error")
		}
		if forexCandles.S != "ok" {
			log.Println("No data received")
			return nil, nil, nil
		}

		rateLength := len(forexCandles.O)
		for i := 0; i < rateLength; i++ {
			intTimestamp := int64(forexCandles.T[i])
			if (intTimestamp > currentTime-startSeconds) && (intTimestamp < currentTime-endSeconds) {
				if i != rateLength-1 && intTimestamp == int64(forexCandles.T[i+1]) {
					continue
				}
				open := forexCandles.O[i]
				high := forexCandles.H[i]
				low := forexCandles.L[i]
				close := forexCandles.C[i]
				volume := forexCandles.V[i]
				timestamp := time.Unix(intTimestamp, 0).UTC()

				newSymbolData := finance.SymbolData{Open: open, High: high, Low: low, Close: close, Volume: volume}
				formattedDate := timestamp.Format("2006-01-02")
				formattedTimestamp := timestamp.Format("15:04:05")
				if searchTimeIndex(&symbolRateItems, formattedDate, formattedTimestamp) == -1 {
					switch symbol {
					case "EURUSD":
						newSymbolRateItem := finance.FinancialDataItem{Date: formattedDate, Timestamp: formattedTimestamp, EURUSD: newSymbolData}
						symbolRateItems = append(symbolRateItems, newSymbolRateItem)
					case "GBPUSD":
						newSymbolRateItem := finance.FinancialDataItem{Date: formattedDate, Timestamp: formattedTimestamp, GBPUSD: newSymbolData}
						symbolRateItems = append(symbolRateItems, newSymbolRateItem)
					case "USDJPY":
						newSymbolRateItem := finance.FinancialDataItem{Date: formattedDate, Timestamp: formattedTimestamp, USDJPY: newSymbolData}
						symbolRateItems = append(symbolRateItems, newSymbolRateItem)
					case "AUDCAD":
						newSymbolRateItem := finance.FinancialDataItem{Date: formattedDate, Timestamp: formattedTimestamp, AUDCAD: newSymbolData}
						symbolRateItems = append(symbolRateItems, newSymbolRateItem)
					}
				} else {
					index := searchTimeIndex(&symbolRateItems, formattedDate, formattedTimestamp)
					switch symbol {
					case "EURUSD":
						symbolRateItems[index].EURUSD = newSymbolData
					case "GBPUSD":
						symbolRateItems[index].GBPUSD = newSymbolData
					case "USDJPY":
						symbolRateItems[index].USDJPY = newSymbolData
					case "AUDCAD":
						symbolRateItems[index].AUDCAD = newSymbolData
					}
				}
			}
		}
	}

	latestNews, err := getMarketNews(finnhubClient, &auth, headline)
	if err != nil {
		log.Println(err)
		return &symbolRateItems, nil, nil
	}

	return &symbolRateItems, latestNews, nil
}

func searchTimeIndex(symbolData *[]finance.FinancialDataItem, date string, timestamp string) int {
	for index, data := range *symbolData {
		if data.Date == date && data.Timestamp == timestamp {
			return index
		}
	}
	return -1
}
