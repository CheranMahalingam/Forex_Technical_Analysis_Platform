package subscription

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"websocket-server-lambda/broadcast"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

// Controller function for allowing user to subscribe and unsubscribe to data channels
// Parses websocket message payload to check which data channel the user subscribes to
// New subscribers are sent past data since it will take time for new data to be received and
// no data is streamed during the time the forex markets are closed
func HandleSubscription(connectionId string, event events.APIGatewayWebsocketProxyRequest, subscribe bool) error {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	msg := SubscriptionMessage{}
	// Decodes json message so that data can be parsed
	err := json.Unmarshal([]byte(event.Body), &msg)
	if err != nil {
		log.Println(err)
		return errors.New("json unmarshaling error")
	}

	svc := dynamodb.New(sess)

	// Updates user websocket subscription information
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

	if subscribe && permissionToGetRates(msg.Data) {
		// Occurs if user subscribes to ohlc data channel
		keyCondList := keyConditionDayData()
		for _, condition := range *keyCondList {
			newDataQuery, err := broadcast.InitialOhlcData(msg.Data, condition)
			if err != nil {
				return err
			}

			subscriptionData, err := svc.Query(newDataQuery)
			if err != nil {
				log.Printf("dynamodb querying error: %s", err)
				return err
			}

			exchangeRateMessage := broadcast.CreateExchangeRate()
			if err = exchangeRateMessage.Broadcast(event, sess, subscriptionData, connectionId, &msg.Data); err != nil {
				return err
			}
		}
	} else if subscribe && permissionToGetInference(msg.Data) {
		// Occurs if user subscribes to an exchange rate forecast channel
		symbol := msg.Data[:6]
		newInferenceQuery, err := broadcast.InitialInferenceData(symbol)
		if err != nil {
			return err
		}

		inferenceData, err := svc.Query(newInferenceQuery)
		if err != nil {
			log.Printf("dynamodb querying error: %s", err)
			return err
		}

		inferenceMessage := broadcast.CreateInference()
		if err = inferenceMessage.Broadcast(event, sess, inferenceData, connectionId, &symbol); err != nil {
			return err
		}
	} else if subscribe && permissionToGetMarketNews(msg.Data) {
		// Occurs if user subscribes to market news channel
		keyCondList := keyConditionDayData()
		for _, condition := range *keyCondList {
			newNewsQuery, err := broadcast.InitialMarketNews(condition)
			if err != nil {
				return err
			}

			marketNews, err := svc.Query(newNewsQuery)
			if err != nil {
				log.Printf("dynamodb querying error: %s", err)
				return err
			}

			newsMessage := broadcast.CreateMarketNews()
			if err = newsMessage.Broadcast(event, sess, marketNews, connectionId, nil); err != nil {
				return err
			}
		}
	}

	return nil
}

// Generates key conditions to query DynamoDB
// Generally data from the past 24 hours is needed, however, forex markets are closed from
// Friday 5:00PM EST to Sunday 5:00PM EST
// During the weekend data from Friday and Thursday may be queried to provided 24 hours of data
// Due to the db structure two reads are needed to fetch all data
func keyConditionDayData() *[]expression.KeyConditionBuilder {
	currentDay := int(time.Now().Weekday())
	currentHour := time.Now().Hour()

	if currentDay == 6 {
		// On Saturday data from Friday and data from Thursday above the current hour is provided to users
		currentDate := time.Now().Add(-time.Hour * 24)
		previousDate := time.Now().Add(-time.Hour * 48)
		keyCond := expression.Key("Date").Equal(expression.Value(currentDate.Format("2006-01-02")))
		prevKeyCond := expression.Key("Date").Equal(expression.Value(previousDate.Format("2006-01-02")))
		keyCondList := []expression.KeyConditionBuilder{prevKeyCond, keyCond}
		return &keyCondList

	} else if currentDay == 0 && currentHour > 22 {
		// On Sunday after 5:00PM data from Sunday, Friday, and Thursday above the current hour is
		// provided to users since Sunday data will be minimal when the markets open
		currentDate := time.Now().Add(-time.Hour * 0)
		previousDate := time.Now().Add(-time.Hour * 48)
		extraDate := time.Now().Add(-time.Hour * 72)
		keyCond := expression.Key("Date").Equal(expression.Value(currentDate.Format("2006-01-02")))
		prevKeyCond := expression.Key("Date").Equal(expression.Value(previousDate.Format("2006-01-02")))
		extraKeyCond := expression.KeyAnd(
			expression.Key("Date").Equal(expression.Value(extraDate.Format("2006-01-02"))),
			expression.Key("Timestamp").GreaterThanEqual(expression.Value(extraDate.Format("03:04:05"))),
		)
		keyCondList := []expression.KeyConditionBuilder{extraKeyCond, prevKeyCond, keyCond}
		return &keyCondList

	} else if currentDay == 0 {
		// On Sunday data from Friday and Thursday above the current hour is provided to users
		currentDate := time.Now().Add(-time.Hour * 48)
		previousDate := time.Now().Add(-time.Hour * 72)
		keyCond := expression.Key("Date").Equal(expression.Value(currentDate.Format("2006-01-02")))
		prevKeyCond := expression.Key("Date").Equal(expression.Value(previousDate.Format("2006-01-02")))
		keyCondList := []expression.KeyConditionBuilder{prevKeyCond, keyCond}
		return &keyCondList

	} else if currentDay == 1 {
		// On Monday data from Monday, Sunday, and Friday above the current hour is provided to users
		currentDate := time.Now().Add(-time.Hour * 0)
		previousDate := time.Now().Add(-time.Hour * 24)
		extraDate := time.Now().Add(-time.Hour * 72)
		keyCond := expression.Key("Date").Equal(expression.Value(currentDate.Format("2006-01-02")))
		prevKeyCond := expression.Key("Date").Equal(expression.Value(previousDate.Format("2006-01-02")))
		extraKeyCond := expression.Key("Date").Equal(expression.Value(extraDate.Format("2006-01-02")))
		keyCondList := []expression.KeyConditionBuilder{extraKeyCond, prevKeyCond, keyCond}
		return &keyCondList

	} else {
		// For all other days data from the current day and data from the previous day above the current
		// hour is provided to users, results in 24 hours of data
		currentDate := time.Now().Add(-time.Hour * 0)
		previousDate := time.Now().Add(-time.Hour * 24)
		keyCond := expression.Key("Date").Equal(expression.Value(currentDate.Format("2006-01-02")))
		prevKeyCond := expression.KeyAnd(
			expression.Key("Date").Equal(expression.Value(previousDate.Format("2006-01-02"))),
			expression.Key("Timestamp").GreaterThanEqual(expression.Value(previousDate.Format("03:04:05"))),
		)
		keyCondList := []expression.KeyConditionBuilder{prevKeyCond, keyCond}
		return &keyCondList
	}
}
