"""
Module generates technical indicator data such as exponential moving averages,
and accumulation/distribution lines to help LSTM model generate accurate predictions
"""

import boto3
from boto3.dynamodb.conditions import Key
import datetime

# LSTM models only supported for these currency pairs for now
SYMBOLS = ["EURUSD", "GBPUSD", "USDJPY", "AUDCAD"]


def generate_technical_indicators(new_ohlc_data):
    """
    Controller function for calculating new economic indicator data.

    Args:
        new_ohlc_data: Dictionary from DynamoDB containing latest ohlc data and
        timestamp

    Returns:
        Dictionary containing 10 day ema, 50 day ema, accumulation/distribution
        and closing price for each currency pair
    """
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
    """
    Calculates changes in accumulation/distribution line.
    A/D = ((Close - Low) - (High - Close))/(High - Low)

    Args:
        open: Float representing exchange rate at the beginning of an interval
        high: Float representing the highest exchange rate during the interval
        low: Float respresenting the lowest exchange rate during the interval
        close: Float representing the exchange rate at the end of an interval
        volume: Float representing the number of trades during the interval

    Returns:
        Float representing the change in accumulation/distribution
    """
    if high == low:
        # Prevent x/0 undefined error
        return 0
    return ((2*close - low - high)/(high - low))*volume


def calculate_ema(close, periods, previous_ema):
    """
    Calculates the exponential moving average.
    EMA = Price(t)*weighting_multipler + previous_ema*(1-weighting_multiplier)
    *weighting_multiplier is given by 2/(periods + 1)

    Args:
        close: Float representing the exchange rate at the end of an interval
        periods: Integer representing the number of days in the EMA period (commonly 12 or 26)
        previous_ema: Float representing the last calculated EMA

    Returns:
        Float representing the new EMA
    """
    return close*(2/(periods + 1)) + previous_ema*(1-(2/(periods + 1)))


def get_previous_ema(close, symbol, interval):
    """
    Searches DynamoDB for the last calculated EMA to avoid recalculating
    EMAs for multiple periods.

    Args:
        close: Float representing the exchange rate at the end of an interval
        symbol: String representing the name of the currency pair (e.g. EURUSD)
        interval: Integer representing the number of periods used in the EMA calculation
    
    Returns:
        Float representing the previous EMA
    """
    dynamodb = boto3.resource('dynamodb')
    table = dynamodb.Table('TechnicalAnalysisTable')

    # Many different EMAs were calculated based on the number of
    # periods they considered (10 day EMA and 50 day EMAs are stored)
    column_name = symbol + "Ema" + str(interval)

    current_date = datetime.datetime.now()
    formatted_date = current_date.strftime("%Y-%m-%d")

    # Searches for data from current date
    # Searches timestamps in reverse to get data quickly
    # Uses query instead of scan to speed up operation
    response = table.query(
        KeyConditionExpression=Key('Date').eq(formatted_date),
        ScanIndexForward=False,
        Limit=1
    )

    # If no data exists just use current closing price as EMA
    if response['Count'] == 0:
        return close
    else:
        return float(response['Items'][0][column_name])
