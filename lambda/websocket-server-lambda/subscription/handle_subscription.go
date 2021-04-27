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

type exchangeRate struct {
	Open   float32
	High   float32
	Low    float32
	Close  float32
	Volume float32
}

type ExchangeRateTable struct {
	Date      string
	Timestamp string
	EURUSD    exchangeRate
	GBPUSD    exchangeRate
}

func HandleSubscription(connectionId string, event events.APIGatewayWebsocketProxyRequest, subscribe bool) error {
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
		currentDate := time.Now()
		previousDay := currentDate.Add(-time.Hour * 48)
		log.Println(currentDate)
		log.Println(previousDay)
		keyCond := expression.Key("Date").Equal(expression.Value(currentDate.Format("2006-01-02")))
		prevKeyCond := expression.KeyAnd(
			expression.Key("Date").Equal(expression.Value(previousDay.Format("2006-01-02"))),
			expression.Key("Timestamp").GreaterThan(expression.Value(previousDay.Format("03:04:05"))),
		)
		keyCondList := []expression.KeyConditionBuilder{prevKeyCond, keyCond}

		for _, condition := range keyCondList {
			newDataQuery, err := initialSubscriptionData(msg.Data, connectionId, condition)
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
	}

	return nil
}

func initialSubscriptionData(symbol string, connectionId string, key expression.KeyConditionBuilder) (*dynamodb.QueryInput, error) {
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

func broadcastSubscribedData(connectionId string, event events.APIGatewayWebsocketProxyRequest, symbolRate *dynamodb.QueryOutput, sess *session.Session, symbol string) error {
	endpointUrl := "https://" + event.RequestContext.DomainName + "/" + event.RequestContext.Stage
	apigw := apigatewaymanagementapi.New(sess, aws.NewConfig().WithEndpoint(endpointUrl))

	websocketMessage := map[string][]CallbackMessageData{}
	rateList := make([]ExchangeRateTable, *symbolRate.Count)
	err := dynamodbattribute.UnmarshalListOfMaps(symbolRate.Items, &rateList)
	if err != nil {
		log.Println("Could not unmarshal", err)
		return errors.New("json unmarshalling error")
	}
	log.Println(rateList)
	websocketMessage[symbol] = *createCallbackMessage(&rateList, symbol)

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

func createCallbackMessage(symbolRate *[]ExchangeRateTable, symbol string) *[]CallbackMessageData {
	var newCallbackMessageData []CallbackMessageData
	for _, rate := range *symbolRate {
		symbolField := getSymbolStructField(symbol, &rate)
		if checkIsRateValid(symbolField) {
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
	return &newCallbackMessageData
}

func getSymbolStructField(symbol string, symbolRate *ExchangeRateTable) *exchangeRate {
	switch symbol {
	case "EURUSD":
		return &symbolRate.EURUSD
	case "GBPUSD":
		return &symbolRate.GBPUSD
	default:
		return nil
	}
}

func checkIsRateValid(rate *exchangeRate) bool {
	if rate.Open == 0 ||
		rate.High == 0 ||
		rate.Low == 0 ||
		rate.Close == 0 {
		return false
	}
	return true
}
