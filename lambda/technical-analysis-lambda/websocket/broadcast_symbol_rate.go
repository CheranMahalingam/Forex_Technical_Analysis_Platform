package websocket

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"technical-analysis-lambda/finance"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/apigatewaymanagementapi"
)

// Sends new ohlc data to all users subscribed to a particular symbol's data channel
// Subscription occurs when on the charts webpage and user selects a currency pair
func BroadcastSymbolRate(connectionList *[]finance.Connection, symbolRate *[]finance.FinancialDataItem, symbol string) error {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	apigw := apigatewaymanagementapi.New(sess, aws.NewConfig().WithEndpoint(os.Getenv("API_GATEWAY_URI")))

	newData, err := createCallbackSymbolMessage(symbolRate, symbol)
	if err != nil {
		return err
	}
	if newData == nil {
		log.Println("No new data")
		return errors.New("no data to broadcast")
	}

	newSymbolData := map[string][]CallbackSymbolMessage{}
	newSymbolData[symbol] = *newData

	// Convert websocket payload to json format
	byteMessage, err := json.Marshal(newSymbolData)
	if err != nil {
		log.Println("Marshalling error", err)
		return errors.New("callback message marshalling error")
	}

	for _, conn := range *connectionList {
		// Send new ohlc data to all users subscribed to the currency pair data channel over websocket api
		if _, err = apigw.PostToConnection(&apigatewaymanagementapi.PostToConnectionInput{ConnectionId: &conn.ConnectionId, Data: byteMessage}); err != nil {
			log.Println("Could not send to api", err)
			return errors.New("could not send to apigateway")
		}
	}

	return nil
}

// Creates websocket payload containing ohlc data
func createCallbackSymbolMessage(symbolRate *[]finance.FinancialDataItem, symbol string) (*[]CallbackSymbolMessage, error) {
	var newCallbackMessageData []CallbackSymbolMessage
	for _, rate := range *symbolRate {
		// Selects field of FinancialDataItem according to currency pair
		symbolField := getSymbolRateField(symbol, &rate)
		if checkIsRateValid(symbolField) {
			newData := CallbackSymbolMessage{
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

// Ohlc data is only valid if open, high, low, and close prices are non-zero
func checkIsRateValid(symbolField *finance.SymbolData) bool {
	if symbolField.Open == 0 ||
		symbolField.High == 0 ||
		symbolField.Low == 0 ||
		symbolField.Close == 0 {
		return false
	}
	return true
}

// Gets correct struct field using the selected symbol
// Avoids compile time type issues
func getSymbolRateField(symbol string, symbolRate *finance.FinancialDataItem) *finance.SymbolData {
	switch symbol {
	case "EURUSD":
		return &symbolRate.EURUSD
	case "GBPUSD":
		return &symbolRate.GBPUSD
	case "USDJPY":
		return &symbolRate.USDJPY
	case "AUDCAD":
		return &symbolRate.AUDCAD
	default:
		return nil
	}
}
