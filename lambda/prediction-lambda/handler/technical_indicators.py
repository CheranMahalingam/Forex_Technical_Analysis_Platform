import boto3
from boto3.dynamodb.conditions import Key
import datetime

SYMBOLS = ["EURUSD", "GBPUSD"]


def generate_technical_indicators(new_ohlc_data):
    date = new_ohlc_data['Date']['S']
    timestamp = new_ohlc_data['Timestamp']['S']
    indicator_data = {}
    for symbol in SYMBOLS:
        symbolData = new_ohlc_data[symbol]
        open_price = float(symbolData['M']['Open']['N'])
        high_price = float(symbolData['M']['High']['N'])
        low_price = float(symbolData['M']['Low']['N'])
        close_price = float(symbolData['M']['Close']['N'])
        trade_volume = float(symbolData['M']['Volume']['N'])
        indicator_data[symbol + 'AccumulationDistribution'] = calculate_accumulation_distribution(
            open_price,
            high_price,
            low_price,
            close_price,
            trade_volume
        )
        previous_ema_10 = get_previous_ema(close_price, symbol, 10)
        previous_ema_50 = get_previous_ema(close_price, symbol, 50)
        indicator_data[symbol + 'Ema10'] = calculate_ema(close_price, 10, previous_ema_10)
        indicator_data[symbol  + 'Ema50'] = calculate_ema(close_price, 50, previous_ema_50)
        indicator_data[symbol] = close_price
    return indicator_data


def calculate_accumulation_distribution(open, high, low, close, volume):
    if high == low:
        return 0
    return ((2*close - low - high)/(high - low))*volume


def calculate_ema(close, periods, previous_ema):
    return close*(2/(periods + 1)) + previous_ema*(1-(2/(periods + 1)))


def get_previous_ema(close, symbol, interval):
    dynamodb = boto3.resource('dynamodb')
    table = dynamodb.Table('TechnicalAnalysisTable')

    column_name = symbol + "Ema" + str(interval)

    current_date = datetime.datetime.now()
    formatted_date = current_date.strftime("%Y-%m-%d")

    response = table.query(
        KeyConditionExpression=Key('Date').eq(formatted_date),
        ScanIndexForward=False,
        Limit=1
    )

    if response['Count'] == 0:
        return close
    else:
        return float(response['Items'][0][column_name])
