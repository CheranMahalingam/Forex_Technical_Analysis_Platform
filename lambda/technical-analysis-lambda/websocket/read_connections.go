package websocket

import (
	"errors"
	"log"
	"technical-analysis-lambda/finance"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

// Creates a slice of all websocket connections subscribed to a data channel
func ReadConnections(dbColumn string) (*[]finance.Connection, error) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := dynamodb.New(sess)

	// Filters out users who are subscribed to the dbColumn data channel
	filt := expression.Name(dbColumn).Equal(expression.Value("true"))
	expr, err := expression.NewBuilder().WithFilter(filt).Build()
	if err != nil {
		return nil, errors.New("could not create filter expression")
	}
	params := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		TableName:                 aws.String("WebsocketConnectionsTable"),
	}

	result, err := svc.Scan(params)
	if err != nil {
		log.Println("Scan Failed", err)
		return nil, errors.New("DynamoDB scan failed")
	}

	connectionList := make([]finance.Connection, *result.Count)
	if err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &connectionList); err != nil {
		return nil, errors.New("connection list unmarshaling error")
	}

	return &connectionList, nil
}
