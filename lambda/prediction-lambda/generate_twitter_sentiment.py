"""Generates dataframe containing sentiment scores for individual currencies"""

from TwitterAPI import TwitterAPI, TwitterPager
import datetime
from dotenv import load_dotenv
import re
import pandas as pd
from vaderSentiment.vaderSentiment import SentimentIntensityAnalyzer
import os
from sentiment_keyword_defs import SENTIMENT_KEYWORDS
load_dotenv()

api = TwitterAPI(consumer_key=os.getenv("TWITTER_CONSUMER_KEY"), consumer_secret=os.getenv("TWITTER_CONSUMER_SECRET"),
                 access_token_key=os.getenv("TWITTER_ACCESS_TOKEN_KEY"), access_token_secret=os.getenv("TWITTER_ACCESS_TOKEN_SECRET"), api_version='2')


def generate_twitter_sentiment(hours, window):
    """
    Generates a dataframe containing twitter sentiment scores for each currency over a time interval.

    Args:
        hours: Integer representing number of hours in the past that tweets should be collected
        window: Integer of number of minutes over which to calculate an average score for each currency

    Returns:
        Dataframe with "Time" column and sentiment scores for all major currencies
    """
    # Twitter API only accepts dates in format "YYYY-MM-DDTHH:MM:SS.%fZ"
    tweets = get_twitter_data(
        (datetime.datetime.now() - datetime.timedelta(hours=hours)).isoformat("T") + "Z")
    print(tweets, "TWEETS")
    sentiment = tweet_sentiment(tweets)
    print(sentiment, "SENTIMENT")
    return country_sentiment_df(sentiment, hours, window)


def get_twitter_data(start_time):
    """
    Collects tweets from prominent forex accounts over specified interval.

    Args:
        start_time: String of RFC33339 formatted date

    Returns:
        List with dictionaries containing tweet text, when the were created, and public metrics
    """
    # Get tweets in batches of 100 for speed
    # 5 second delay between pages to prevent rate limiting
    pager = TwitterPager(api, 'tweets/search/recent', {
        'query': 'from:FXstreetNews OR from:forexcom',
        'tweet.fields': 'public_metrics,created_at',
        'start_time': str(start_time),
        'max_results': 100
    }
    )
    tweet_data = []
    counter = 0
    for item in pager.get_iterator(new_tweets=False):
        tweet_data.append(
            {"text": item['text'], "created_at": item['created_at']})
        print(item)
        counter += 1
    print(counter)
    return tweet_data


def country_sentiment_df(tweets, start, window):
    """
    Uses sentiment scores from each tweets and maps them to average scores for each currency.
    Certain keywords are used to link the tweet to a currency as positive or negative. The
    sentiment scores are averaged over intervals to minimize outliers in score.

    Args:
        tweets: List of dictionaries containing tweets and their sentiment scores
        start: String of RFC3339 formatted date
        window: Integer representing the number of rows used to calculate a rolling average for sentiment scores

    Returns:
        Dataframe containing "Time" and sentiment scores for each currency around a certain time 
    """
    tweet_df = pd.DataFrame()
    # Reformats dates to ISO
    tweet_df['Time'] = [datetime.datetime.strptime(
        tweet['created_at'], "%Y-%m-%dT%H:%M:%S.%fZ") for tweet in tweets]
    tweet_df['Time'] = tweet_df['Time'].dt.strftime("%Y-%m-%d %H:%M:00")
    tweet_df['Twitter_Sentiment'] = [tweet['score'] for tweet in tweets]
    # Lowercase posts makes it easier to find keyword matches
    tweet_df['Post'] = [tweet['text'].lower() for tweet in tweets]

    # The seconds digits are insignificant
    start = str(datetime.datetime.strftime(datetime.datetime.now() -
                                           datetime.timedelta(hours=start), "%Y-%m-%d %H:%M:00"))

    country_df = pd.DataFrame()
    for currency in SENTIMENT_KEYWORDS:
        for entity in SENTIMENT_KEYWORDS[currency]["positive"]:
            # Creates mask of rows which contain the keyword
            currency_df = tweet_df[tweet_df['Post'].str.contains(entity)]

            currency_df = currency_df[{"Time", "Twitter_Sentiment"}]
            currency_df = currency_df.rename(
                columns={"Twitter_Sentiment": currency.upper()}
            )
            if country_df.empty:
                country_df = currency_df
            elif not currency.upper() in country_df.columns:
                country_df = country_df.merge(
                    currency_df, how="outer", on="Time")
            else:
                country_df = country_df.merge(
                    currency_df, how="outer", on=["Time", currency.upper()]
                )
        for entity in SENTIMENT_KEYWORDS[currency]['negative']:
            # Creates mask of rows which contain the keyword
            currency_df = tweet_df[tweet_df['Post'].str.contains(entity)]

            currency_df = currency_df[{"Time", "Twitter_Sentiment"}]
            # The score sign is switched since the keyword is inversely related to the currency
            currency_df["Twitter_Sentiment"] = currency_df[
                "Twitter_Sentiment"
            ].transform(lambda score: -score)
            currency_df = currency_df.rename(
                columns={"Twitter_Sentiment": currency.upper()}
            )
            if country_df.empty:
                country_df = currency_df
            elif not currency.upper() in country_df.columns:
                country_df = country_df.merge(
                    currency_df, how="outer", on="Time")
            else:
                country_df = country_df.merge(
                    currency_df, how="outer", on=["Time", currency.upper()]
                )

    print(country_df)
    time_frame = pd.date_range(
        start=start, freq="1T", end=str(datetime.datetime.now())
    )
    time_frame = pd.DataFrame(time_frame, columns=["Time"])
    time_frame["Time"] = time_frame["Time"].dt.strftime("%Y-%m-%d %H:%M:%S")
    country_df = country_df.reset_index(drop=True)
    # Sentiment scores are integrated with time column which goes up by 1 minute intervals
    country_df = combine_dates(country_df)

    country_df = time_frame.merge(country_df, how="left", on="Time")
    country_df = country_df.sort_values(by="Time", ascending=True)

    for currency in SENTIMENT_KEYWORDS:
        # Uses a rolling window to calculate average sentiment
        country_df[currency.upper()] = (
            country_df[currency.upper()].rolling(window, min_periods=1).mean()
        )
    country_df = country_df.fillna(0)
    print(country_df, "COUNTRY DATAFRAME")
    return country_df


def tweet_sentiment(tweets):
    """
    Generates sentiment score from tweets using VADER to get score on -1 to 1 scale.

    Args:
        tweets: List containing dictionaries with tweet data

    Returns:
        Identical list that was passed but with the text property with useful information
        and property for sentiment score for each tweet
    """
    sid = SentimentIntensityAnalyzer()
    for tweet_data in tweets:
        # Removes @s
        tweet_data["text"] = remove_pattern(tweet_data["text"], "RT @[\w]*:")
        tweet_data["text"] = remove_pattern(tweet_data["text"], "@[\w]*")
        # Removes hyperlinks
        tweet_data["text"] = remove_pattern(
            tweet_data["text"], "https?://[A-Za-z0-9./]*")
        # Removes other users
        tweet_data["text"] = tweet_data["text"].replace("[^a-zA-Z]", " ")
        # Removes newline characters
        tweet_data["text"] = tweet_data["text"].replace("\n", " ")
        tweet_data["score"] = sid.polarity_scores(
            tweet_data["text"])["compound"]
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
                    tweets[currency.upper()].iloc[i - length: i].mean()
                )
        elif current == tweets.at[i - length, "Time"]:
            length += 1
        elif length > 1:
            for currency in currencies:
                tweets.at[i - length, currency.upper()] = (
                    tweets[currency.upper()].iloc[i - length: i].mean()
                )
            length = 1
    tweets.drop_duplicates(subset=["Time"], inplace=True)
    return tweets


if __name__ == "__main__":
    print(generate_twitter_sentiment(30, 60))
