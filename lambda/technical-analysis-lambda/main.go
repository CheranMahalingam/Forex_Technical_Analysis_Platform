package main

import (
	"encoding/json"
	"log"
	"technical-analysis-lambda/exchangerate"
	"technical-analysis-lambda/finance"
	"technical-analysis-lambda/websocket"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var currencyPairs = [4]string{"EURUSD", "GBPUSD", "USDJPY", "AUDCAD"}

func handler(request events.CloudWatchEvent) (events.CloudWatchEvent, error) {
	prevHeadline, err := finance.GetNewsHeadline()
	if err != nil {
		log.Println(err)
		return events.CloudWatchEvent{}, err
	}
	newRate, newNews, err := exchangerate.CreateNewSymbolRate(&currencyPairs, 120, 0, "1", prevHeadline)
	if err != nil {
		log.Println(err)
		return events.CloudWatchEvent{}, err
	}
	if newRate == nil {
		return events.CloudWatchEvent{}, nil
	}
	finance.SendRateToDB(newRate, newNews)
	byteNewRate, err := json.Marshal(newRate)
	if err != nil {
		log.Println("Could not marshal struct", err)
	}

	for i := 0; i < len(currencyPairs); i++ {
		connectionList, err := websocket.ReadSymbolConnections(currencyPairs[i], &byteNewRate)
		if err != nil {
			log.Println("read symbol rate error:", err)
		}
		if err = websocket.BroadcastSymbolRate(connectionList, newRate, currencyPairs[i]); err != nil {
			log.Println("Websocket Exchange Rate Broadcasting Error", err)
		}
	}

	if websocket.ValidateNewMarketNews(newNews, *prevHeadline) {
		if err == nil {
			newsConnectionList, err := websocket.ReadNewsConnections()
			if err != nil {
				log.Println("read market news error:", err)
			}
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
