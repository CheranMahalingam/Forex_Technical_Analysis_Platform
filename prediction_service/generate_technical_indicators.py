"""Generates technical analysis data used to predict future closing prices"""

from talib import abstract
import pandas as pd
import numpy as np


def generate_technical_indicators(pair_df):
    """
    Builds dataframe containing Time and indicators such as 10-day and 50-day EMAs,
    RSI, and A/D index. Also uses differencing to remove seasonality trends.

    Args:
        pair_df: Dataframe containing "Time" column and columns for ohlc data

    Returns:
        Dataframe containing technical indicator data for a specific currency pair
    """
    pair_df["EMA_10"] = pd.DataFrame(abstract.EMA(
        pair_df["Close"], timeperiod=10))  # 960))
    pair_df["EMA_50"] = pd.DataFrame(abstract.EMA(
        pair_df["Close"], timeperiod=50))  # 4800))
    pair_df["RSI"] = pd.DataFrame(
        abstract.RSI(pair_df["Close"], timeperiod=14))
    pair_df["A/D Index"] = pd.DataFrame(
        abstract.AD(
            pair_df["High"], pair_df["Low"], pair_df["Close"], pair_df["Volume"]
        )
    )
    # Take difference to remove trends
    pair_df["A/D Index"] = pair_df["A/D Index"] - pair_df["A/D Index"].shift(1)
    pair_df = stationary_log_returns(pair_df)
    return pair_df


def stationary_log_returns(pair_df):
    """
    Calculates log returns for EMA and closing price to make data stationary.

    Args:
        pair_df: Dataframe containing OHLC data, Time, and technical indicators

    Returns:
        Dataframe with EMA and closing prices substituted with log returns
    """
    pair_df = pair_df.copy()
    # Create new column to remember real closing price for symbol or else log returns are non-reversible
    pair_df["Real Close"] = pair_df["Close"]
    pair_df["Close"] = np.log(pair_df["Close"] / pair_df["Close"].shift(1))
    pair_df["EMA_10"] = np.log(pair_df["EMA_10"] / pair_df["EMA_10"].shift(1))
    pair_df["EMA_50"] = np.log(pair_df["EMA_50"] / pair_df["EMA_50"].shift(1))
    return pair_df
