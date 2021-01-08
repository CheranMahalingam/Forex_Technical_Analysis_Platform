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

def search(users):
    for user in users:
        data = {}
        data['headline'] = []
        data['date'] = []

        for year in range(2018, 2021):
            for month in range(1, 13):
                feb_days = "29" if year%4 == 0 else "28"
                days_per_month = {1: "31", 2: feb_days, 3: "31", 4: "30", 5: "31", 6: "30", 7: "31", 8: "31", 9: "30", 10: "31", 11: "30", 12: "31"}
                str_month = str(month) if month > 9 else "0" + str(month)
                url = "https://www.twitter.com/search?lang=en&q=(from%3A" + user + ")%20until%3A" + str(year) + "-" + str_month \
                        + "-" + days_per_month[month] + "%20since%3A" + str(year) + "-" + str_month + "-01%20&src=typed_query&f=live"
                scrape_data(url, data)
        
        store_data(user, data)

def scrape_data(url, data):
    tweet_count = 0
    TWEET_MAX = 3000

    driver.get(url)
    sleep_for(5, 9)

    tweet_store = set()
    post_element_xpath = "//div[@data-testid='tweet']"
    current_time = time()

    while tweet_count < TWEET_MAX:
        if time() - current_time > 20:
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
                data['headline'].append(post_text)
                data['date'].append(post_date)
                tweet_store.add(post_text)
                tweet_count += 1
                print(tweet_count)
                if tweet_count == TWEET_MAX:
                    break
        
        new_post_list = driver.find_elements_by_xpath(post_element_xpath)

        while new_post_list == post_list:
            if(time() - current_time) > 20:
                tweet_count == TWEET_MAX
                break
            
            driver.execute_script("window.scrollBy(0, 200);")
            new_post_list = driver.find_elements_by_xpath(post_element_xpath)

        sleep_for(1, 2)

def sleep_for(opt1, opt2):
    time_for = random.uniform(opt1, opt2)
    time_for_int = int(round(time_for))
    sleep(abs(time_for_int - time_for))
    for i in range(time_for_int, 0, -1):
        sleep(1)

def store_data(user, data):
    with open(str(user) + "_historical.json", 'w') as outfile:
        json.dump(data, outfile)

def close_driver():
    driver.close()


if __name__ == "__main__":
    users = ["ReutersGMF", "forexcom", "FXstreetNews", "FTMarkets", "markets", "WSJmarkets"]
    search(users)
    close_driver()
