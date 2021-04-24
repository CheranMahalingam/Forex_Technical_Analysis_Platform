package websocket

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"technical-analysis-lambda/dynamosymbol"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/apigatewaymanagementapi"
)

func BroadcastSymbolRate(connectionList *[]dynamosymbol.Connection, symbolRate *[]dynamosymbol.SymbolRateItem, symbol string) error {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	apigw := apigatewaymanagementapi.New(sess, aws.NewConfig().WithEndpoint(os.Getenv("ApiGatewayUri")))

	for _, conn := range *connectionList {
		log.Println(conn.ConnectionId)
		for _, rate := range *symbolRate {
			newData, err := createCallbackMessage(&rate, symbol)
			if err != nil {
				return err
			}
			if newData == nil {
				break
			}
			byteMessage, err := json.Marshal(*newData)
			if err != nil {
				log.Println("Marshalling error", err)
				return errors.New("callback message marshalling error")
			}
			if _, err = apigw.PostToConnection(&apigatewaymanagementapi.PostToConnectionInput{ConnectionId: &conn.ConnectionId, Data: byteMessage}); err != nil {
				log.Println("Could not send to api", err)
				return errors.New("could not send to apigatway")
			}
		}
	}

	return nil
}

func createCallbackMessage(symbolRate *dynamosymbol.SymbolRateItem, symbol string) (*CallbackMessageData, error) {
	if !checkIsRateValid(symbolRate, symbol) {
		log.Printf("No data for symbol %s", symbol)
		return nil, nil
	}
	switch symbol {
	case "EurUsd":
		newCallbackMessageData := CallbackMessageData{
			Symbol:    symbol,
			Timestamp: symbolRate.Date + symbolRate.Timestamp,
			Open:      symbolRate.EurUsd.Open,
			High:      symbolRate.EurUsd.High,
			Low:       symbolRate.EurUsd.Low,
			Close:     symbolRate.EurUsd.Close,
			Volume:    symbolRate.EurUsd.Volume,
		}
		return &newCallbackMessageData, nil
	case "GbpUsd":
		newCallbackMessageData := CallbackMessageData{
			Symbol:    symbol,
			Timestamp: symbolRate.Date + symbolRate.Timestamp,
			Open:      symbolRate.GbpUsd.Open,
			High:      symbolRate.GbpUsd.High,
			Low:       symbolRate.GbpUsd.Low,
			Close:     symbolRate.GbpUsd.Close,
			Volume:    symbolRate.GbpUsd.Volume,
		}
		return &newCallbackMessageData, nil
	}
	return nil, errors.New("unknown symbol")
}

func checkIsRateValid(symbolRate *dynamosymbol.SymbolRateItem, symbol string) bool {
	switch symbol {
	case "EurUsd":
		if symbolRate.EurUsd.Open == 0 ||
			symbolRate.EurUsd.High == 0 ||
			symbolRate.EurUsd.Low == 0 ||
			symbolRate.EurUsd.Close == 0 {
			return false
		}
	case "GbpUsd":
		if symbolRate.GbpUsd.Open == 0 ||
			symbolRate.GbpUsd.High == 0 ||
			symbolRate.GbpUsd.Low == 0 ||
			symbolRate.GbpUsd.Close == 0 {
			return false
		}
	}
	return true
}
