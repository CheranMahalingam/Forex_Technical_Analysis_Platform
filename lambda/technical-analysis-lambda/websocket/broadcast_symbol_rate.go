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

	apigw := apigatewaymanagementapi.New(sess, aws.NewConfig().WithEndpoint(os.Getenv("API_GATEWAY_URI")))

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
		symbolField := getSymbolRateField(symbol, &rate)
		if checkIsRateValid(symbolField, symbol) {
			newData := CallbackMessageData{
				Timestamp: rate.Date + " " + rate.Timestamp,
				Open:      symbolField.Open,
				High:      symbolField.High,
				Low:       symbolField.Low,
				Close:     symbolField.Close,
				Volume:    symbolField.Volume,
			}
			newCallbackMessageData = append(newCallbackMessageData, newData)
		}
	}
	return &newCallbackMessageData, nil
}

func checkIsRateValid(symbolField *dynamosymbol.SymbolData, symbol string) bool {
	if symbolField.Open == 0 ||
		symbolField.High == 0 ||
		symbolField.Low == 0 ||
		symbolField.Close == 0 {
		return false
	}
	return true
}

func getSymbolRateField(symbol string, symbolRate *dynamosymbol.SymbolRateItem) *dynamosymbol.SymbolData {
	switch symbol {
	case "EURUSD":
		return &symbolRate.EURUSD
	case "GBPUSD":
		return &symbolRate.GBPUSD
	default:
		return nil
	}
}
