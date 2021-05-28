package finance

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

// Controller for processing financial data into a form such that it can be stored in db
func SendRateToDB(newSymbolRates *[]FinancialDataItem, latestNews *[]NewsItem) {
	// Market news is added to the financial data struct
	newFinancialData := mergeNewsItems(newSymbolRates, latestNews)
	for _, newEntry := range *newFinancialData {
		putNewSymbolRate(&newEntry)
	}
}

// Stores new ohlc data in DynamoDB
func putNewSymbolRate(newSymbolRate *FinancialDataItem) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := dynamodb.New(sess)

	newRate, err := dynamodbattribute.MarshalMap(newSymbolRate)
	if err != nil {
		log.Fatalf("Error marshalling exchange rate: %s", err)
	}

	input := &dynamodb.PutItemInput{
		Item:      newRate,
		TableName: aws.String("SymbolRateTable"),
	}

	_, err = svc.PutItem(input)
	if err != nil {
		log.Fatalf("Error pushing data to dynamoDB: %s", err)
	}
}

// Adds market news to financial data slice
func mergeNewsItems(newSymbolRates *[]FinancialDataItem, latestNews *[]NewsItem) *[]FinancialDataItem {
	// Avoids null pointer dereferencing error
	if latestNews == nil {
		return newSymbolRates
	}
	// MarketNews fields are filled with collected headlines
	for index, news := range *latestNews {
		if len(*newSymbolRates) > index {
			(*newSymbolRates)[index].MarketNews = news
		} else {
			break
		}
	}
	return newSymbolRates
}
