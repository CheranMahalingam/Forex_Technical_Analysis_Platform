package dynamosymbol

import (
	"errors"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/apigatewaymanagementapi"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

func ReadSymbolRate(symbol string, symbolRate *[]byte) error {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := dynamodb.New(sess)

	filt := expression.Name("EurUsd").Equal(expression.Value("true"))
	expr, err := expression.NewBuilder().WithFilter(filt).Build()
	if err != nil {
		return errors.New("Could not create filter expression")
	}
	params := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		TableName:                 aws.String("WebsocketConnectionsTable"),
	}

	apigw := apigatewaymanagementapi.New(sess, aws.NewConfig().WithEndpoint(os.Getenv("ApiGatewayUri")))

	result, err := svc.Scan(params)
	if err != nil {
		log.Println("Scan Failed LMAO", err)
		return errors.New("DynamoDB scan failed")
	}

	connectionList := make([]Connection, *result.Count)
	if err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &connectionList); err != nil {
		return errors.New("Connection list unmarshaling error")
	}

	for _, c := range connectionList {
		log.Println(c.ConnectionId)
		_, err = apigw.PostToConnection(&apigatewaymanagementapi.PostToConnectionInput{ConnectionId: &c.ConnectionId, Data: *symbolRate})
		if err != nil {
			log.Println("Could not send to api", err)
			return errors.New("Could not send to apigatway")
		}
	}

	return nil
}
