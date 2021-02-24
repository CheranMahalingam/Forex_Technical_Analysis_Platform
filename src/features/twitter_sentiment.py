import pandas as pd
import numpy as np
from vaderSentiment.vaderSentiment import SentimentIntensityAnalyzer
from datetime import datetime
from functools import reduce
import re


class TwitterSentiment:

    currencies = {
        "usd": {"positive": ["usd/", "u.s.", "greenback", "buck", "barnie", "america", "united states"], "negative": ["/usd", "cable"]},
        "aud": {"positive": ["aud/", "gold", "aussie", "australia"], "negative": ["/aud"]},
        "gbp": {"positive": ["gbp/", "sterling", "pound", "u.k.", "united kingdom", "cable", "guppy"], "negative": ["/gbp"]},
        "nzd": {"positive": ["nzd/", "gold", "kiwi", "new zealand"], "negative": ["/nzd"]},
        "cad": {"positive": ["cad/", "oil", "loonie", "canada"], "negative": ["/cad"]},
        "chf": {"positive": ["chf/", "swiss"], "negative": ["/chf"]},
        "jpy": {"positive": ["jpy/", "asian", "japan"], "negative": ["/jpy", "guppy"]},
        "eur": {"positive": ["eur/", "fiber", "euro"], "negative": ["/eur"]}
    }

    def twitter_data_selection(self, account):
        tweets = pd.read_json(
            "../../data/raw/tweets/{}_historical.json".format(account))
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

    def generate_sentiment_score(self, tweets):
        sid = SentimentIntensityAnalyzer()
        score = []
        for post in tweets["Post"]:
            score.append(sid.polarity_scores(post)["compound"])
        tweet_score = pd.DataFrame({"Twitter_Sentiment": score})
        tweets["Twitter_Sentiment"] = tweet_score
        return tweets

    def clean_tweets(self, tweets):
        tweets = np.vectorize(remove_pattern)(tweets, "RT @[\w]*:")
        tweets = np.vectorize(remove_pattern)(tweets, "@[\w]*")
        tweets = np.vectorize(remove_pattern)(
            tweets, "https?://[A-Za-z0-9./]*")
        tweets = np.core.defchararray.replace(tweets, "[^a-zA-Z]", " ")
        tweets = np.core.defchararray.replace(tweets, "\n", " ")
        return tweets

    def remove_pattern(self, input_text, pattern):
        matches = re.findall(pattern, input_text)
        for i in matches:
            input_text = re.sub(i, "", input_text)
        return input_text

    def tweets_merge(self, tweet_list):
        if len(tweet_list) == 0:
            return
        elif len(tweet_list) == 1:
            return tweet_list[0]
        else:
            merged_tweets = reduce(lambda left, right: pd.merge(
                left, right, how="outer", on=["Time", "Post", "Twitter_Sentiment"]), tweet_list)
            merged_tweets.sort_values(by=["Time"], inplace=True)
            merged_tweets = merged_tweets.reset_index(drop=True)
            return merged_tweets

    def currency_sentiment(currencies_dict):
        country_df = pd.DataFrame()
        for currency in currencies_dict:
            for entity in currencies_dict[currency]["positive"]:
                tweet_lower = tweets["Post"].transform(
                    lambda post: post.lower())
                currency_df = tweets[tweet_lower.str.contains(entity)]
                currency_df = currency_df[{"Time", "Twitter_Sentiment"}]
                currency_df = currency_df.rename(
                    columns={"Twitter_Sentiment": currency.upper()})
                if country_df.empty:
                    country_df = currency_df
                elif not currency.upper() in country_df.columns:
                    country_df = country_df.merge(
                        currency_df, how="outer", on="Time")
                else:
                    country_df = country_df.merge(currency_df, how="outer", on=[
                                                  "Time", currency.upper()])
            for entity in currencies_dict[currency]["negative"]:
                tweet_lower = tweets["Post"].transform(
                    lambda post: post.lower())
                currency_df = tweets[tweet_lower.str.contains(entity)]
                currency_df = currency_df[{"Time", "Twitter_Sentiment"}]
                currency_df["Twitter_Sentiment"] = currency_df["Twitter_Sentiment"].transform(
                    lambda score: -score)
                currency_df = currency_df.rename(
                    columns={"Twitter_Sentiment": currency.upper()})
                if country_df.empty:
                    country_df = currency_df
                elif not currency.upper() in country_df.columns:
                    country_df = country_df.merge(
                        currency_df, how="outer", on="Time")
                else:
                    country_df = country_df.merge(currency_df, how="outer", on=[
                                                  "Time", currency.upper()])

        time_frame = pd.date_range(
            start="2018-01-01 22:00:00", freq="1T", end="2020-12-31 21:59:00")
        time_frame = pd.DataFrame(time_frame, columns=["Time"])
        time_frame["Time"] = time_frame["Time"].dt.strftime(
            "%Y-%m-%d %H:%M:%S")
        country_df = country_df.reset_index(drop=True)
        country_df = combine_tweet_dates(country_df)

        country_df = time_frame.merge(country_df, how="left", on="Time")
        country_df = country_df.fillna(0)
        country_df = country_df.sort_values(by='Time', ascending=True)

        for currency in currencies_dict:
            country_df[currency.upper()] = country_df[currency.upper()
                                                      ].rolling(1440, min_periods=1).mean()

        country_df.to_csv(
            "../../data/interim/tweets/tweets_sentiment.csv", index=False)

    def combine_tweet_dates(self, tweets):
        currencies = ["eur", "usd", "jpy", "cad", "gbp", "aud", "nzd", "chf"]
        length = 1
        for i in range(1, len(tweets.index)):
            current = tweets.at[i, "Time"]
            if current == tweets.at[i - length, "Time"] and i == len(tweets.index) - 1:
                for currency in currencies:
                    tweets.at[i - length, currency.upper()] = tweets[currency.upper()
                                                                     ].iloc[i - length: i].mean()
            elif current == tweets.at[i - length, "Time"]:
                length += 1
            elif length > 1:
                for currency in currencies:
                    tweets.at[i - length, currency.upper()] = tweets[currency.upper()
                                                                     ].iloc[i - length: i].mean()
                length = 1
        tweets.drop_duplicates(subset=["Time"], inplace=True)
        return tweets
