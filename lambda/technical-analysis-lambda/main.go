package main

import (
	"encoding/json"
	"log"
	"technical-analysis-lambda/dynamosymbol"
	"technical-analysis-lambda/exchangerate"
	"technical-analysis-lambda/websocket"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var currencyPairs = [2]string{"EURUSD", "GBPUSD"}

func handler(request events.CloudWatchEvent) (events.CloudWatchEvent, error) {
	newRate, err := exchangerate.CreateNewSymbolRate(&currencyPairs, 120, 0, "1")
	if err != nil {
		log.Println(err)
		return events.CloudWatchEvent{}, err
	}
	if newRate == nil {
		return events.CloudWatchEvent{}, nil
	}
	log.Println(newRate)
	dynamosymbol.SendRateToDB(newRate)
	byteNewRate, err := json.Marshal(newRate)
	if err != nil {
		log.Println("Could not marshal struct", err)
	}
	for i := 0; i < len(currencyPairs); i++ {
		connectionList, err := websocket.ReadConnections(currencyPairs[i], &byteNewRate)
		if err != nil {
			log.Println("read symbol rate error:", err)
		}
		if err = websocket.BroadcastSymbolRate(connectionList, newRate, currencyPairs[i]); err != nil {
			log.Println("Websocket broadcasting error", err)
		}
	}
	return events.CloudWatchEvent{}, nil
}

func main() {
	lambda.Start(handler)
}
