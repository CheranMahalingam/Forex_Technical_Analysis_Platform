"""Module for getting past technical indicator and market sentiment data from DynamoDB."""

from dynamodb_json import json_util as json
import pandas as pd
import boto3
from boto3.dynamodb.conditions import Key
import datetime


def get_technical_analysis_data(date, timestamp):
    """
    Gets past 24 hours of technical indicator and market sentiment data from DynamoDB.

    Args:
        date: String of the form YYYY-mm-dd
        timestamp: String of the form HH:MM:SS
    
    Returns:
        Pandas dataframe containing technical indicator and market sentiment data
    """
    dynamodb = boto3.resource('dynamodb')
    table = dynamodb.Table('TechnicalAnalysisTable')

    # Gets all table data for current date
    response_present = table.query(
        KeyConditionExpression=Key('Date').eq(date)
    )

    previous_date = datetime.datetime.strptime(date, "%Y-%m-%d") - datetime.timedelta(days=1)
    previous_date = datetime.datetime.strftime(previous_date, "%Y-%m-%d")
    # Gets all table data from previous day where the hour is greater than
    # the current hour
    response_past = table.query(
        KeyConditionExpression=
            Key('Date').eq(previous_date) & Key('Timestamp').gte(timestamp)
    )

    present_df = convert_to_dataframe(response_present['Items'])
    previous_df = convert_to_dataframe(response_past['Items'])

    if present_df.empty and previous_df.empty:
        raise Exception("No data received from db")
    elif present_df.empty:
        return previous_df
    elif previous_df.empty:
        return present_df
    else:
        # Concatenate both dataframes if they both contain data
        new_df = pd.concat([previous_df, present_df], ignore_index=True)
        return new_df


def convert_to_dataframe(dynamo_db_data):
    """
    Converts data from DynamoDB into Pandas dataframe.

    Args:
        dynamo_db_data: Raw data from DynomoDB containing technical indicators
        and market sentiment data
    
    Returns:
        Pandas dataframe containing technical indicator and market sentiment data
    """
    if not dynamo_db_data:
        return pd.DataFrame()
    return pd.DataFrame(json.loads(dynamo_db_data))
