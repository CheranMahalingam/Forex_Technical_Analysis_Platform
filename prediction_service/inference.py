"""Generates predictions for closing price of a symbol for 4 next periods"""

import tensorflow as tf
import pandas as pd
import numpy as np
import joblib
import random
from datetime import datetime
import time
import math
from generate_twitter_sentiment import generate_twitter_sentiment
from generate_technical_indicators import generate_technical_indicators

np.random.seed(42)
tf.random.set_seed(42)


def generate_prediction(pair, window_size, time_steps):
    """
    Loads trained model for a symbol and generates a dataframe containing twitter sentiment,
    ohlc data, and technical indicators. The model is used to predict future closing prices
    of the symbol for 4 next periods.

    Args:
        pair: String representing the currency pair symbol(e.g EURUSD)
        window_size: Integer representing how many periods of data are needed to make a prediction
        time_steps: Integer representing the number of periods the model predicts in advance

    Returns:
        Numpy array with closing price predictions for symbol
    """
    # Used temporarily to see if there are bottlenecks in performance
    startTime = time.time()

    # Loads scalers used during training and trained model for specific symbol
    fScaler = joblib.load("lstm_model/scalers/{}/features.bin".format(pair))
    scaler = joblib.load("lstm_model/scalers/{}/close.bin".format(pair))
    model = tf.keras.models.load_model("lstm_model/models/{}".format(pair))

    # Names of currency being bought and sold
    buy = pair[:3]
    sell = pair[3:]

    twitter_df = generate_twitter_sentiment(48, 60)
    ohlc_df = generate_fake_ohlc_data(twitter_df)
    technical_analysis_df = generate_technical_indicators(ohlc_df)
    inference_df = configure_time(
        15, technical_analysis_df, technical_analysis_df.loc[0, 'Time'])
    # Removes rows with NaN in a column due to calculation of EMA_50
    inference_df = inference_df.loc[5:, :]
    print(inference_df)

    predicted_df = inference_df[['Close', 'EMA_10', 'EMA_50', 'RSI', 'A/D Index',
                                 '{}'.format(buy), '{}'.format(sell),
                                 ]]
    predicted_df = predicted_df.rename(columns={'{}'.format(buy): '{} Twitter Sentiment'.format(buy),
                                                '{}'.format(sell): '{} Twitter Sentiment'.format(sell)})

    predicted_df.loc[:, ['Close']] = scaler.transform(predicted_df[['Close']])
    predicted_df.loc[:, ~predicted_df.columns.isin(['Close'])] = fScaler.transform(
        predicted_df.loc[:, ~predicted_df.columns.isin(['Close'])]
    )
    print(predicted_df)

    X, _ = create_dataset(predicted_df, window_size)

    y_pred = model.predict(X)
    print(y_pred)

    inverse_y_pred = scaler.inverse_transform(y_pred)
    print(inverse_y_pred)

    # Reverses stationary log returns from preprocessing steps
    new_prediction = inference_df['Real Close'].iloc[-1] * \
        np.exp(inverse_y_pred[-1])

    elapsedTime = time.time() - startTime
    print(elapsedTime)

    return new_prediction


def generate_fake_ohlc_data(data_df):
    """
    Generates fake ohlc data for each row of twitter sentiment dataframe. Used for testing
    without creating database or making api calls to finnhub.

    Args:
        data_df: Dataframe containing time and twitter sentiment socres for each currency

    Returns:
        Identical twitter sentiment dataframe with additional columns for fake ohlc data
    """
    data_df.loc[0, 'Close'] = random.random()*0.001 + 1
    data_df.loc[0, 'Open'] = random.random()*0.001 + 1
    data_df.loc[0, 'High'] = data_df.loc[0, 'Open'] + random.random()*0.001 if data_df.loc[0,
                                                                                           'Open'] > data_df.loc[0, 'Close'] else data_df.loc[0, 'Close'] + random.random()*0.001
    data_df.loc[0, 'Low'] = data_df.loc[0, 'Open'] - random.random()*0.001 if data_df.loc[0,
                                                                                          'Open'] < data_df.loc[0, 'Close'] else data_df.loc[0, 'Close'] - random.random()*0.001
    data_df.loc[0, 'Volume'] = random.randrange(10, 100, 1)
    for i in range(1, len(data_df)):
        data_df.loc[i, 'Open'] = data_df.loc[i - 1,
                                             'Close'] + (random.random() - 0.5)*0.001
        data_df.loc[i, 'Close'] = data_df.loc[i, 'Open'] + \
            (random.random() - 0.5)*0.001
        data_df.loc[i, 'Volume'] = random.randrange(10, 100, 1)
        if data_df.loc[i, 'Open'] > data_df.loc[i, 'Close']:
            data_df.loc[i, 'High'] = data_df.loc[i,
                                                 'Open'] + random.random()*0.001
            data_df.loc[i, 'Low'] = data_df.loc[i,
                                                'Close'] - random.random()*0.001
        else:
            data_df.loc[i, 'High'] = data_df.loc[i,
                                                 'Close'] + random.random()*0.001
            data_df.loc[i, 'Low'] = data_df.loc[i,
                                                'Open'] - random.random()*0.001
    return data_df


def configure_time(minutes, dataframe, start):
    """
    Compresses dataframe to a specific time interval (e.g. 15 minutes).

    Args:
        minutes: Integer for time interval expected in dataframe in minutes
        dataframe: Generic dataframe containing "Time" column with ISO formatted dates
        start: ISO formatted date representing "Time" row for which the new dataframe's index begins

    Returns:
        Compressed dataframe with "Time" column going up by intervals of "minutes"
    """
    time_frame = pd.date_range(
        start=start,
        freq="{}T".format(minutes),
        end=str(datetime.now())
    )
    time_frame = pd.DataFrame(time_frame, columns=["Time"])
    # Format "Time" column to match given dataframe
    time_frame["Time"] = time_frame["Time"].dt.strftime("%Y-%m-%d %H:%M:%S")

    # Should only contain rows that match time_frame and the given dataframe
    configured_df = time_frame.merge(dataframe, how="inner", on="Time")
    return configured_df


def create_dataset(df, window_size):
    """
    Builds numpy array for inputs and output column from given dataframe for time series forecasting.
    Arrays must contain past data for the model to determine the future closing price, the amount of
    past data is specified by the "window_size".

    Args:
        df: Generic dataframe with "Close" column and other input columns
        window_size: Integer representing the number of previous rows of the dataframe used to make a prediction

    Returns:
        Numpy array of model inputs with windows of past data and numpy array of expected outputs
    """
    X, y = [], []
    # Ignore rows without # previous rows = window_size so there is enough past data to make prediction
    for i in range(len(df) - window_size):
        # Past data is used as an input
        v = df.iloc[i:(i + window_size)].values
        X.append(v)
        y.append(df["Close"].iloc[i + window_size])
    return np.array(X), np.array(y)


if __name__ == "__main__":
    generate_prediction("GBPUSD", 96, 4)
