"""Utility functions used to organize data for model training and testing"""

import math


def find_batch_gcd(length, batch_size):
    while length % batch_size != 0:
        length -= 1
    return length


def create_split(df, pct_train, pct_val, batch_size, window_size):
    length = df.shape[0]
    temp_train_size = find_batch_gcd(math.floor(pct_train * length), batch_size)
    test_size = length - temp_train_size
    train_size = find_batch_gcd(math.floor((1 - pct_val) * temp_train_size), batch_size)
    val_size = temp_train_size - train_size
    df_train = df[: -val_size - test_size]
    df_val = df[-val_size - test_size - window_size : -test_size]
    df_test = df[-test_size - window_size :]
    return df_train, df_val, df_test
