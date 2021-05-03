import pandas as pd
import numpy as np

CURRENCIES = ["USD", "EUR", "CAD", "NZD", "CHF", "AUD", "JPY", "GBP"]
SYMBOLS = ["EURUSD", "GBPUSD"]


def preprocess(analysis_df):
    processed_df = combine_date_timestamp(analysis_df)
    analysis_df = convert_to_nan(analysis_df, CURRENCIES, -100)
    analysis_df = convert_to_nan(analysis_df, SYMBOLS, 0)

    for symbol in SYMBOLS:
        analysis_df[symbol] = analysis_df[symbol].fillna(method="ffill")
        processed_df[symbol] = stationary_log_returns(analysis_df, symbol)
        ema10_col = symbol + "Ema10"
        ema50_col = symbol + "Ema50"
        processed_df[ema10_col] = stationary_log_returns(analysis_df, ema10_col)
        processed_df[ema50_col] = stationary_log_returns(analysis_df, ema50_col)
        ad_col = symbol + "AccumulationDistribution"
        processed_df[ad_col] = analysis_df[ad_col]
        rsi_col = symbol + "Rsi"
        processed_df[rsi_col] = calculate_rsi(analysis_df, symbol, 14)

    for currency in CURRENCIES:
        processed_df[currency] = average_tweet_sentiment(analysis_df, currency, 60)
        print(processed_df[currency])
    
    processed_df = compress_dataframe_time_interval(processed_df, 15)

    print(processed_df)
    return processed_df


def stationary_log_returns(analysis_df, column):
    """
    Calculates log returns for EMA and closing price to make data stationary.

    Args:
        analysis_df: Dataframe containing OHLC data, Time, and technical indicators
        column: Name of column for which log returns must be calculated

    Returns:
        Dataframe with column substituted with log returns
    """
    analysis_df[column] = np.log(analysis_df[column] / analysis_df[column].shift(1))
    return analysis_df[column]


def combine_date_timestamp(analysis_df):
    new_df = pd.DataFrame()
    new_df['Time'] = analysis_df['Date'] + analysis_df['Timestamp']
    new_df['Time'] = pd.to_datetime(new_df['Time'], format="%Y-%m-%d%H:%M:%S")
    return new_df


def convert_to_nan(analysis_df, column_names, value):
    for name in column_names:
        analysis_df[name] = analysis_df[name].transform(lambda val: np.nan if val == value else val)
    return analysis_df


def calculate_rsi(analysis_df, column, window):
    delta = analysis_df[column]
    up_periods = delta.copy()
    up_periods[delta<=0] = 0.0
    down_periods = abs(delta.copy())
    down_periods[delta>0] = 0.0
    rs_up = up_periods.rolling(window).mean()
    rs_down = down_periods.rolling(window).mean()
    rsi = 100 - 100/(1+rs_up/rs_down)
    return rsi


def average_tweet_sentiment(analysis_df, column, window):
    return analysis_df[column].rolling(window, min_periods=1).mean()


def compress_dataframe_time_interval(processed_df, interval):
    cool = processed_df.resample('{}min'.format(interval), on='Time').mean()
    print(cool)
    return cool
