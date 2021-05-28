"""Module for preprocessing data"""

import pandas as pd
import numpy as np

CURRENCIES = ["USD", "EUR", "CAD", "NZD", "CHF", "AUD", "JPY", "GBP"]
SYMBOLS = ["EURUSD", "GBPUSD", "USDJPY", "AUDCAD"]


def preprocess(analysis_df):
    """
    Controller for running all preprocessing methods in the correct sequence.

    Args:
        analysis_df: Pandas dataframe containing symbol exchange rates, ema,
        accumulation/distribution, and twitter sentiment scores for each symbol
    
    Returns:
        Pandas dataframe with preprocessed data
    """
    # Combines Date and Timestamp columns
    processed_df = combine_date_timestamp(analysis_df)
    # Converts empty data to nan to avoid issues when calculating mean
    analysis_df = convert_to_nan(analysis_df, CURRENCIES, -100)
    analysis_df = convert_to_nan(analysis_df, SYMBOLS, 0)

    for symbol in SYMBOLS:
        analysis_df[symbol] = analysis_df[symbol].fillna(method="ffill")
        # Calculates log returns for ema and closing price
        # Removes trends and makes data stationary to improve performance of LSTM
        processed_df[symbol] = stationary_log_returns(analysis_df, symbol)
        ema10_col = symbol + "Ema10"
        ema50_col = symbol + "Ema50"
        processed_df[ema10_col] = stationary_log_returns(analysis_df, ema10_col)
        processed_df[ema50_col] = stationary_log_returns(analysis_df, ema50_col)
        ad_col = symbol + "AccumulationDistribution"
        processed_df[ad_col] = analysis_df[ad_col]
        rsi_col = symbol + "Rsi"
        # Calculates 14 period RSI
        processed_df[rsi_col] = calculate_rsi(analysis_df, symbol, 14)

    for currency in CURRENCIES:
        # Gets average tweet sentiment using past hour of scores
        processed_df[currency] = average_tweet_sentiment(analysis_df, currency, 60)

    # Compresses all data from 1 minute intervals to 15 minute intervals
    processed_df = compress_dataframe_time_interval(processed_df, 15)
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
    analysis_df[column] = analysis_df[column].fillna(method="bfill")
    return analysis_df[column]


def combine_date_timestamp(analysis_df):
    """
    Creates a dataframe with time in the format YYYY-mm-ddHH:MM:SS.

    Args:
        analysis_df: Pandas dataframe containing Date and Timestamp columns

    Returns:
        Pandas dataframe with a Time column
    """
    new_df = pd.DataFrame()
    new_df['Time'] = analysis_df['Date'] + analysis_df['Timestamp']
    new_df['Time'] = pd.to_datetime(new_df['Time'], format="%Y-%m-%d%H:%M:%S")
    return new_df


def convert_to_nan(analysis_df, column_names, value):
    """
    Allows a certain value along a dataframe column to be converted to nan.

    Args:
        analysis_df: Pandas dataframe
        column_names: List of strings which are names of columns in the dataframe
        value: Any value that needs to be substituted with nan

    Returns:
        Identical pandas dataframe with specific value substituted with nan
    """
    for name in column_names:
        analysis_df[name] = analysis_df[name].transform(lambda val: np.nan if val == value else val)
    return analysis_df


def calculate_rsi(analysis_df, column, window):
    """
    Calculates relative stength index.

    Args:
        analysis_df: Pandas dataframe with a closing price column
        column: String representing the name of the closing price column
        window: Integer representing the number of periods used in the RSI calculation

    Returns:
        Pandas dataframe containing RSI
    """
    delta = analysis_df[column]
    up_periods = delta.copy()
    up_periods[delta<=0] = 0.0
    down_periods = abs(delta.copy())
    down_periods[delta>0] = 0.0
    rs_up = up_periods.rolling(window, min_periods=1).mean()
    rs_down = down_periods.rolling(window, min_periods=1).mean()
    rsi = 100 - 100/(1+rs_up/rs_down)
    # Impute nan rows
    rsi = rsi.fillna(method="bfill")
    return rsi


def average_tweet_sentiment(analysis_df, column, window):
    """
    Calculates a rolling mean for tweet sentiment scores.

    Args:
        analysis_df: Pandas dataframe
        column: String representing the name of the Pandas column with sentiment
        scores
        window: Integer representing the size of the rolling window

    Returns:
        Pandas dataframe column containing rolling mean of twitter sentiment
    """
    average_sentiment = analysis_df[column].rolling(window, min_periods=1).mean()
    average_sentiment = average_sentiment.replace(np.nan).fillna(0)
    return average_sentiment


def compress_dataframe_time_interval(processed_df, interval):
    """
    Resamples dataframe according to time interval. If data is originally in 1
    minute intervals the number of rows can be reduced by making the interval 15 minutes.
    To maintain data quality, an average is taken when compressing the dataframe.

    Args:
        processed_df: Pandas dataframe containing a "Time" column with date ranges
        interval: Integer representing the new date range interval for the compressed dataframe
    
    Returns:
        Pandas dataframe with compressed time interval
    """
    resampled_df = processed_df.resample('{}min'.format(interval), on='Time').mean()
    return resampled_df
