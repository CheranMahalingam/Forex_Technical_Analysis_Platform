import scrapy
from scrapy.loader import ItemLoader
from forex_data.items import DailyFXItem

class DailyFXHistoricalSpider(scrapy.Spider):

    name = "dailyfx_data"
    allowed_domains = ["www.dailyfx.com"]
    start_urls = ["https://www.dailyfx.com/archive/2018/01",
                "https://www.dailyfx.com/archive/2018/02",
                "https://www.dailyfx.com/archive/2018/03",
                "https://www.dailyfx.com/archive/2018/04",
                "https://www.dailyfx.com/archive/2018/05",
                "https://www.dailyfx.com/archive/2018/06",
                "https://www.dailyfx.com/archive/2018/07",
                "https://www.dailyfx.com/archive/2018/08",
                "https://www.dailyfx.com/archive/2018/09",
                "https://www.dailyfx.com/archive/2018/10",
                "https://www.dailyfx.com/archive/2018/11",
                "https://www.dailyfx.com/archive/2018/12",
                "https://www.dailyfx.com/archive/2019/01",
                "https://www.dailyfx.com/archive/2019/02",
                "https://www.dailyfx.com/archive/2019/03",
                "https://www.dailyfx.com/archive/2019/04",
                "https://www.dailyfx.com/archive/2019/05",
                "https://www.dailyfx.com/archive/2019/06",
                "https://www.dailyfx.com/archive/2019/07",
                "https://www.dailyfx.com/archive/2019/08",
                "https://www.dailyfx.com/archive/2019/09",
                "https://www.dailyfx.com/archive/2019/10",
                "https://www.dailyfx.com/archive/2019/11",
                "https://www.dailyfx.com/archive/2019/12",
                "https://www.dailyfx.com/archive/2020/01",
                "https://www.dailyfx.com/archive/2020/02",
                "https://www.dailyfx.com/archive/2020/03",
                "https://www.dailyfx.com/archive/2020/04",
                "https://www.dailyfx.com/archive/2020/05",
                "https://www.dailyfx.com/archive/2020/06",
                "https://www.dailyfx.com/archive/2020/07",
                "https://www.dailyfx.com/archive/2020/08",
                "https://www.dailyfx.com/archive/2020/09",
                "https://www.dailyfx.com/archive/2020/10",
                "https://www.dailyfx.com/archive/2020/11",
                "https://www.dailyfx.com/archive/2020/12"]

    def parse(self, response):
        loader = ItemLoader(DailyFXItem(), response)
        loader.add_xpath('headline', "//section[@class='my-6']//span[@class='dfx-articleListItem__title jsdfx-articleListItem__title font-weight-bold align-middle']/text()")
        loader.add_xpath('date', "//section[@class='my-6']//span[@class='jsdfx-articleListItem__date text-nowrap']/text()")

        yield loader.load_item()
