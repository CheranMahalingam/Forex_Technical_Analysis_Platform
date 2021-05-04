from TwitterAPI import TwitterAPI, TwitterPager
import boto3
import re
from vaderSentiment.vaderSentiment import SentimentIntensityAnalyzer
import os
import datetime
from sentiment_keyword_defs import SENTIMENT_KEYWORDS

api = TwitterAPI(
    consumer_key=os.getenv("TWITTER_CONSUMER_KEY"),
    consumer_secret=os.getenv("TWITTER_CONSUMER_SECRET"),
    access_token_key=os.getenv("TWITTER_ACCESS_TOKEN_KEY"),
    access_token_secret=os.getenv("TWITTER_ACCESS_TOKEN_SECRET"),
    api_version='2'
)


def generate_twitter_sentiment(end_time_key):
    testing_interval = datetime.timedelta(minutes=2)

    date = end_time_key['Date']['S'] + " " + end_time_key['Timestamp']['S']
    end_time = datetime.datetime.strptime(date, "%Y-%m-%d %H:%M:%S")
    start_time = end_time - testing_interval
    formatted_start = start_time.isoformat('T') + 'Z'
    formatted_end = end_time.isoformat('T') + 'Z'

    tweets = get_twitter_data(formatted_start, formatted_end)
    sentiment = tweet_sentiment(tweets)
    sentiment = currency_sentiment(sentiment)
    return sentiment


def get_twitter_data(start_time, end_time):
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
            'query': 'from:FXstreetNews OR from:forexcom OR from:markets OR from:ReutersGMF OR from:FTMarkets OR from:WSJmarkets',
            'tweet.fields': 'public_metrics,created_at',
            'start_time': str(start_time),
            'end_time': str(end_time),
            'max_results': 100
        }
    )
    tweet_data = []
    counter = 0
    for item in pager.get_iterator(new_tweets=False):
        tweet_data.append(
            {"text": item['text'], "created_at": item['created_at']})
        counter += 1
    print(counter)
    return tweet_data


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
        input_text: String representing of a twitter post
        pattern: Regex pattern to search for in twitter post

    Returns:
        String with pattern stripped.
    """
    match = re.findall(pattern, input_text)
    for i in match:
        input_text = re.sub(i, "", input_text)
    return input_text


def currency_sentiment(tweets):
    for tweet in tweets:
        tweet['text'] = tweet['text'].lower()

    scores = {"USD": [], "EUR": [], "GBP": [], "CHF": [], "AUD": [], "JPY": [], "CAD": [], "NZD": []}
    for currency in SENTIMENT_KEYWORDS:
        for keyword in SENTIMENT_KEYWORDS[currency]["positive"]:
            currency_score = search_tweets_for_keyword(tweets, keyword, True)
            scores[currency] += currency_score
        for keyword in SENTIMENT_KEYWORDS[currency]["negative"]:
            currency_score = search_tweets_for_keyword(tweets, keyword, False)
            scores[currency] += currency_score

    for currency in scores:
        scores[currency] = list_average(scores[currency])
    return scores


def search_tweets_for_keyword(tweets, keyword, positive):
    score = []
    for tweet in tweets:
        if keyword in tweet['text']:
            score.append(tweet['score']) if positive else score.append(-tweet['score'])
    return score


def list_average(score_list):
    if len(score_list) == 0:
        return -100
    return sum(score_list) / len(score_list)
