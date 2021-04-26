package subscription

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/apigatewaymanagementapi"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

type ExchangeRate struct {
	Open   float32
	High   float32
	Low    float32
	Close  float32
	Volume float32
}

type dynamodbEurUsd struct {
	Date      string
	Timestamp string
	EURUSD    ExchangeRate
}

type dynamodbGbpUsd struct {
	Date      string
	Timestamp string
	GBPUSD    ExchangeRate
}

func HandleSubscription(connectionId string, event events.APIGatewayWebsocketProxyRequest, subscribe bool) error {
	log.Println(event.Body)
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	msg := SubscriptionMessage{}
	err := json.Unmarshal([]byte(event.Body), &msg)
	if err != nil {
		log.Println(err)
		return errors.New("json unmarshaling error")
	}

	svc := dynamodb.New(sess)

	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":symbol": {
				S: aws.String(strconv.FormatBool(subscribe)),
			},
		},
		TableName: aws.String("WebsocketConnectionsTable"),
		Key: map[string]*dynamodb.AttributeValue{
			"ConnectionId": {
				S: aws.String(connectionId),
			},
		},
		UpdateExpression: aws.String(fmt.Sprintf("set %s = :symbol", msg.Data)),
	}

	_, err = svc.UpdateItem(input)
	if err != nil {
		log.Printf("Error updating websocket subscription: %s", err)
		return errors.New("DynamoDB Update Error")
	}

	if subscribe {
		newDataQuery, err := initialSubscriptionData(msg.Data, connectionId, 1)
		if err != nil {
			return err
		}

		subscriptionData, err := svc.Query(newDataQuery)
		if err != nil {
			log.Printf("dynamodb querying error: %s", err)
			return err
		}

		if err = broadcastSubscribedData(connectionId, event, subscriptionData, sess, msg.Data); err != nil {
			return err
		}
	}

	return nil
}

func initialSubscriptionData(symbol string, connectionId string, dayCount int) (*dynamodb.QueryInput, error) {
	currentDate := time.Now().Format("2006-01-02")
	log.Println(symbol)
	keyCond := expression.Key("Date").Equal(expression.Value(currentDate))
	proj := expression.NamesList(expression.Name(symbol), expression.Name("Date"), expression.Name("Timestamp"))
	expr, err := expression.NewBuilder().
		WithKeyCondition(keyCond).
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

func broadcastSubscribedData(connectionId string, event events.APIGatewayWebsocketProxyRequest, symbolRate *dynamodb.QueryOutput, sess *session.Session, symbol string) error {
	endpointUrl := "https://" + event.RequestContext.DomainName + "/" + event.RequestContext.Stage
	apigw := apigatewaymanagementapi.New(sess, aws.NewConfig().WithEndpoint(endpointUrl))

	websocketMessage := map[string][]CallbackMessageData{}
	switch symbol {
	case "EURUSD":
		rateList := make([]dynamodbEurUsd, *symbolRate.Count)
		err := dynamodbattribute.UnmarshalListOfMaps(symbolRate.Items, &rateList)
		if err != nil {
			log.Println("Could not unmarshal", err)
			return errors.New("json unmarshalling error")
		}
		log.Println(rateList)
		websocketMessage[symbol] = *createCallbackMessageEurUsd(&rateList)
	case "GBPUSD":
		rateList := make([]dynamodbGbpUsd, *symbolRate.Count)
		err := dynamodbattribute.UnmarshalListOfMaps(symbolRate.Items, &rateList)
		if err != nil {
			log.Println("Could not unmarshal", err)
			return errors.New("json unmarshalling error")
		}
		log.Println(rateList)
		websocketMessage[symbol] = *createCallbackMessageGbpUsd(&rateList)
	default:
		return errors.New("invalid symbol")
	}

	byteMessage, err := json.Marshal(websocketMessage)
	if err != nil {
		log.Println("Marshalling error", err)
		return errors.New("callback message marshalling error")
	}

	if _, err = apigw.PostToConnection(&apigatewaymanagementapi.PostToConnectionInput{ConnectionId: &connectionId, Data: byteMessage}); err != nil {
		log.Println("Could not send message to api", err)
		return errors.New("could not send to apigateway")
	}

	return nil
}

func createCallbackMessageEurUsd(symbolRate *[]dynamodbEurUsd) *[]CallbackMessageData {
	var newCallbackMessageData []CallbackMessageData
	for _, rate := range *symbolRate {
		if checkIsRateValidEurUsd(&rate) {
			newData := CallbackMessageData{
				Timestamp: rate.Date + " " + rate.Timestamp,
				Open:      rate.EURUSD.Open,
				High:      rate.EURUSD.High,
				Low:       rate.EURUSD.Low,
				Close:     rate.EURUSD.Close,
			}
			newCallbackMessageData = append(newCallbackMessageData, newData)
		}
	}
	return &newCallbackMessageData
}

func createCallbackMessageGbpUsd(symbolRate *[]dynamodbGbpUsd) *[]CallbackMessageData {
	var newCallbackMessageData []CallbackMessageData
	for _, rate := range *symbolRate {
		if checkIsRateValidGbpUsd(&rate) {
			newData := CallbackMessageData{
				Timestamp: rate.Date + " " + rate.Timestamp,
				Open:      rate.GBPUSD.Open,
				High:      rate.GBPUSD.High,
				Low:       rate.GBPUSD.Low,
				Close:     rate.GBPUSD.Close,
			}
			newCallbackMessageData = append(newCallbackMessageData, newData)
		}
	}
	return &newCallbackMessageData
}

func checkIsRateValidEurUsd(rate *dynamodbEurUsd) bool {
	if rate.EURUSD.Open == 0 ||
		rate.EURUSD.High == 0 ||
		rate.EURUSD.Low == 0 ||
		rate.EURUSD.Close == 0 {
		return false
	}
	return true
}

func checkIsRateValidGbpUsd(rate *dynamodbGbpUsd) bool {
	if rate.GBPUSD.Open == 0 ||
		rate.GBPUSD.High == 0 ||
		rate.GBPUSD.Low == 0 ||
		rate.GBPUSD.Close == 0 {
		return false
	}
	return true
}
