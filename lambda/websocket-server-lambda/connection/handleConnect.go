package connection

import (
	"errors"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

func HandleConnect(connectionId string) error {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	connectionItem := Connection{ConnectionId: connectionId, EurUsd: false, GbpUsd: false}

	svc := dynamodb.New(sess)

	newConnection, err := dynamodbattribute.MarshalMap(connectionItem)
	if err != nil {
		log.Printf("Error marshalling connection store: %s", err)
		return errors.New("Connection Error")
	}

	input := &dynamodb.PutItemInput{
		Item:      newConnection,
		TableName: aws.String("WebsocketConnectionsTable"),
	}

	_, err = svc.PutItem(input)
	if err != nil {
		log.Printf("Error pushing data to dynamoDB: %s", err)
		return errors.New("DynamoDB Push Error")
	}

	return nil
}
