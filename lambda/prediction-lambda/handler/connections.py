"""
Module provides methods to get websocket connections and broadcast generated forecasts
to users.
"""

import boto3
from boto3.dynamodb.conditions import Attr
from botocore.exceptions import ClientError
import os
import json
import datetime
from decimal import Decimal


def scan_connections(symbol):
    """
    Gets the connectionId of websockets where the user is subscribed to the specified
    currency pair.

    Args:
        symbol: String representing a currency pair
    
    Returns:
        List of strings representing connectionIds
    """
    dynamodb = boto3.resource('dynamodb')
    table = dynamodb.Table('WebsocketConnectionsTable')

    # Scans for connections where the user has subscribed to this currency pair
    response = table.scan(
        FilterExpression=Attr('{}Inference'.format(symbol)).eq('true'),
        ProjectionExpression="ConnectionId"
    )

    connection_ids = []
    for item in response['Items']:
        connection_ids.append(item['ConnectionId'])
    return connection_ids


def broadcast_inference(symbol, connections, inference, date):
    """
    Method to broadcast new currency pair forecasts to users over websockets. Data is
    only sent to users who subscribed to a specific currency pair.

    Args:
        symbol: String representing a specific currency pair (e.g. EURUSD, GBPUSD...)
        connections: List of strings representing websocket connectionIds
        inference: Numpy array containing exchange rate forecast
        date: String representing current time in the RFC3339 format
    """
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
    """
    Stores inference in DynamoDB for quick future retrieval.

    Args:
        symbol: String representing a specific currency pair (e.g. EURUSD, GBPUSD...)
        inference: Numpy array containing exchange rate forecast
        timestamp: String representing time in the format HH:MM:SS
    """
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
