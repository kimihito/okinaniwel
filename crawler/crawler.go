package crawler

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/gocolly/colly/v2"
)

var dateRe = regexp.MustCompile(`[0-9]{1}\/[0-9]{1,2}\s*(\W*)（\W*）`)

func Run() {
	// Instantiate default collector
	c := colly.NewCollector(
		// Visit only domains: aniwel-pref.okinawa
		colly.AllowedDomains("www.aniwel-pref.okinawa"),
	)

	detailCollector := c.Clone()

	// On every a element which has href attribute call callback
	c.OnHTML("#main .lists > a", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		// Print link
		fmt.Printf("Link found: %q -> %s\n", strings.TrimSpace(e.Text), link)
		detailCollector.Visit(e.Request.AbsoluteURL(link))
	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	detailCollector.OnHTML("body", func(e *colly.HTMLElement) {
		title := dateRe.FindStringSubmatch(e.DOM.Find(".single_title").Text())
		fmt.Println(title)
		// e.DOM.Find(".box").First().Find("tr>td").Each(func(i int, s *goquery.Selection) {
		// 	fmt.Println(i, s.Text())
		// })
	})

	// Start scraping on https://hackerspaces.org
	c.Visit("https://www.aniwel-pref.okinawa/animals/accommodate/dogs")
	fmt.Println("Running crawl")
}
