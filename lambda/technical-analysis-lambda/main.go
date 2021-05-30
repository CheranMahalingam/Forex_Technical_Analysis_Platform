package main

import (
	"log"
	"technical-analysis-lambda/exchangerate"
	"technical-analysis-lambda/finance"
	"technical-analysis-lambda/websocket"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// Currency pairs for which lstm models are currently supported
var currencyPairs = [4]string{"EURUSD", "GBPUSD", "USDJPY", "AUDCAD"}

func handler(request events.CloudWatchEvent) (events.CloudWatchEvent, error) {
	// Gets previous headline from DynamoDB
	prevHeadline, err := finance.GetNewsHeadline()
	if err != nil {
		log.Println(err)
		return events.CloudWatchEvent{}, err
	}
	// Gets latest exchange rates for each symbol and market news
	finnhubClient := exchangerate.NewForexApi()
	newRate, newNews, err := exchangerate.CreateNewSymbolRate(finnhubClient, &currencyPairs, 120, 0, "1", prevHeadline)
	if err != nil {
		log.Println(err)
		return events.CloudWatchEvent{}, err
	}
	if newRate == nil {
		return events.CloudWatchEvent{}, nil
	}
	// Sends market news and exchange rates to DynamoDB
	finance.SendRateToDB(newRate, newNews)

	for i := 0; i < len(currencyPairs); i++ {
		connectionList, err := websocket.ReadConnections(currencyPairs[i])
		if err != nil {
			log.Println("read symbol rate error:", err)
		}
		// Sends new exchange rate data to all users subscribed to the exchange rate data channel
		if err = websocket.BroadcastSymbolRate(connectionList, newRate, currencyPairs[i]); err != nil {
			log.Println("Websocket Exchange Rate Broadcasting Error", err)
		}
	}

	// Must ensure that the headline is different from previous
	if websocket.ValidateNewMarketNews(newNews, *prevHeadline) {
		if err == nil {
			newsConnectionList, err := websocket.ReadConnections("News")
			if err != nil {
				log.Println("read market news error:", err)
			}
			// Sends new market news to users subscribed to the market news channel
			if err = websocket.BroadcastMarketNews(newsConnectionList, newNews); err != nil {
				log.Println("Websocket News Broadcasting Error", err)
			}
		} else {
			log.Println(err)
		}
	}

	return events.CloudWatchEvent{}, nil
}

func main() {
	lambda.Start(handler)
}
