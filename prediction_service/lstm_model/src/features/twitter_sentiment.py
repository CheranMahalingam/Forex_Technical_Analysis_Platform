"""Feature engineering for sentiment score of a currency using tweets."""

from functools import reduce
import re
from vaderSentiment.vaderSentiment import SentimentIntensityAnalyzer
import pandas as pd
import numpy as np


def twitter_unpacking(account):
    """
    Converts json file containing an account's tweets to a dataframe.

    Args:
        account: String representing the name of a twitter account

    Returns:
        Dataframe containing tweets from an account with post and date.
    """
    tweets = pd.read_json(
        "lstm_model/data/raw/tweets/{}_historical.json".format(account)
    )
    headline_arr = []
    date_arr = []

    for post in range(tweets["headline"].count()):
        headline_arr.append(tweets["headline"].iloc[post])
        date_arr.append(tweets["date"].iloc[post])
    tweets_df = pd.DataFrame({"Time": date_arr, "Post": headline_arr})

    tweets_df["Post"] = clean_tweets(tweets_df["Post"])
    tweets_df = generate_sentiment_score(tweets_df)
    tweets_df["Time"] = tweets_df["Time"].dt.strftime("%Y-%m-%d %H:%M:00")
    return tweets_df


def generate_sentiment_score(tweets):
    """
    Generates sentiment scores for each tweet.

    Args:
        tweets: Dataframe containing tweet date and post

    Returns:
        Tweet dataframe with additional column for sentiment score.
    """
    sid = SentimentIntensityAnalyzer()
    score = []
    for post in tweets["Post"]:
        score.append(sid.polarity_scores(post)["compound"])
    tweet_score = pd.DataFrame({"Twitter_Sentiment": score})
    tweets["Twitter_Sentiment"] = tweet_score
    return tweets


def clean_tweets(tweets):
    """
    Strips twitter posts of uninformative special characters.

    Args:
        tweets: Dataframe containing tweet date, post, and sentiment score

    Returns:
        Tweet dataframe with machine readable tweets.
    """
    # pylint: disable=anomalous-backslash-in-string
    tweets = np.vectorize(remove_pattern)(tweets, "RT @[\w]*:")
    tweets = np.vectorize(remove_pattern)(tweets, "@[\w]*")
    tweets = np.vectorize(remove_pattern)(tweets, "https?://[A-Za-z0-9./]*")
    tweets = np.core.defchararray.replace(tweets, "[^a-zA-Z]", " ")
    tweets = np.core.defchararray.replace(tweets, "\n", " ")
    return tweets


def remove_pattern(input_text, pattern):
    """
    Finds patterns in posts and substitutes them with blank space.

    Args:
        input_text: String representing a twitter post
        pattern: Regex pattern to search for in twitter post

    Returns:
        String with pattern stripped.
    """
    match = re.findall(pattern, input_text)
    for i in match:
        input_text = re.sub(i, "", input_text)
    return input_text


def tweets_merge(tweet_list):
    """
    Merges a list of tweet dataframes.

    Args:
        tweet_list: Array of tweet dataframes containing time, post, and sentiment scores

    Returns:
        Single dataframe containing tweets + scores for all accounts.
    """
    if len(tweet_list) == 0:
        return None
    if len(tweet_list) == 1:
        return tweet_list[0]
    merged_tweets = reduce(
        lambda left, right: pd.merge(
            left, right, how="outer", on=["Time", "Post", "Twitter_Sentiment"]
        ),
        tweet_list,
    )
    merged_tweets.sort_values(by=["Time"], inplace=True)
    merged_tweets = merged_tweets.reset_index(drop=True)
    return merged_tweets


def currency_sentiment(currencies_dict, tweets):
    """
    Generates csv file containing tweet sentiment scores for each country according to date.

    Args:
        currencies_dict: Dictionary holding keywords to indicate that a post affects the
                         strength of a currency and whether is positive or negative
        tweets: Dataframe containing twitter posts with their sentiment scores
    """
    country_df = pd.DataFrame()
    for currency in currencies_dict:
        for entity in currencies_dict[currency]["positive"]:
            tweet_lower = tweets["Post"].transform(lambda post: post.lower())
            currency_df = tweets[tweet_lower.str.contains(entity)]
            currency_df = currency_df[{"Time", "Twitter_Sentiment"}]
            currency_df = currency_df.rename(
                columns={"Twitter_Sentiment": currency.upper()}
            )
            if country_df.empty:
                country_df = currency_df
            elif not currency.upper() in country_df.columns:
                country_df = country_df.merge(currency_df, how="outer", on="Time")
            else:
                country_df = country_df.merge(
                    currency_df, how="outer", on=["Time", currency.upper()]
                )
        for entity in currencies_dict[currency]["negative"]:
            tweet_lower = tweets["Post"].transform(lambda post: post.lower())
            currency_df = tweets[tweet_lower.str.contains(entity)]
            currency_df = currency_df[{"Time", "Twitter_Sentiment"}]
            currency_df["Twitter_Sentiment"] = currency_df[
                "Twitter_Sentiment"
            ].transform(lambda score: -score)
            currency_df = currency_df.rename(
                columns={"Twitter_Sentiment": currency.upper()}
            )
            if country_df.empty:
                country_df = currency_df
            elif not currency.upper() in country_df.columns:
                country_df = country_df.merge(currency_df, how="outer", on="Time")
            else:
                country_df = country_df.merge(
                    currency_df, how="outer", on=["Time", currency.upper()]
                )

    time_frame = pd.date_range(
        start="2018-01-01 22:00:00", freq="1T", end="2020-12-31 21:59:00"
    )
    time_frame = pd.DataFrame(time_frame, columns=["Time"])
    time_frame["Time"] = time_frame["Time"].dt.strftime("%Y-%m-%d %H:%M:%S")
    country_df = country_df.reset_index(drop=True)
    country_df = combine_dates(country_df)

    country_df = time_frame.merge(country_df, how="left", on="Time")
    country_df = country_df.sort_values(by="Time", ascending=True)

    for currency in currencies_dict:
        country_df[currency.upper()] = (
            country_df[currency.upper()].rolling(1440, min_periods=1).mean()
        )
    country_df = country_df.fillna(0)

    country_df.to_csv(
        "lstm_model/data/interim/tweets/tweets_sentiment.csv", index=False
    )


def combine_dates(tweets):
    """
    Merge sentiment scores according to date.

    Args:
        tweets: Dataframe containing countries and their sentiment scores at a certain time

    Returns:
        Dataframe with a country's sentiment score with sequential time.
    """
    currencies = ["eur", "usd", "jpy", "cad", "gbp", "aud", "nzd", "chf"]
    length = 1
    for i in range(1, len(tweets.index)):
        current = tweets.at[i, "Time"]
        if current == tweets.at[i - length, "Time"] and i == len(tweets.index) - 1:
            for currency in currencies:
                tweets.at[i - length, currency.upper()] = (
                    tweets[currency.upper()].iloc[i - length : i].mean()
                )
        elif current == tweets.at[i - length, "Time"]:
            length += 1
        elif length > 1:
            for currency in currencies:
                tweets.at[i - length, currency.upper()] = (
                    tweets[currency.upper()].iloc[i - length : i].mean()
                )
            length = 1
    tweets.drop_duplicates(subset=["Time"], inplace=True)
    return tweets
