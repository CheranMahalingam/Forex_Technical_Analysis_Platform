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

func BroadcastMarketNews(connectionList *[]finance.Connection, marketNews *[]finance.NewsItem) error {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	apigw := apigatewaymanagementapi.New(sess, aws.NewConfig().WithEndpoint(os.Getenv("API_GATEWAY_URI")))

	newData := createCallbackNewsMessage(marketNews)

	newMarketData := map[string][]CallbackNewsMessage{}
	newMarketData["News"] = *newData

	byteMessage, err := json.Marshal(newMarketData)
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

func ValidateNewMarketNews(marketNews *[]finance.NewsItem, headline string) bool {
	if len(*marketNews) == 0 {
		return false
	}

	if len(*marketNews) == 1 && (*marketNews)[0].Headline == headline {
		return false
	}

	return true
}

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
