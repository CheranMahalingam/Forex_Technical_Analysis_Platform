import scrapy
from scrapy.loader import ItemLoader
from forex_data.items import ForexLiveItem

class ForexLiveHistoricalWebscraper(scrapy.Spider):

    name = "forexlive_data"
    allowed_domains = ["www.forexlive.com"]
    start_urls = ["https://www.forexlive.com/technical-analysis/Headlines/1"]

    def parse(self, response):
        if response.xpath("//time/text()[2]"):
            loader = ItemLoader(ForexLiveItem(), response)
            loader.add_xpath('headline', "//article[@class='row axa-aprev']//h3/a/text()")
            loader.add_xpath('date', "//time/text()[2]")
            yield loader.load_item()

            next_page = response.xpath("//h3/a[2]/@href").get()
            yield scrapy.Request(response.urljoin(next_page))
