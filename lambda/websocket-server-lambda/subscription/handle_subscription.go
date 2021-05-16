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
	"github.com/aws/aws-sdk-go/service/dynamodb"
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
	USDJPY    exchangeRate
	AUDCAD    exchangeRate
}

type Inference struct {
	Latest          string
	EURUSDInference []float32
	GBPUSDInference []float32
	USDJPYInference []float32
	AUDCADInference []float32
}

type newsItem struct {
	Timestamp string
	Headline  string
	Image     string
	Source    string
	Summary   string
	NewsUrl   string
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

	log.Println(msg.Data)
	if subscribe && permissionToGetRates(msg.Data) {
		keyCondList := keyConditionDayData()
		for _, condition := range *keyCondList {
			newDataQuery, err := initialSubscriptionData(msg.Data, condition)
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
	} else if subscribe && permissionToGetInference(msg.Data) {
		symbol := msg.Data[:6]
		newInferenceQuery, err := initialInferenceData(symbol)
		if err != nil {
			return err
		}

		inferenceData, err := svc.Query(newInferenceQuery)
		if err != nil {
			log.Printf("dynamodb querying error: %s", err)
			return err
		}
		log.Println(inferenceData)

		if err = broadcastPrediction(connectionId, event, inferenceData, sess, symbol); err != nil {
			return err
		}
	} else if subscribe && permissionToGetMarketNews(msg.Data) {
		keyCondList := keyConditionDayData()
		for _, condition := range *keyCondList {
			newNewsQuery, err := initialMarketNews(condition)
			if err != nil {
				return err
			}

			marketNews, err := svc.Query(newNewsQuery)
			if err != nil {
				log.Printf("dynamodb querying error: %s", err)
				return err
			}

			if err = broadcastMarketNews(connectionId, event, marketNews, sess); err != nil {
				return err
			}
		}
	}

	return nil
}

func initialSubscriptionData(symbol string, key expression.KeyConditionBuilder) (*dynamodb.QueryInput, error) {
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

func initialInferenceData(symbol string) (*dynamodb.QueryInput, error) {
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

func initialMarketNews(key expression.KeyConditionBuilder) (*dynamodb.QueryInput, error) {
	proj := expression.NamesList(expression.Name("MarketNews"))
	expr, err := expression.NewBuilder().
		WithKeyCondition(key).
		WithProjection(proj).
		Build()
	if err != nil {
		log.Printf("Market news expression error: %s", err)
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

func getSymbolStructField(symbol string, symbolRate *ExchangeRateTable) *exchangeRate {
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

func getInferenceStructField(symbol string, inference Inference) *[]float32 {
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

func checkIsRateValid(rate *exchangeRate) bool {
	if rate.Open == 0 ||
		rate.High == 0 ||
		rate.Low == 0 ||
		rate.Close == 0 {
		return false
	}
	return true
}

func keyConditionDayData() *[]expression.KeyConditionBuilder {
	currentDay := int(time.Now().Weekday())
	var currentDate time.Time
	var previousDate time.Time

	if currentDay == 6 {
		currentDate = time.Now().Add(-time.Hour * 24)
		previousDate = time.Now().Add(-time.Hour * 48)
	} else if currentDay == 0 {
		currentDate = time.Now().Add(-time.Hour * 48)
		previousDate = time.Now().Add(-time.Hour * 72)
	} else if currentDay == 1 {
		currentDate = time.Now().Add(-time.Hour * 0)
		previousDate = time.Now().Add(-time.Hour * 72)
	} else {
		currentDate = time.Now().Add(-time.Hour * 0)
		previousDate = time.Now().Add(-time.Hour * 24)
	}

	keyCond := expression.Key("Date").Equal(expression.Value(currentDate.Format("2006-01-02")))
	prevKeyCond := expression.KeyAnd(
		expression.Key("Date").Equal(expression.Value(previousDate.Format("2006-01-02"))),
		expression.Key("Timestamp").GreaterThanEqual(expression.Value(previousDate.Format("03:04:05"))),
	)
	keyCondList := []expression.KeyConditionBuilder{prevKeyCond, keyCond}

	return &keyCondList
}
