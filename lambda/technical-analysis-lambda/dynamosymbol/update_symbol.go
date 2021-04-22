package dynamosymbol

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

func PutNewSymbolRate(newSymbolRate *SymbolRateItem) {
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

func SendRateToDB(newSymbolRates *[]SymbolRateItem) {
	for _, newEntry := range *newSymbolRates {
		PutNewSymbolRate(&newEntry)
	}
}
