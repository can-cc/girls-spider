package main

import (
	"github.com/gocolly/colly"
	"log"
)

type MMPage struct {
	Url string
}

type MMItem struct {
	Url   string
	Title string
}

type MMDetailImage struct {
	MMItemTitle string
	ImageUrl    string
}

func main() {

	mmPageCollector := colly.NewCollector(
		colly.AllowedDomains("m.mm131.net"),
		colly.CacheDir("./colly_cache"),
		colly.UserAgent("Mozilla/5.0 (iPhone; CPU iPhone OS 13_2_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/13.0.3 Mobile/15E148 Safari/604.1"),

	)

	//mmPages := make([]MMPage, 0, 200)

	mmPageCollector.OnHTML("#xbtn", func(e *colly.HTMLElement) {
		// If attribute class is this long string return from callback
		// As this a is irrelevant

		link := "https://m.mm131.net/xinggan/" + e.Attr("href")
		// If link start with browse or includes either signup or login return from callback
		log.Println(link)
		e.Request.Visit(link)
	})

	mmPageCollector.OnResponse(func(response *colly.Response) {
		log.Println(response.Request.URL)
	})

	mmPageCollector.Visit("https://m.mm131.net/xinggan/")
}
