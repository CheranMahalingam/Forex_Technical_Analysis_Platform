from selenium import webdriver
from selenium.webdriver.common.keys import Keys
from selenium.webdriver.chrome.options import Options
from time import sleep, time
from dotenv import load_dotenv
import os
import json
import random

load_dotenv()

browser_path = os.getenv("SELENIUM_PATH")
chrome_options = Options()
chrome_options.add_argument("--headless")

driver = webdriver.Chrome(executable_path=browser_path, options=chrome_options)

urls = ["https://www.twitter.com/ReutersGMF?lang=en",
        "https://twitter.com/forexcom?lang=en",
        "https://twitter.com/FXstreetNews?lang=en"]

data = {}
data['headline'] = []
data['date'] = []

def sleep_for(opt1, opt2):
    time_for = random.uniform(opt1, opt2)
    time_for_int = int(round(time_for))
    sleep(abs(time_for_int - time_for))
    for i in range(time_for_int, 0, -1):
        sleep(1)

def scrape_data(url):
    tweet_count = 0
    TWEET_MAX = 1

    driver.get(url)
    sleep_for(10, 15)

    tweet_store = set()
    post_element_xpath = "//div[@data-testid='tweet']"
    current_time = time()

    while tweet_count < TWEET_MAX:
        if time() - current_time > 120:
            break
        post_list = driver.find_elements_by_xpath(post_element_xpath)
        for post in post_list:
            try:
                post_date = post.find_element_by_xpath(".//time").get_attribute("datetime")
                post_text = post.find_element_by_xpath(".//div[2]/div[2]/div[1]").text
            except:
                post_date = 0
                post_text = 0

            if post_date != 0 and post_text not in tweet_store:
                current_time = time()
                print(tweet_count)
                data['headline'].append(post_text)
                data['date'].append(post_date)
                tweet_store.add(post_text)
                tweet_count += 1
                if tweet_count == TWEET_MAX:
                    break
        
        new_post_list = driver.find_elements_by_xpath(post_element_xpath)
        scroll_count = 0

        while new_post_list == post_list:
            if(time() - current_time) > 120:
                tweet_count == TWEET_MAX
                break
            driver.execute_script("window.scrollBy(0, 200);")
            new_post_list = driver.find_elements_by_xpath(post_element_xpath)
            scroll_count += 1

        sleep_for(1, 3)
        store_data()

def store_data():
    with open('twitter.json', 'w') as outfile:
        json.dump(data, outfile)

def close_driver():
    driver.close()


if __name__ == "__main__":
    for i in urls:
        scrape_data(i)
    close_driver()