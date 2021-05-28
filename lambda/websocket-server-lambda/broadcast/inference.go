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

type Inference struct {
	message map[string]CallbackMessageInference
}

type inferenceTable struct {
	Latest          string
	EURUSDInference []float32
	GBPUSDInference []float32
	USDJPYInference []float32
	AUDCADInference []float32
}

func (i *Inference) Broadcast(event events.APIGatewayWebsocketProxyRequest, sess *session.Session, output *dynamodb.QueryOutput, connectionId string, optionalArg *string) error {
	symbol := *optionalArg

	inferenceList := make([]inferenceTable, *output.Count)
	err := dynamodbattribute.UnmarshalListOfMaps(output.Items, &inferenceList)
	if err != nil {
		log.Println("Could not unmarshal", err)
		return errors.New("json unmarshalling error")
	}
	i.message[symbol] = *createCallbackInferenceMessage(&inferenceList, symbol)

	// Converts payload to json
	byteMessage, err := json.Marshal(i.message)
	if err != nil {
		log.Println("Marshalling error", err)
		return errors.New("callback message marshalling error")
	}

	if err = sendMessageToWebsocketConnection(connectionId, event, sess, &byteMessage); err != nil {
		return err
	}

	return nil
}

func CreateInference() *Inference {
	return &Inference{message: make(map[string]CallbackMessageInference)}
}

// Generates a query to DynamoDB to get latest exchange rate forecast
func InitialInferenceData(symbol string) (*dynamodb.QueryInput, error) {
	colName := symbol + "Inference"
	proj := expression.NamesList(expression.Name(colName), expression.Name("Time"))
	keyCond := expression.Key("Date").Equal(expression.Value("inference"))
	expr, err := expression.NewBuilder().
		WithKeyCondition(keyCond).
		WithProjection(proj).
		Build()
	if err != nil {
		log.Printf("Subscription to inference data expression error: %s", err)
		return nil, errors.New("expression error")
	}

	input := &dynamodb.QueryInput{
		TableName:                 aws.String("TechnicalAnalysisTable"),
		ProjectionExpression:      expr.Projection(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
	}

	return input, nil
}

// Gets field from Inference according to symbol
func getInferenceStructField(symbol string, inference inferenceTable) *[]float32 {
	switch symbol {
	case "EURUSD":
		return &inference.EURUSDInference
	case "GBPUSD":
		return &inference.GBPUSDInference
	case "USDJPY":
		return &inference.USDJPYInference
	case "AUDCAD":
		return &inference.AUDCADInference
	default:
		return nil
	}
}
