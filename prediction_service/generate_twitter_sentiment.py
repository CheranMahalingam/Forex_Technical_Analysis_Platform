from TwitterAPI import TwitterAPI, TwitterPager
from dotenv import load_dotenv
import os
load_dotenv()

api = TwitterAPI(consumer_key=os.getenv("TWITTER_CONSUMER_KEY"), consumer_secret=os.getenv("TWITTER_CONSUMER_SECRET"),
                 access_token_key=os.getenv("TWITTER_ACCESS_TOKEN_KEY"), access_token_secret=os.getenv("TWITTER_ACCESS_TOKEN_SECRET"), api_version='2')

def get_twitter_data(start_time):
    pager = TwitterPager(api, 'tweets/search/recent', {
        'query': 'from:FXstreetNews OR from:forexcom',
        'tweet.fields': 'public_metrics,created_at',
        'start_time': str(start_time),
        'max_results': 100
        }
    )
    tweet_data = []
    for item in pager.get_iterator(new_tweets=False):
        tweet_data.append({"text": item['text'], "created_at": item['created_at']})
        print(item)
    return tweet_data
