package finance

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

func SendRateToDB(newSymbolRates *[]FinancialDataItem, latestNews *[]NewsItem) {
	newFinancialData := mergeNewsItems(newSymbolRates, latestNews)
	for _, newEntry := range *newFinancialData {
		putNewSymbolRate(&newEntry)
	}
}

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

func mergeNewsItems(newSymbolRates *[]FinancialDataItem, latestNews *[]NewsItem) *[]FinancialDataItem {
	for index, news := range *latestNews {
		if len(*newSymbolRates) > index {
			(*newSymbolRates)[index].MarketNews = news
		} else {
			break
		}
	}

	return newSymbolRates
}
