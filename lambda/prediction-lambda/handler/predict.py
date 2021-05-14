import pandas as pd
import numpy as np
import joblib
import tensorflow as tf


def generate_prediction(symbol, processed_df, window_size, time_steps, prev_close):
    """
    Loads trained model for a symbol and generates a dataframe containing twitter sentiment,
    ohlc data, and technical indicators. The model is used to predict future closing prices
    of the symbol for x next periods.

    Args:
        pair: String representing the currency pair symbol(e.g EURUSD)
        window_size: Integer representing how many periods of data are needed to make a prediction
        time_steps: Integer representing the number of periods the model predicts in advance

    Returns:
        Numpy array with closing price predictions for symbol
    """

    # Loads scalers used during training and trained model for specific symbol
    fScaler = joblib.load("scalers/{}/features.bin".format(symbol))
    scaler = joblib.load("scalers/{}/close.bin".format(symbol))
    model = tf.keras.models.load_model("models/{}".format(symbol))

    # Names of currency being bought and sold
    buy = symbol[:3]
    sell = symbol[3:]

    ema10_col = symbol + "Ema10"
    ema50_col = symbol + "Ema50"
    rsi_col = symbol + "Rsi"
    ad_col = symbol + "AccumulationDistribution"
    predicted_df = processed_df[[
        symbol,
        ema10_col,
        ema50_col,
        rsi_col,
        ad_col,
        buy,
        sell
    ]]
    predicted_df = predicted_df.copy()
    predicted_df.loc[:, [symbol]] = scaler.transform(predicted_df[[symbol]])
    predicted_df.loc[:, ~predicted_df.columns.isin([symbol])] = fScaler.transform(
        predicted_df.loc[:, ~predicted_df.columns.isin([symbol])]
    )

    X = create_dataset(predicted_df, window_size)

    y_pred = model.predict(X)

    inverse_y_pred = scaler.inverse_transform(y_pred)

    new_prediction = prev_close * np.exp(inverse_y_pred[-1])
    print(new_prediction, "PREDICTION")

    return new_prediction


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
    X = []
    print(len(df))
    if len(df) <= window_size:
        df = extend_df(df, window_size)
        print(len(df), "NEW LENGTH")
    # Ignore rows without # previous rows = window_size so there is enough past data to make prediction
    for i in range(len(df) - window_size):
        # Past data is used as an input
        v = df.iloc[i:(i + window_size)].values
        X.append(v)
    return np.array(X)


def extend_df(df, window_size):
    new_rows = window_size - len(df) + 2
    new_time = pd.date_range(start=df.index[len(df) - 1], freq='15min', periods=new_rows, closed='right')
    new_df = pd.DataFrame(index=new_time)
    extended_df = pd.concat([df, new_df]).ffill()
    print("EXTENDED")
    return extended_df
