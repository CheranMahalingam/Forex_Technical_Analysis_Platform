"""Evaluates lstm model for a specific currency using mse and mae"""

from sklearn.preprocessing import MinMaxScaler
from sklearn.metrics import mean_squared_error
from sklearn.metrics import mean_absolute_error
import tensorflow as tf
import pandas as pd
import numpy as np
import joblib
import math
import os
import sys
import inspect

currentdir = os.path.dirname(os.path.abspath(
    inspect.getfile(inspect.currentframe())))
parentdir = os.path.dirname(currentdir)
sys.path.insert(0, parentdir)

from utils.model_utils import find_batch_gcd, create_split

np.random.seed(42)
tf.random.set_seed(42)


def create_dataset(df, window_size):
    X, y = [], []
    for i in range(len(df) - window_size):
        v = df.iloc[i:(i + window_size)].values
        X.append(v)
        y.append(df["Close"].iloc[i + window_size])
    return np.array(X), np.array(y)


def flatten_prediction(pred, pred_count, time_steps):
    print(pred_count, pred.shape[0])
    pred = pred[::time_steps]
    pred = pred.flatten()
    if pred_count < pred.shape[0]:
        pred = pred[:pred_count - pred.shape[0]]
    return pred


def evaluate_forecast(pred, actual):
    mse = mean_squared_error(pred, actual)
    print("Test Mean Squared Error:", mse)
    mae = mean_absolute_error(pred, actual)
    print("Test Mean Absolute Error:", mae)


def test_model(pair, window_size, batch_size, time_steps):
    fScaler = joblib.load("lstm_model/scalers/{}/features.bin".format(pair))
    scaler = joblib.load("lstm_model/scalers/{}/close.bin".format(pair))
    model = tf.keras.models.load_model("lstm_model/models/{}".format(pair))

    buy = pair[:3]
    sell = pair[3:]

    series = pd.read_csv(
        "lstm_model/data/processed/{}_processed.csv".format(pair))
    series = series[series.shape[0] % batch_size:]
    close = series[['Time', 'Real Close', 'Close']]
    close = close.copy()

    series = series.drop(['Time', 'Real Close'], axis=1)
    series = series[['Close', 'EMA_10', 'EMA_50', 'RSI', 'A/D Index',
                     '{} Interest Rate'.format(buy), '{} Interest Rate'.format(
                         sell), '{}_CPI'.format(buy), '{}_CPI'.format(sell),
                     '{} Twitter Sentiment'.format(
                         buy), '{} Twitter Sentiment'.format(sell),
                     '{} News Sentiment'.format(
                         buy), '{} News Sentiment'.format(sell),
                     #'EUR_GDP', 'USD_GDP', 'EUR Unemployment Rate', 'USD Unemployment Rate', 'EUR_PPI', 'USD_PPI'
                     ]]

    df_train, df_val, df_test = create_split(
        series, 0.75, 0.1, batch_size, window_size)
    df_test = df_test.copy()
    df_test.loc[:, ['Close']] = scaler.transform(df_test[['Close']])
    df_test.loc[:, ~df_test.columns.isin(['Close'])] = fScaler.transform(
        df_test.loc[:, ~df_test.columns.isin(['Close'])])

    X_test, y_test = create_dataset(df_test, window_size)

    y_pred = model.predict(X_test)

    multi_pred = flatten_prediction(y_pred, y_test.shape[0], time_steps)
    evaluate_forecast(multi_pred, y_test)

    df = pd.DataFrame(close[-multi_pred.shape[0] - window_size:])
    df = df[window_size:]
    df.reset_index(inplace=True, drop=True)

    index = [i for i in range(multi_pred.shape[0])]
    df_predicted = pd.DataFrame(scaler.inverse_transform(
        multi_pred.reshape(-1, 1)), columns=['Close'], index=index)
    df_actual = pd.DataFrame(scaler.inverse_transform(
        y_test.reshape(-1, 1)), columns=['Close'], index=index)

    df['Multi'] = df_predicted
    df['Multi'] = df['Multi'].shift(1)

    df['Prediction'] = np.nan

    for i in range(1, df_predicted.shape[0]):
        if (i - 1) % time_steps == 0:
            df.loc[i, 'Prediction'] = df.loc[i-1, 'Real Close'] * \
                math.exp(df.loc[i, 'Multi'])
        else:
            df.loc[i, 'Prediction'] = df.loc[i-1, 'Prediction'] * \
                math.exp(df.loc[i, 'Multi'])

    df_actual = df['Real Close'].mul(
        np.exp(df_actual['Close'].shift(-1))).shift(1)
    print(df)

    evaluate_forecast(df['Prediction'].iloc[1:], df['Real Close'].iloc[1:])

    return df['Prediction'].iloc[1:], df['Real Close'].iloc[1:]


if __name__ == "__main__":
    test_model("EURUSD", 96, 64, 4)
