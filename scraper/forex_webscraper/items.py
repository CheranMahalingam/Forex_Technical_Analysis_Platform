# Define here the models for your scraped items
#
# See documentation in:
# https://docs.scrapy.org/en/latest/topics/items.html

from scrapy.item import Item, Field
from itemloaders.processors import MapCompose
from w3lib.html import strip_html5_whitespace
from datetime import datetime, timezone
from time import mktime

import scrapy

def unix_timestamp_daily_fx(value):
    try:
        unix_time = datetime.strptime(value, "%Y-%m-%d %H:%M:%S").replace(tzinfo=timezone.utc).timestamp()
    except:
        unix_time = 0
    return unix_time

def unix_timestamp_forex_live(value):
    if len(value) == 24:
        try:
            month_number = datetime.strptime(value[7:10], "%b").month
            unix_time = datetime.strptime(value[4:7] + str(month_number) + value[10:], "%d %m %Y %H:%M:%S").replace(tzinfo=timezone.utc).timestamp()
        except:
            unix_time = 0
    else:
        try:
            month_number = datetime.strptime(value[6:9], "%b").month
            unix_time = datetime.strptime(value[4:6] + str(month_number) + value[9:], "%d %m %Y %H:%M:%S").replace(tzinfo=timezone.utc).timestamp()
        except:
            unix_time = 0
    return unix_time

def unix_timestamp_reuters(value):
    try:
        month_number = datetime.strptime(value[0:3], "%b").month
        unix_time = datetime.strptime(str(month_number) + value[3:], "%m %d %Y").replace(tzinfo=timezone.utc).timestamp()
    except:
        unix_time = 0
    return unix_time


class DailyFXItem(Item):
    
    headline = Field()

    date = Field(
        input_processor=MapCompose(strip_html5_whitespace, unix_timestamp_daily_fx)
    )

class ForexLiveItem(Item):

    headline = Field()

    date = Field(
        input_processor=MapCompose(strip_html5_whitespace, unix_timestamp_forex_live)
    )

class ReutersItem(Item):

    headline = Field(
        input_processor=MapCompose(strip_html5_whitespace)
    )

    summary = Field(
        input_processor=MapCompose(strip_html5_whitespace)
    )

    date = Field(
        input_processor=MapCompose(unix_timestamp_reuters)
    )
