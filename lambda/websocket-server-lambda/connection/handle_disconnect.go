package connection

import (
	"errors"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// Removes user from DynamoDB WebsocketConnections table
func HandleDisconnect(connectionId string) error {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := dynamodb.New(sess)

	// Remove user by connectionId
	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"ConnectionId": {
				S: aws.String(connectionId),
			},
		},
		TableName: aws.String("WebsocketConnectionsTable"),
	}

	_, err := svc.DeleteItem(input)
	if err != nil {
		log.Printf("Connection String Deletion Error: %s", err)
		return errors.New("DynamoDB Delete Error")
	}

	return nil
}
