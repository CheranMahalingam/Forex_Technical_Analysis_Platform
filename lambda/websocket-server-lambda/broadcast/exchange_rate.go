package broadcast

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

type ExchangeRate struct {
	message map[string][]CallbackMessageSymbol
}

type exchangeRateTable struct {
	Date      string
	Timestamp string
	EURUSD    ohlcData
	GBPUSD    ohlcData
	USDJPY    ohlcData
	AUDCAD    ohlcData
}

type ohlcData struct {
	Open   float32
	High   float32
	Low    float32
	Close  float32
	Volume float32
}

func (er *ExchangeRate) Broadcast(event events.APIGatewayWebsocketProxyRequest, sess *session.Session, output *dynamodb.QueryOutput, connectionId string, optionalArg *string) error {
	symbol := *optionalArg

	rateList := make([]exchangeRateTable, *output.Count)
	err := dynamodbattribute.UnmarshalListOfMaps(output.Items, &rateList)
	if err != nil {
		log.Println("Could not unmarshal", err)
		return errors.New("json unmarshalling error")
	}
	er.message[symbol] = *createCallbackSymbolMessage(&rateList, symbol)

	// Converts payload to json
	byteMessage, err := json.Marshal(er.message)
	if err != nil {
		log.Println("Marshalling error", err)
		return errors.New("callback message marshalling error")
	}

	if err = sendMessageToWebsocketConnection(connectionId, event, sess, &byteMessage); err != nil {
		return err
	}

	return nil
}

func CreateExchangeRate() *ExchangeRate {
	return &ExchangeRate{message: make(map[string][]CallbackMessageSymbol)}
}

// Generates a query to DynamoDB to get 24 hours of ohlc data
func InitialOhlcData(symbol string, key expression.KeyConditionBuilder) (*dynamodb.QueryInput, error) {
	proj := expression.NamesList(expression.Name(symbol), expression.Name("Date"), expression.Name("Timestamp"))
	expr, err := expression.NewBuilder().
		WithKeyCondition(key).
		WithProjection(proj).
		Build()
	if err != nil {
		log.Printf("Subscription data expression error: %s", err)
		return nil, errors.New("expression error")
	}

	input := &dynamodb.QueryInput{
		TableName:                 aws.String("SymbolRateTable"),
		ProjectionExpression:      expr.Projection(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
	}

	return input, nil
}

// Gets field from ExchangeRateTable according to symbol
func getSymbolStructField(symbol string, symbolRate *exchangeRateTable) *ohlcData {
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

// Validates ohlc data
// open, high, low, and close price can never be zero
func checkIsRateValid(rate *ohlcData) bool {
	if rate.Open == 0 ||
		rate.High == 0 ||
		rate.Low == 0 ||
		rate.Close == 0 {
		return false
	}
	return true
}
