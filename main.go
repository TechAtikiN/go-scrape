package main

// import packages
import (
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// type for SearchResult
type SearchResult struct {
	ResultRank int
	ResultURL string
	ResultTitle string
	ResultDesc string
}

// map for bing domains
var bingDomains = map[string]string {
	"com": "",
	"ar": "&cc=ar",
	"at": "&cc=at",
	"au": "&cc=au",
	"be": "&cc=be",
	"br": "&cc=br",
	"ca": "&cc=ca",
	"cl": "&cc=cl",
	"dk": "&cc=dk",
	"fi": "&cc=fi",
	"fr": "&cc=fr",
	"de": "&cc=de",
	"hk": "&cc=hk",
	"in": "&cc=in",
	"ie": "&cc=ie",
	"it": "&cc=it",
	"jp": "&cc=jp",
	"mx": "&cc=mx",
	"nl": "&cc=nl",
	"nz": "&cc=nz",
	"no": "&cc=no",
	"cn": "&cc=cn",
	"pl": "&cc=pl",
	"pt": "&cc=pt",
	"za": "&cc=za",
	"es": "&cc=es",
	"se": "&cc=se",
	"ch": "&cc=ch",
	"tw": "&cc=tw",
	"tr": "&cc=tr",
	"gb": "&cc=gb",
	"us": "&cc=us",
	"ae": "&cc=ae",
	"ve": "&cc=ve",
	"vn": "&cc=vn",
	"bg": "&cc=bg",
	"hr": "&cc=hr",
	"cz": "&cc=cz",
	"ee": "&cc=ee",
	"gr": "&cc=gr",
	"hu": "&cc=hu",
	"id": "&cc=id",
	"il": "&cc=il",
	"lv": "&cc=lv",
	"lt": "&cc=lt",
	"my": "&cc=my",
	"ph": "&cc=ph",
	"ro": "&cc=ro",
	"ru": "&cc=ru",
	"sa": "&cc=sa",
	"sg": "&cc=sg",
	"sk": "&cc=sk",
	"si": "&cc=si",
	"th": "&cc=th",
	"ua": "&cc=ua",
	"uy": "&cc=uy",
}

// slice for userAgents
var userAgents = []string {
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_6)" +
	"AppleWebKit/537.36 (KHTML, like Gecko) Chrome/53.0.2785.143 Safari/537.36",
	"Mozilla/5.0 (Windows NT 6.1; WOW64)" +
	"AppleWebKit/537.36 (KHTML, like Gecko) Chrome/53.0.2785.143 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; WOW64)" +
	"AppleWebKit/537.36 (KHTML, like Gecko) Chrome/53.0.2785.143 Safari/537.36",
	"Mozilla/5.0 (Windows NT 6.1; WOW64; rv:47.0)" +
	"Gecko/20100101 Firefox/47.0",
	"Mozilla/5.0 (Windows NT 10.0; WOW64; rv:47.0)" +
	"Gecko/20100101 Firefox/47.0",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_6)" +
	"AppleWebKit/601.7.7 (KHTML, like Gecko) Version/9.1.2 Safari/601.7.7",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_6)" +
	"AppleWebKit/537.36 (KHTML, like Gecko) Chrome/52.0.2743.116 Safari/537.36",
	"Mozilla/5.0 (Windows NT 6.1; WOW64)" +
	"AppleWebKit/537.36 (KHTML, like Gecko) Chrome/52.0.2743.116 Safari/537.36",
}

// main function
func main() {
	res, err := BingScrape("Nelson Mandela", "com", nil, 1, 15, 30)

	if err == nil {
		for _, r := range res {
			fmt.Println(r)
		}
	} else {
		fmt.Println(err)
	}
}

// scraping bing search results
func BingScrape(searchTerm, country string, proxyString interface{}, pages, count, backoff int) ([]SearchResult, error) {
	results := []SearchResult{}

	bindPages, err := buildBingUrls(searchTerm, country, pages, count)

	if err != nil {
		return nil, err
	}

	for _, page := range bindPages {

		rank := len(results)
		res, err := scrapeClientRequest(page, proxyString)

		if err != nil {
			return nil, err
		}

		data, err := bingResultParser(res, rank)
		if err != nil {
			return nil, err
		}

		for _, result := range data {
			results = append(results, result)
		}
	time.Sleep(time.Duration(backoff)*time.Second)
	}
	return results, nil
}

// building bing urls
func buildBingUrls(searchTerm, country string, pages, count int) ([]string, error) {
  toScrape := []string{}
	searchTerm = strings.Trim(searchTerm, " ")
	searchTerm = strings.Replace(searchTerm, " ", "+", -1)

	if countryCode, found := bingDomains[country]; found {
		for i := 0; i < pages; i++ {
			first := firstParameter(i, count);
			scrapeUrl := fmt.Sprintf("https://bing.com/search?q=%s&first=%d&count=%d%s",
			searchTerm, first, count, countryCode)
			toScrape = append(toScrape, scrapeUrl)
		}
	} else {
		err := fmt.Errorf("Country code not found: %s", country)
		return nil, err
	}

	return toScrape, nil
}

// getting first parameter
func firstParameter(number, count int) int {
	if number == 0 {
		return number + 1
	}
	return number * count + 1
}

// scraping client request
func scrapeClientRequest(searchURL string, proxyString interface{}) (*http.Response, error) {
	baseClient := getScrapeClient(proxyString)
	req , _ := http.NewRequest("GET", searchURL, nil)
	req.Header.Set("User-Agent", randomUserAgent())

	res, err := baseClient.Do(req)
	if res.StatusCode != 200 {
		err := fmt.Errorf("Status code error: %d %s", res.StatusCode, res.Status)
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	return res, nil
}

// getting scrape client
func getScrapeClient(proxyString interface{}) *http.Client {
	switch V := proxyString.(type) {
	case string:
		proxyURL, _ := url.Parse(V)
		return &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyURL)}}

	default:
		return &http.Client{}
	}
}

// getting random user agent
func randomUserAgent() string {
	rand.Seed(time.Now().Unix())
	randNum := rand.Intn(len(userAgents))
	return userAgents[randNum]
}

// bing result parser
func bingResultParser(response *http.Response, rank int)([]SearchResult, error) {

	doc, err := goquery.NewDocumentFromResponse(response)
	if err != nil {
		return nil, err
	}

	results := []SearchResult{}
	sel := doc.Find("li.b_algo")
	rank++

	for i := range sel.Nodes {
		item := sel.Eq(i)
		linkTag := item.Find("a")
		link, _ := linkTag.Attr("href")
		titleTag := item.Find("h2")
		descTag := item.Find("div.b_caption p")
		desc := descTag.Text()
		title := titleTag.Text()
		link = strings.Trim(link, " ")

		if link != "" && title != "#" && !strings.HasPrefix(link, "/") {
			result := SearchResult{rank, link, title, desc}
			results = append(results, result)
			rank++
		}
	}
	return results, err
}