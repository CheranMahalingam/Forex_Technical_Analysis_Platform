import pandas as pd
import numpy as np
import joblib


def generate_prediction(symbol, window, time_steps):
    return


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