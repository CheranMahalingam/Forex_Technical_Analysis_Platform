from dynamodb_json import json_util as json
import pandas as pd
import boto3
from boto3.dynamodb.conditions import Key
import datetime


def get_technical_analysis_data(date, timestamp):
    dynamodb = boto3.resource('dynamodb')
    table = dynamodb.Table('TechnicalAnalysisTable')

    response_present = table.query(
        KeyConditionExpression=Key('Date').eq(date)
    )

    previous_date = datetime.datetime.strptime(date, "%Y-%m-%d") - datetime.timedelta(days=1)
    previous_date = datetime.datetime.strftime(previous_date, "%Y-%m-%d")
    response_past = table.query(
        KeyConditionExpression=
            Key('Date').eq(previous_date) & Key('Timestamp').gt(timestamp)
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
        #new_df = previous_df.merge(present_df, how="outer")
        new_df = pd.concat([previous_df, present_df], ignore_index=True)
        print(new_df)
        return new_df


def convert_to_dataframe(dynamo_db_data):
    if not dynamo_db_data:
        return pd.DataFrame()
    return pd.DataFrame(json.loads(dynamo_db_data))