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

	newData, err := createCallbackMessage(symbolRate, symbol)
	if err != nil {
		return err
	}
	if newData == nil {
		log.Println("No new data")
		return errors.New("no data to broadcast")
	}

	newSymbolData := map[string][]CallbackMessageData{}
	newSymbolData[symbol] = *newData

	byteMessage, err := json.Marshal(newSymbolData)
	if err != nil {
		log.Println("Marshalling error", err)
		return errors.New("callback message marshalling error")
	}

	for _, conn := range *connectionList {
		log.Println(conn.ConnectionId)
		if _, err = apigw.PostToConnection(&apigatewaymanagementapi.PostToConnectionInput{ConnectionId: &conn.ConnectionId, Data: byteMessage}); err != nil {
			log.Println("Could not send to api", err)
			return errors.New("could not send to apigateway")
		}
	}

	return nil
}

func createCallbackMessage(symbolRate *[]dynamosymbol.SymbolRateItem, symbol string) (*[]CallbackMessageData, error) {
	var newCallbackMessageData []CallbackMessageData
	for _, rate := range *symbolRate {
		if !checkIsRateValid(&rate, symbol) {
			log.Printf("No data for symbol %s", symbol)
			continue
		}
		switch symbol {
		case "EURUSD":
			newData := CallbackMessageData{
				Timestamp: rate.Date + " " + rate.Timestamp,
				Open:      rate.EURUSD.Open,
				High:      rate.EURUSD.High,
				Low:       rate.EURUSD.Low,
				Close:     rate.EURUSD.Close,
				Volume:    rate.EURUSD.Volume,
			}
			newCallbackMessageData = append(newCallbackMessageData, newData)
		case "GBPUSD":
			newData := CallbackMessageData{
				Timestamp: rate.Date + " " + rate.Timestamp,
				Open:      rate.GBPUSD.Open,
				High:      rate.GBPUSD.High,
				Low:       rate.GBPUSD.Low,
				Close:     rate.GBPUSD.Close,
				Volume:    rate.GBPUSD.Volume,
			}
			newCallbackMessageData = append(newCallbackMessageData, newData)
		default:
			return nil, errors.New("unknown symbol")
		}
	}
	return &newCallbackMessageData, nil
}

func checkIsRateValid(symbolRate *dynamosymbol.SymbolRateItem, symbol string) bool {
	switch symbol {
	case "EURUSD":
		if symbolRate.EURUSD.Open == 0 ||
			symbolRate.EURUSD.High == 0 ||
			symbolRate.EURUSD.Low == 0 ||
			symbolRate.EURUSD.Close == 0 {
			return false
		}
	case "GBPUSD":
		if symbolRate.GBPUSD.Open == 0 ||
			symbolRate.GBPUSD.High == 0 ||
			symbolRate.GBPUSD.Low == 0 ||
			symbolRate.GBPUSD.Close == 0 {
			return false
		}
	}
	return true
}
