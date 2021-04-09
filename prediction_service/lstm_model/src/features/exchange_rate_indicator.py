"""Feature engineering for currency pair exchange rates"""

from talib import abstract
import pandas as pd
import numpy as np


def indicators_preprocess(pair):
    """
    Reads 1 minute interval exchange rate data and calculates 10 day EMA, 50 day EMA,
    RSI, and A/D index. EMA and Close prices are transformed by using log returns
    to make the data stationary. Log returns are calculated by,

    Returns = ln(Closing price / Previous closing price)

    Args:
        pair: String representing the currency pair symbol(e.g EURUSD)
    """
    currency = pd.read_csv(
        "lstm_model/data/external/exchange_rates/{}_M1.csv".format(pair)
    )
    currency = convert_date(currency)

    currency = configure_time(15, currency)

    # pylint: disable=no-member
    currency["EMA_10"] = pd.DataFrame(abstract.EMA(currency["Close"], timeperiod=960))
    currency["EMA_50"] = pd.DataFrame(abstract.EMA(currency["Close"], timeperiod=4800))
    currency["RSI"] = pd.DataFrame(abstract.RSI(currency["Close"], timeperiod=14))
    currency["A/D Index"] = pd.DataFrame(
        abstract.AD(
            currency["High"], currency["Low"], currency["Close"], currency["Volume"]
        )
    )
    currency["A/D Index"] = currency["A/D Index"] - currency["A/D Index"].shift(1)
    currency = stationary_log_returns(currency)

    currency["Time"] = currency["Time"].dt.strftime("%Y-%m-%d %H:%M:%S")

    currency.to_csv(
        "lstm_model/data/interim/exchange_rate/{}_exchange.csv".format(pair),
        index=False,
    )


def convert_date(exchange):
    """
    Converts dates of raw exchange rate data to pandas datetime objects.

    Args:
        exchange: Dataframe for currency pair containing date and OHLC data

    Returns:
        Dataframe for currency pair with datetime objects in Time column
    """
    exchange = exchange.rename(columns={"DateTime": "Time"})
    exchange["Time"] = pd.to_datetime(exchange["Time"]).dt.strftime("%Y-%m-%d %H:%M:%S")
    return exchange


def stationary_log_returns(pair_df):
    """
    Calculates log returns for EMA and closing price to make data stationary.

    Args:
        pair_df: Dataframe containing OHLC data, Time, and technical indicators

    Returns:
        Dataframe with EMA and closing prices substituted with log returns
    """
    pair_df = pair_df.copy()
    pair_df["Real Close"] = pair_df["Close"]
    pair_df["Close"] = np.log(pair_df["Close"] / pair_df["Close"].shift(1))
    pair_df["EMA_10"] = np.log(pair_df["EMA_10"] / pair_df["EMA_10"].shift(1))
    pair_df["EMA_50"] = np.log(pair_df["EMA_50"] / pair_df["EMA_50"].shift(1))
    return pair_df


def configure_time(minutes, pair_df):
    """
    Merges dataframe with data in 1 minute intervals to x minute intervals.

    Args:
        minutes: Integer representing desired time interval for training lstm
        pair_df: Dataframe with technical indicators on 1 minute intervals

    Returns:
        Dataframe with technical indicators on x minute intervals.
    """
    time_frame = pd.date_range(
        start="2017-01-01 22:00:00",
        freq="{}T".format(minutes),
        end="2020-12-31 21:59:00",
    )
    time_frame = pd.DataFrame(time_frame, columns=["Time"])
    time_frame["Time"] = time_frame["Time"].dt.strftime("%Y-%m-%d %H:%M:%S")
    time_frame["Time"] = pd.to_datetime(time_frame["Time"], utc=True)
    pair_df["Time"] = pair_df["Time"].astype("datetime64[ns, UTC]")
    configured_df = time_frame.merge(pair_df, how="inner", on="Time")
    return configured_df
