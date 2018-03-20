# -*- coding: utf-8 -*-
import scrapy


class Mm131XingganSpider(scrapy.Spider):
    name = 'mm131_xinggan'
    allowed_domains = ['http://www.mm131.com/xinggan/']
    start_urls = ['http://http://www.mm131.com/xinggan//']

    def parse(self, response):
        pass
