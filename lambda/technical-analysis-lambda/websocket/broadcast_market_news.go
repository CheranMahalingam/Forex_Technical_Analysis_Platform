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

// Broadcasts new market news to all subscribed users
// All users on the market news webpage are automatically subscribed to the market news channel
func BroadcastMarketNews(connectionList *[]finance.Connection, marketNews *[]finance.NewsItem) error {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	apigw := apigatewaymanagementapi.New(sess, aws.NewConfig().WithEndpoint(os.Getenv("API_GATEWAY_URI")))

	newData := createCallbackNewsMessage(marketNews)

	newMarketData := map[string][]CallbackNewsMessage{}
	newMarketData["News"] = *newData

	// Converts websocket payload to json format
	byteMessage, err := json.Marshal(newMarketData)
	if err != nil {
		log.Println("Marshalling error", err)
		return errors.New("callback message marshalling error")
	}

	for _, conn := range *connectionList {
		// Send new market news to all users subscribed to the global market news data channel
		if _, err = apigw.PostToConnection(&apigatewaymanagementapi.PostToConnectionInput{ConnectionId: &conn.ConnectionId, Data: byteMessage}); err != nil {
			log.Println("Could not send to api", err)
			return errors.New("could not send to apigateway")
		}
	}

	return nil
}

// Verfies that the market news received is different from the previous news headline
func ValidateNewMarketNews(marketNews *[]finance.NewsItem, headline string) bool {
	// Avoid dereferencing null pointer later
	if marketNews == nil {
		return false
	}

	// No news received
	if len(*marketNews) == 0 {
		return false
	}

	// News headline matches previous news headline
	if len(*marketNews) == 1 && (*marketNews)[0].Headline == headline {
		return false
	}

	return true
}

// Creates webscoket payloads containing relevant news information
func createCallbackNewsMessage(marketNews *[]finance.NewsItem) *[]CallbackNewsMessage {
	var newCallbackMessageData []CallbackNewsMessage
	for _, news := range *marketNews {
		newData := CallbackNewsMessage{
			Timestamp: news.Timestamp,
			Headline:  news.Headline,
			Image:     news.Image,
			Source:    news.Source,
			Summary:   news.Summary,
			NewsUrl:   news.NewsUrl,
		}
		newCallbackMessageData = append(newCallbackMessageData, newData)
	}
	return &newCallbackMessageData
}
