package main

import (
	"log"
	"technical-analysis-lambda/dynamosymbol"
	"technical-analysis-lambda/exchangerate"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var currencyPairs = [2]string{"EURUSD", "GBPUSD"}

func handler(request events.CloudWatchEvent) (events.CloudWatchEvent, error) {
	newRate := exchangerate.CreateNewSymbolRate(&currencyPairs, 120, 0, "1")
	log.Println(newRate)
	dynamosymbol.SendRateToDB(newRate)
	return events.CloudWatchEvent{}, nil
}

func main() {
	lambda.Start(handler)
}
