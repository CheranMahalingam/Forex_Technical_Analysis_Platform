package finance

import (
	"errors"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

// Controller function used to get latest market news headline
func GetNewsHeadline() (*string, error) {
	currentTime := time.Now().Format("2006-01-02")
	// Search for latest headline for current date
	newsHeadline, err := queryMarketNewsHeadline(currentTime)
	if err != nil {
		return nil, err
	} else if newsHeadline == nil {
		// If there were no headlines for current date check for headline during previous day
		previousTime := time.Now().Add(-24 * time.Hour).Format("2006-01-02")
		newsHeadline, err = queryMarketNewsHeadline(previousTime)
		if err != nil {
			return nil, err
		} else if newsHeadline == nil {
			return nil, nil
		}
		return newsHeadline, nil
	}
	return newsHeadline, nil
}

// Queries DynamoDB to get latest news headline
func queryMarketNewsHeadline(date string) (*string, error) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := dynamodb.New(sess)

	// Search along sort key in reverse to get latest news and first entry will be latest headline
	scanIndexForward := false
	itemLimit := int64(1)

	expr, err := createExpression(date)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	params := &dynamodb.QueryInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		KeyConditionExpression:    expr.KeyCondition(),
		TableName:                 aws.String("SymbolRateTable"),
		ScanIndexForward:          &scanIndexForward,
		Limit:                     &itemLimit,
	}

	result, err := svc.Query(params)
	if err != nil {
		log.Println("Scan Failed", err)
		return nil, nil
	}

	marketData := make([]FinancialDataItem, *result.Count)
	if err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &marketData); err != nil {
		return nil, errors.New("financial data unmarshalling error")
	}

	// There may not be a latest news headline for the date given if a new day has started
	if len(marketData) > 0 {
		return &marketData[0].MarketNews.Headline, nil
	}
	return nil, nil
}

// Create DynamoDB query expression for data gotten during given date
func createExpression(formattedDate string) (*expression.Expression, error) {
	keyCond := expression.Key("Date").Equal(expression.Value(formattedDate))
	expr, err := expression.NewBuilder().WithKeyCondition(keyCond).Build()
	if err != nil {
		return nil, errors.New("could not create filter expression")
	}

	return &expr, nil
}
