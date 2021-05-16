"""Trains model for specified currency pair and saves model"""

from sklearn.preprocessing import MinMaxScaler, StandardScaler
from tensorflow.keras.models import Sequential
from tensorflow.keras.layers import Dense, Dropout, LSTM
import tensorflow as tf
import pandas as pd
import numpy as np
import joblib
import datetime
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


def create_multi_pred_dataset(df, window_size, time_steps):
    X, y = [], []
    for i in range(len(df) - window_size - time_steps - 1):
        v = df.iloc[i:(i + window_size)].values
        X.append(v)
        y.append(df["Close"].iloc[i + window_size:i +
                                  window_size + time_steps].values)
    return np.array(X), np.array(y)


def create_model(nodes, optimizer, dropout, X_train):
    model = Sequential()
    model.add(LSTM(nodes[0], input_shape=(
        X_train.shape[1], X_train.shape[2]), return_sequences=True))
    model.add(LSTM(nodes[1], return_sequences=True))
    model.add(LSTM(nodes[2]))
    model.add(Dropout(dropout))
    model.add(Dense(nodes[3]))
    model.compile(loss="mse", optimizer=optimizer, metrics=['mae'])
    return model


def train_model(pair, batch_size, window_size, nodes_arr, optimizer, dropout, epochs):
    series = pd.read_csv(
        "lstm_model/data/processed/{}_processed.csv".format(pair))

    buy = pair[:3]
    sell = pair[3:]

    series = series[series.shape[0] % batch_size:]

    series = series.drop(['Time', 'Real Close'], axis=1)
    series = series[['Close', 'EMA_10', 'EMA_50', 'RSI', 'A/D Index',
                     #'{} Interest Rate'.format(buy), '{} Interest Rate'.format(
                     #    sell), '{}_CPI'.format(buy), '{}_CPI'.format(sell),
                     '{} Twitter Sentiment'.format(
                         buy), '{} Twitter Sentiment'.format(sell),
                     #'{} News Sentiment'.format(
                     #    buy), '{} News Sentiment'.format(sell),
                     #'EUR_GDP', 'USD_GDP', 'EUR_PPI', 'USD_PPI', 'USD Unemployment Rate', 'EUR Unemployment Rate'
                     ]]

    df_train, df_val, df_test = create_split(
        series, 0.75, 0.1, batch_size, window_size)
    print(
        f'df_train.shape {df_train.shape}, df_validation.shape {df_val.shape}, df_test.shape {df_test.shape}')

    closeScaler = MinMaxScaler(feature_range=(0, 1))
    featureScaler = MinMaxScaler(feature_range=(0, 1))
    df_train = df_train.copy()
    df_val = df_val.copy()
    df_train.loc[:, ['Close']] = closeScaler.fit_transform(df_train[['Close']])
    df_train.loc[:, ~df_train.columns.isin(['Close'])] = featureScaler.fit_transform(
        df_train.loc[:, ~df_train.columns.isin(['Close'])])
    df_val.loc[:, ['Close']] = closeScaler.transform(df_val[['Close']])
    df_val.loc[:, ~df_val.columns.isin(['Close'])] = featureScaler.transform(
        df_val.loc[:, ~df_val.columns.isin(['Close'])])

    X_train, y_train = create_multi_pred_dataset(
        df_train, window_size, nodes_arr[3])
    X_val, y_val = create_multi_pred_dataset(df_val, window_size, nodes_arr[3])

    model = create_model(nodes_arr, optimizer, dropout, X_train)

    current_time = datetime.datetime.now().strftime("%Y%m%d-%H%M%S")

    log_dir = "logs/tuning/" + current_time
    tensorboard_callback = tf.keras.callbacks.TensorBoard(
        log_dir=log_dir, update_freq='epoch', profile_batch=0, histogram_freq=1)

    history = model.fit(X_train, y_train,
                        validation_data=(X_val, y_val),
                        epochs=epochs,
                        batch_size=batch_size,
                        shuffle=False,
                        callbacks=[tensorboard_callback]
                        )
    model.save("lstm_model/models/{}".format(pair))
    joblib.dump(featureScaler, "lstm_model/scalers/{}/features.bin".format(pair))
    joblib.dump(closeScaler, "lstm_model/scalers/{}/close.bin".format(pair))
    del model

    return history


if __name__ == "__main__":
    batch_size = 64
    window_size = 96
    nodes = [80, 64, 32, 4]
    optimizer = tf.keras.optimizers.Adam(learning_rate=0.0005)
    dropout = 0.2
    epochs = 15
    train_model("AUDCAD", batch_size, window_size,
                nodes, optimizer, dropout, epochs)
