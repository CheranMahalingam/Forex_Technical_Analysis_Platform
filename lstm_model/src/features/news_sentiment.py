"""Feature engineering for sentiment around a currency using news headlines."""

from datetime import datetime
from functools import reduce
from vaderSentiment.vaderSentiment import SentimentIntensityAnalyzer
import pandas as pd


def news_unpacking(org):
    """
    Reads historical news headlines from a json file and writes them into a dataframe.

    Args:
        org: String representing the news organization from which news is read

    Returns:
        Dataframe containing news headlines, time, and summaries.
    """
    org_news = pd.read_json("lstm_model/data/raw/news/{}_historical.json".format(org))
    headline_arr = []
    date_arr = []
    summary_arr = []

    for month in range(org_news["headline"].count()):
        for headline in org_news["headline"][month]:
            headline_arr.append(headline)
        for date in org_news["date"][month]:
            date_arr.append(date)
        if "summary" in org_news:
            for summary in org_news["summary"][month]:
                summary_arr.append(summary)
    if summary_arr:
        org_df = pd.DataFrame({"Time": date_arr, "Headline": headline_arr, "Summary": summary_arr})
    else:
        org_df = pd.DataFrame({"Time": date_arr, "Headline": headline_arr,})
    org_df["Time"] = org_df["Time"].transform(
        lambda time: datetime.utcfromtimestamp(time).strftime("%Y-%m-%d %H:%M:00"))
    org_df.sort_values(by=["Time"], inplace=True)
    org_df = org_df.reset_index(drop=True)
    return org_df

def news_merge(news):
    """
    Merges array of news dataframes into one.

    Args:
        news: Array of news dataframes

    Returns:
        Single dataframe containing all headlines from different news sources.
    """
    if len(news) == 0:
        return None
    if len(news) == 1:
        return news[0]
    merged_news = reduce(
        lambda left, right: pd.merge(left, right, how="outer", on=["Time", "Headline"]), news)
    merged_news.sort_values(by=["Time"], inplace=True)
    merged_news = merged_news.reset_index(drop=True)
    return merged_news

def generate_sentiment_score(news):
    """
    Generates a sentiment score for news headlines and summaries.

    Args:
        news: Dataframe containing time, headline, and summary of news

    Returns:
        Dataframe with additional column containing sentiment score.
    """
    sid = SentimentIntensityAnalyzer()
    score = []
    for headline, summary in zip(news.Headline, news.Summary):
        if pd.isna(summary):
            text = headline
            score.append(sid.polarity_scores(text)["compound"])
        else:
            text = headline + summary
            score.append(sid.polarity_scores(text)["compound"])
    news_score = pd.DataFrame({"News_Sentiment": score})
    news["News_Sentiment"] = news_score
    return news

def currency_sentiment(currencies_dict, news):
    """
    Generates csv file containing news sentiment scores for each country according to date.

    Args:
        currencies_dict: Dictionary holding keywords to indicate that a headline affects the
                         strength of a currency and whether is positive or negative
        news: Dataframe containing news headlines with their sentiment scores
    """
    country_df = pd.DataFrame()
    for currency in currencies_dict:
        for entity in currencies_dict[currency]["positive"]:
            headline_lower = news["Headline"].transform(lambda headline: headline.lower())
            summary_lower = news["Summary"].transform(
                lambda summary: summary.lower() if type(summary) == "object" else summary)
            currency_df = news[(headline_lower.str.contains(entity)) | \
                (summary_lower.str.contains(entity))]
            currency_df = currency_df[{"Time", "News_Sentiment"}]
            currency_df = currency_df.rename(columns={"News_Sentiment": currency.upper()})
            if country_df.empty:
                country_df = currency_df
            elif not currency.upper() in country_df.columns:
                country_df = country_df.merge(currency_df, how="outer", on="Time")
            else:
                country_df = country_df.merge(
                    currency_df, how="outer", on=["Time", currency.upper()])
        for entity in currencies_dict[currency]["negative"]:
            headline_lower = news["Headline"].transform(lambda headline: headline.lower())
            summary_lower = news["Summary"].transform(
                lambda summary: summary.lower() if type(summary) == "object" else summary)
            currency_df = news[(headline_lower.str.contains(entity)) | \
                (summary_lower.str.contains(entity))]
            currency_df = currency_df[{"Time", "News_Sentiment"}]
            if not currency_df["News_Sentiment"].empty:
                currency_df["News_Sentiment"] = currency_df["News_Sentiment"].transform(
                    lambda score: -score)
            currency_df = currency_df.rename(columns={"News_Sentiment": currency.upper()})
            if country_df.empty:
                country_df = currency_df
            elif not currency.upper() in country_df.columns:
                country_df = country_df.merge(currency_df, how="outer", on="Time")
            else:
                country_df = country_df.merge(
                    currency_df, how="outer", on=["Time", currency.upper()])

    country_df = combine_dates(country_df)

    time_frame = pd.date_range(start="2018-01-01 22:00:00", freq="1T", end="2020-12-31 21:59:00")
    time_frame = pd.DataFrame(time_frame, columns=["Time"])
    time_frame["Time"] = time_frame["Time"].dt.strftime("%Y-%m-%d %H:%M:%S")

    country_df = time_frame.merge(country_df, how="outer", on="Time")
    country_df = country_df.sort_values(by='Time', ascending=True)

    for currency in currencies_dict:
        country_df[currency.upper()] = country_df[currency.upper()].rolling(
            1440, min_periods=1).mean()
    country_df = country_df.fillna(0)

    country_df.to_csv("lstm_model/data/interim/news/news_sentiment.csv", index=False)

def combine_dates(news):
    """
    Merge sentiment scores according to date.

    Args:
        news: Dataframe containing countries and their sentiment scores at a certain time

    Returns:
        Dataframe with a country's sentiment score with sequential time.
    """
    currencies = ["eur", "usd", "jpy", "cad", "gbp", "aud", "nzd", "chf"]
    length = 1
    for i in range(1, len(news.index)):
        current = news.at[i, "Time"]
        if current == news.at[i - length, "Time"] and i == len(news.index) - 1:
            for currency in currencies:
                news.at[i - length, currency.upper()] = news[
                    currency.upper()].iloc[i - length: i].mean()
        elif current == news.at[i - length, "Time"]:
            length += 1
        elif length > 1:
            for currency in currencies:
                news.at[i - length, currency.upper()] = news[
                    currency.upper()].iloc[i - length: i].mean()
            length = 1
    news.drop_duplicates(subset=["Time"], inplace=True)
    return news
