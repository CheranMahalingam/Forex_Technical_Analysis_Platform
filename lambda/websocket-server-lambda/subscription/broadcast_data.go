package subscription

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/apigatewaymanagementapi"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

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

func broadcastPrediction(connectionId string, event events.APIGatewayWebsocketProxyRequest, prediction *dynamodb.QueryOutput, sess *session.Session, symbol string) error {
	endpointUrl := "https://" + event.RequestContext.DomainName + "/" + event.RequestContext.Stage
	apigw := apigatewaymanagementapi.New(sess, aws.NewConfig().WithEndpoint(endpointUrl))

	websocketMessage := map[string]CallbackMessageInference{}
	log.Println(*prediction)
	inferenceList := make([]Inference, *prediction.Count)
	err := dynamodbattribute.UnmarshalListOfMaps(prediction.Items, &inferenceList)
	if err != nil {
		log.Println("Could not unmarshal", err)
		return errors.New("json unmarshalling error")
	}
	log.Println(inferenceList)
	websocketMessage[symbol] = *createCallbackInferenceMessage(&inferenceList, symbol)

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

func broadcastMarketNews(connectionId string, event events.APIGatewayWebsocketProxyRequest, marketNews *dynamodb.QueryOutput, sess *session.Session) error {
	endpointUrl := "https://" + event.RequestContext.DomainName + "/" + event.RequestContext.Stage
	apigw := apigatewaymanagementapi.New(sess, aws.NewConfig().WithEndpoint(endpointUrl))

	websocketMessage := map[string][]CallbackMessageNews{}
	newsList := []CallbackMessageNews{}
	err := dynamodbattribute.UnmarshalListOfMaps(marketNews.Items, &newsList)
	if err != nil {
		log.Println("Could not unmarshal", err)
		return errors.New("json unmarshalling error")
	}
	filteredNewsList := filterMarketNews(&newsList)
	websocketMessage["News"] = *filteredNewsList

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

func filterMarketNews(marketNews *[]CallbackMessageNews) *[]CallbackMessageNews {
	if len(*marketNews) == 0 {
		return marketNews
	}
	currentNews := (*marketNews)[0].MarketNews.Headline
	var newMarketNewsList = []CallbackMessageNews{(*marketNews)[0]}
	for _, news := range *marketNews {
		if news.MarketNews.Headline != currentNews {
			newMarketNewsList = append(newMarketNewsList, news)
			currentNews = news.MarketNews.Headline
		}
	}
	return &newMarketNewsList
}
