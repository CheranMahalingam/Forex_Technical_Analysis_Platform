import boto3
from boto3.dynamodb.conditions import Attr
from botocore.exceptions import ClientError
import os
import json
import datetime
from decimal import Decimal


def scan_connections(symbol):
    dynamodb = boto3.resource('dynamodb')
    table = dynamodb.Table('WebsocketConnectionsTable')

    response = table.scan(
        FilterExpression=Attr('{}Inference'.format(symbol)).eq('true'),
        ProjectionExpression="ConnectionId"
    )
    print(response)

    connection_ids = []
    for item in response['Items']:
        connection_ids.append(item['ConnectionId'])
    return connection_ids


def broadcast_inference(symbol, connections, inference, date):
    apig_management_client = boto3.client(
        'apigatewaymanagementapi', endpoint_url=os.getenv("API_GATEWAY_URI")
    )

    current_date = datetime.datetime.strptime(date, "%Y-%m-%d%H:%M:%S")
    current_date = datetime.datetime.strftime(current_date, "%Y-%m-%d %H:%M:%S")
    inference = inference.tolist()
    data = {symbol: {"inference": inference, "date": current_date}}

    for conn_id in connections:
        try:
            send_response = apig_management_client.post_to_connection(
                Data=json.dumps(data), ConnectionId=conn_id
            )
        except ClientError:
            print("Couldn't post to connection", conn_id)
        except apig_management_client.exceptions.GoneException:
            print("Connection {} is gone".format(conn_id))


def store_inference(symbol, inference, timestamp):
    dynamodb = boto3.resource('dynamodb')
    table = dynamodb.Table('TechnicalAnalysisTable')

    inference = inference.tolist()
    for i in range(len(inference)):
        inference[i] = Decimal(str(inference[i]))

    response = table.update_item(
        Key={
            'Date': 'inference',
            'Timestamp': 'timestamp'
        },
        UpdateExpression="set {}Inference=:i, latest=:l".format(symbol),
        ExpressionAttributeValues={
            ':i': inference,
            ':l': timestamp
        }
    )
    print(response)
