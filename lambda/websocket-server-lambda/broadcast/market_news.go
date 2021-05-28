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

type MarketNews struct {
	message map[string][]CallbackMessageNews
}

type marketNewsTable struct {
	Timestamp string
	Headline  string
	Image     string
	Source    string
	Summary   string
	NewsUrl   string
}

func (mn *MarketNews) Broadcast(event events.APIGatewayWebsocketProxyRequest, sess *session.Session, output *dynamodb.QueryOutput, connectionId string, optionalArg *string) error {
	newsList := []CallbackMessageNews{}
	err := dynamodbattribute.UnmarshalListOfMaps(output.Items, &newsList)
	if err != nil {
		log.Println("Could not unmarshal", err)
		return errors.New("json unmarshalling error")
	}
	// Filters out headlines that repeat
	filteredNewsList := filterMarketNews(&newsList)
	mn.message["News"] = *filteredNewsList

	// Converts payload to json
	byteMessage, err := json.Marshal(mn.message)
	if err != nil {
		log.Println("Marshalling error", err)
		return errors.New("callback message marshalling error")
	}

	if err = sendMessageToWebsocketConnection(connectionId, event, sess, &byteMessage); err != nil {
		return err
	}

	return nil
}

func CreateMarketNews() *MarketNews {
	return &MarketNews{message: make(map[string][]CallbackMessageNews)}
}

// Generates a query to DynamoDB to get 24 hours of market news
func InitialMarketNews(key expression.KeyConditionBuilder) (*dynamodb.QueryInput, error) {
	// Reduces payload size and read speed by only getting data from MarketNews column
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

// Filters out market news headlines that repeat
// News must be filtered due to db structure, to avoid scanning db for news headlines latest news is
// stored in each row leading to copies
func filterMarketNews(marketNews *[]CallbackMessageNews) *[]CallbackMessageNews {
	if len(*marketNews) == 0 {
		return marketNews
	}
	currentNews := (*marketNews)[0].MarketNews.Headline
	var newMarketNewsList = []CallbackMessageNews{(*marketNews)[0]}
	for _, news := range *marketNews {
		// Compares current headline with previous unique headline
		if news.MarketNews.Headline != currentNews && len(news.MarketNews.Headline) != 0 {
			newMarketNewsList = append(newMarketNewsList, news)
			currentNews = news.MarketNews.Headline
		}
	}
	return &newMarketNewsList
}
