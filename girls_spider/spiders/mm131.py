# -*- coding: utf-8 -*-
import scrapy
from girls_spider.items import GirlsSpiderItem


class Mm131Spider(scrapy.Spider):
    name = "girls_spider"
    allowed_domains = ["www.mm131.com"]
    start_urls = ['http://www.mm131.com/xinggan/',
                  'http://www.mm131.com/qingchun/',
                  'http://www.mm131.com/xiaohua/',
                  'http://www.mm131.com/chemo/',
                  'http://www.mm131.com/qipao/',
                  'http://www.mm131.com/mingxing/'
                  ]

    def parse(self, response):
        list = response.css(".list-left dd:not(.page)")
        for img in list:
            imgname = img.css("a::text").extract_first()
            imgurl = img.css("a::attr(href)").extract_first()
            imgurl2 = str(imgurl)
            print(imgurl2)
            next_url = response.css(".page-en:nth-last-child(2)::attr(href)").extract_first()
            if next_url is not None:
                # next page
                yield response.follow(next_url, callback=self.parse)

            yield scrapy.Request(imgurl2, callback=self.content)

    def content(self, response):
        item = GirlsSpiderItem()
        item['name'] = response.css(".content h5::text").extract_first()
        item['ImgUrl'] = response.css(".content-pic img::attr(src)").extract()
        yield item
        next_url = response.css(".page-ch:last-child::attr(href)").extract_first()

        if next_url is not None:
            print("下一页", next_url)
            print("下一页", response.url)
            yield response.follow(next_url, callback=self.content)
