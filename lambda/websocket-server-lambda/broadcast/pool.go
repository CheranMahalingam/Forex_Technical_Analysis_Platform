package broadcast

import (
	"errors"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/apigatewaymanagementapi"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type Pool interface {
	Broadcast(event events.APIGatewayWebsocketProxyRequest, sess *session.Session, output *dynamodb.QueryOutput, connectionId string, optionalArg *string) error
}

func sendMessageToWebsocketConnection(connectionId string, event events.APIGatewayWebsocketProxyRequest, sess *session.Session, byteMessage *[]byte) error {
	// Construct api endpoint from apiGateway event data
	endpointUrl := "https://" + event.RequestContext.DomainName + "/" + event.RequestContext.Stage
	apigw := apigatewaymanagementapi.New(sess, aws.NewConfig().WithEndpoint(endpointUrl))

	// Sends payload to user through websocket connection
	_, err := apigw.PostToConnection(&apigatewaymanagementapi.PostToConnectionInput{ConnectionId: &connectionId, Data: *byteMessage})
	if err != nil {
		log.Println("Could not send message to api", err)
		return errors.New("could not send to apigateway")
	}

	return nil
}
