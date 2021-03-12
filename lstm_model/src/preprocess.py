"""Generates processed data for a specified currency pair"""

from datetime import datetime
from functools import reduce
import pandas as pd
import numpy as np
from data import cpi_indicator, gdp_indicator, interest_rate_indicator
from features import exchange_rate_indicator, news_sentiment, twitter_sentiment


def build_currency_dataframe(buy, sell):
    pair = (buy + sell).upper()
    buy_df = combine_indicators(buy, pair)
    buy = buy.upper()
    buy_df = buy_df.rename(columns={
        "CPI": buy + "_CPI",
        "GDP": buy + "_GDP",
        "Interest Rate": buy + " Interest Rate",
        "PPI": buy + "_PPI",
        "Unemployment Rate": buy + " Unemployment Rate",
        "News Sentiment": buy + " News Sentiment",
        "Twitter Sentiment": buy + " Twitter Sentiment",
    })
    buy_df = buy_df.reset_index(drop=True)
    sell_df = combine_indicators(sell, pair)
    sell_df = sell_df[{"Time", "CPI", "GDP", "Interest Rate", "PPI", "Unemployment Rate", "News Sentiment", "Twitter Sentiment"}]
    sell = sell.upper()
    sell_df = sell_df.rename(columns={
        "CPI": sell + "_CPI", 
        "GDP": sell + "_GDP", 
        "Interest Rate": sell + " Interest Rate",
        "PPI": sell + "_PPI",
        "Unemployment Rate": sell + " Unemployment Rate",
        "News Sentiment": sell + " News Sentiment",
        "Twitter Sentiment": sell + " Twitter Sentiment",
    })
    sell_df.reset_index(drop=True)
    pair_df = buy_df.merge(sell_df, how="inner", on=["Time"])
    pair_df['Time'] = pd.to_datetime(pair_df['Time'], utc=True)
    pair_df = configure_time(15, pair_df)
    pair_df.to_csv("lstm_model/data/processed/{}_processed.csv".format(pair), index=False)
    return pair_df

def combine_indicators(currency, pair):
    cpi = pd.read_csv("lstm_model/data/interim/cpi/{}_cpi_processed.csv".format(currency))
    gdp = pd.read_csv("lstm_model/data/interim/gdp/{}_gdp_processed.csv".format(currency))
    ir = pd.read_csv("lstm_model/data/interim/interest_rate/{}_ir_processed.csv".format(currency))
    ppi = pd.read_csv("lstm_model/data/interim/ppi/{}_ppi_processed.csv".format(currency))
    ue = pd.read_csv("lstm_model/data/interim/unemployment_rate/{}_ue_processed.csv".format(currency))
    news = pd.read_csv("lstm_model/data/interim/news/news_sentiment.csv")
    news = news[{"Time", currency.upper()}]
    news = news.rename(columns={currency.upper(): "News Sentiment"})
    tweets = pd.read_csv("lstm_model/data/interim/tweets/tweets_sentiment.csv")
    tweets = tweets[{"Time", currency.upper()}]
    tweets = tweets.rename(columns={currency.upper(): "Twitter Sentiment"})

    combined_df = merge_dataframe([cpi, gdp, ir, ppi, ue, news, tweets])

    if currency.upper() in pair:
        exchange_rate = pd.read_csv("lstm_model/data/interim/exchange_rate/{}_exchange.csv".format(pair))
        combined_df = merge_dataframe([combined_df, exchange_rate])

    combined_df = combined_df[combined_df["RSI"].notnull()]
    return combined_df

def merge_dataframe(data_list):
    merged_data = reduce(
        lambda left, right: pd.merge(left, right, how="outer", on="Time"), data_list)
    merged_data.sort_values(by=["Time"], inplace=True)
    merged_data = merged_data.reset_index(drop=True)
    return merged_data

def configure_time(minutes, dataframe):
    time_frame = pd.date_range(start="2018-01-01 22:00:00", freq="{}T".format(minutes), end="2020-12-31 21:59:00")
    time_frame = pd.DataFrame(time_frame, columns=["Time"])
    time_frame["Time"] = time_frame["Time"].dt.strftime("%Y-%m-%d %H:%M:%S")
    time_frame['Time'] = pd.to_datetime(time_frame['Time'], utc=True)
    
    configured_df = time_frame.merge(dataframe, how="inner", on="Time")
    return configured_df

def convert_date(exchange):
    exchange["Time"] = pd.to_datetime(exchange["Time"], format="%Y-%m-%d %H:%M:%S")
    return exchange


if __name__ == "__main__":
    currency_pair = "AUDCAD"
    sentiment_keyword = {
        "usd": {
            "positive": ["usd/", "u.s.", "greenback", "buck", "barnie", "america", "united states"],
            "negative": ["/usd", "cable"]
        },
        "aud": {
            "positive": ["aud/", "gold", "aussie", "australia"],
            "negative": ["/aud"]
        },
        "gbp": {
            "positive": ["gbp/", "sterling", "pound", "u.k.", "united kingdom", "cable", "guppy"],
            "negative": ["/gbp"]
        },
        "nzd": {
            "positive": ["nzd/", "gold", "kiwi", "new zealand"],
            "negative": ["/nzd"]
        },
        "cad": {
            "positive": ["cad/", "oil", "loonie", "canada"],
            "negative": ["/cad"]
        },
        "chf": {
            "positive": ["chf/", "swiss"],
            "negative": ["/chf"]
        },
        "jpy": {
            "positive": ["jpy/", "asian", "japan"],
            "negative": ["/jpy", "guppy"]
        },
        "eur": {
            "positive": ["eur/", "fiber", "euro"],
            "negative": ["/eur"]
        }
    }

    cpi_indicator.cpi_preprocess()
    print("CPI indicator processed...")

    gdp_indicator.gdp_preprocess()
    print("GDP indicator processed...")

    interest_rate_indicator.interest_rate_preprocess()
    print("Interest Rate indicator processed...")

    exchange_rate_indicator.indicators_preprocess(currency_pair)
    print("Exchange Rate technical indicators generated...")

    reuters = news_sentiment.news_unpacking("reuters")
    daily_fx = news_sentiment.news_unpacking("dailyfx")
    forex_live = news_sentiment.news_unpacking("forexlive")
    news_df = news_sentiment.news_merge([reuters, daily_fx, forex_live])
    news_df = news_sentiment.generate_sentiment_score(news_df)
    news_sentiment.currency_sentiment(sentiment_keyword, news_df)
    print("News sentiment analysis generated...")

    forex_com = twitter_sentiment.twitter_unpacking("forexcom")
    ft_markets = twitter_sentiment.twitter_unpacking("FTMarkets")
    bloomberg = twitter_sentiment.twitter_unpacking("markets")
    reuters = twitter_sentiment.twitter_unpacking("ReutersGMF")
    wsj = twitter_sentiment.twitter_unpacking("WSJmarkets")
    fx_street_1 = twitter_sentiment.twitter_unpacking("FXstreetNews")
    fx_street_2 = twitter_sentiment.twitter_unpacking("FXstreetNews2")
    tweets_df = twitter_sentiment.tweets_merge([forex_com, ft_markets, bloomberg, reuters, wsj, fx_street_1, fx_street_2])
    twitter_sentiment.currency_sentiment(sentiment_keyword, tweets_df)
    print("Twitter sentiment analysis generated...")

    build_currency_dataframe(currency_pair[:3].lower(), currency_pair[3:].lower())
    print("Feature engineering & Data cleaning complete!")
    print("{}_processed.csv generated!".format(currency_pair))
