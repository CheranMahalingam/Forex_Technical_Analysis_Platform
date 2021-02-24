import scrapy
from scrapy.loader import ItemLoader
from forex_data.items import ReutersItem

class ReutersHistoricalSpider(scrapy.Spider):

    name = "reuters_data"
    allowed_domains = ["www.reuters.com"]
    start_urls = ["https://www.reuters.com/news/archive/GCA-ForeignExchange"]
    COUNT_MAX = 231
    count = 0

    def parse(self, response):
        self.count += 1
        if self.count < self.COUNT_MAX:
            loader = ItemLoader(ReutersItem(), response)
            loader.add_xpath('headline', "//article[@class='story ']//h3[@class='story-title']/text()")
            loader.add_xpath('summary', "//article[@class='story ']//p/text()")
            loader.add_xpath('date', "//article[@class='story ']//span[@class='timestamp']/text()")
            yield loader.load_item()

            next_page = response.xpath("//a[@class='control-nav-next']/@href").get()
            yield scrapy.Request(response.urljoin(next_page))
