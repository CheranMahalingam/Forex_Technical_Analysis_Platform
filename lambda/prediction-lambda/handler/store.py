import boto3
from decimal import Decimal

DB_ANALYSIS_COLUMNS = [
    "EURUSD",
    "GBPUSD",
    "USDJPY",
    "AUDCAD",
    "EURUSDEma10",
    "EURUSDEma50",
    "GBPUSDEma10",
    "GBPUSDEma50",
    "USDJPYEma10",
    "USDJPYEma50",
    "AUDCADEma10",
    "AUDCADEma50",
    "EURUSDAccumulationDistribution",
    "GBPUSDAccumulationDistribution",
    "USDJPYAccumulationDistribution",
    "AUDCADAccumulationDistribution"
]


def store_indicator_data(date, timestamp, analysis, sentiment):
    dynamodb = boto3.resource('dynamodb')
    table = dynamodb.Table('TechnicalAnalysisTable')

    item = {
        'Date': date,
        'Timestamp': timestamp,
    }

    for name in DB_ANALYSIS_COLUMNS:
        item[name] = Decimal(str(analysis[name]))
    
    for currency in sentiment:
        item[currency] = Decimal(str(sentiment[currency]))

    table.put_item(Item=item)
