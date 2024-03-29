package main

import (
	"github.com/gocolly/colly"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"fmt"
)

type MMItem struct {
	Url         string
	Title       string
	ItemID      string
	ImageNumber int
	Tags        []string
}

func downloadImage(url string, filePath string, referer string) (bool, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("User-Agent", `"Mozilla/5.0 (iPhone; CPU iPhone OS 13_2_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/13.0.3 Mobile/15E148 Safari/604.1"`)
	req.Header.Add("Accept", "image/webp,*/*")
	req.Header.Add("Host", "img1.mmmw.net")
	req.Header.Add("Referer", "https://m.mm131.net/")
	response, err := client.Do(req)
	if err != nil {
		return false, err
	}

	defer response.Body.Close()
	if response.StatusCode == 404 || response.StatusCode == 599 {
		return true, nil
	}

	//open a file for writing
	file, err := os.Create(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Use io.Copy to just dump the response body to the file. This supports huge files
	_, err = io.Copy(file, response.Body)
	if err != nil {
		panic(err)
	}
	return false, nil
}

func handleItemDetect(item *MMItem, resultDir string) error {
	dirPath := resultDir + "/" + item.Title
	if _, err := os.Stat(dirPath); !os.IsNotExist(err) {
		return nil
	}
	if err := os.Mkdir(dirPath, 0755); err != nil {
		return err
	}
	is3 := false
	successNum := 0
	for i := 1; ; i++ {
		filePath := dirPath + "/" + strconv.Itoa(i) + ".jpg"
		mmIndexNumber := strconv.Itoa(i)
		if is3 {
			mmIndexNumber = fmt.Sprintf("/%03d", i)
		}
		url := "https://img1.hnllsy.com/pic/" + item.ItemID + "/" + mmIndexNumber + ".jpg"
		done, err := downloadImage(url, filePath, item.Url)
		if done && i == 1 {
			is3 = true
			i = 0
			continue
		}
		successNum = successNum + 1
		if done {
			break
		}
		if err != nil {
			panic(err)
		}
	}
	fmt.Println("success number = ", + successNum)
	return nil
}

func main() {
	var resultDir string
	if len(os.Args) > 1 {
		resultDir = os.Args[1]
	} else {
		resultDir = "result"
	}

	if _, err := os.Stat(resultDir); os.IsNotExist(err) {
		if err := os.Mkdir(resultDir, 0755); err != nil {
			panic(err)
		}
	}

	mmPageCollector := colly.NewCollector(
		colly.AllowedDomains("m.mm131.net"),
		colly.CacheDir("./colly_cache"),
		colly.DetectCharset(),
		colly.Async(true),
		colly.UserAgent("Mozilla/5.0 (iPhone; CPU iPhone OS 13_2_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/13.0.3 Mobile/15E148 Safari/604.1"),
	)
	mmPageCollector.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: 100})

	mmPageCollector.OnHTML(".post-content", func(e *colly.HTMLElement) {
		url := e.ChildAttr("a", "href")
		title := e.ChildAttr("img", "alt")
		re := regexp.MustCompile(`([0-9]+)\.html$`)
		id := re.FindStringSubmatch(url)[1]

		err := handleItemDetect(&MMItem{
			Url:    url,
			Title:  title,
			ItemID: id,
		}, resultDir)
		if err != nil {
			panic(err)
		}
		log.Println(title, url)
	})

	mmPageCollector.OnHTML("#xbtn", func(e *colly.HTMLElement) {
		link := "https://m.mm131.net/xinggan/" + e.Attr("href")
		log.Println("link", link)
		e.Request.Visit(link)
	})

	mmPageCollector.Visit("https://m.mm131.net/xinggan/")
	mmPageCollector.Wait()
}
