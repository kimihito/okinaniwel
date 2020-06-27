package crawler

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
	"github.com/kimihito/okinaniwel/model"
	"github.com/mitchellh/mapstructure"
)

var dateRe = regexp.MustCompile(`[0-9]{1}\/[0-9]{1,2}\s*(\W*)（\W*）`)

var dogs = []model.Dog{}

var indexToMapKey = []string{"id", "addedOn", "deadlineOn", "place", "color", "sex", "size", "age", "collar", "note", "image"}

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
		titleText := e.DOM.Find(".single_title").Text()
		// 返還した犬に関する情報はスキップする
		if strings.Contains(titleText, "返還しました") {
			return
		}

		title := dateRe.FindStringSubmatch(titleText)[1]

		items := map[string]string{"breed": strings.TrimSpace(title)}
		e.DOM.Find(".box").First().Find("tr>td").Each(func(i int, s *goquery.Selection) {
			if s.HasClass("photo") {
				if src, ok := s.Find("img").First().Attr("src"); ok {
					items["image"] = src
				}
			} else {
				text := s.Text()
				// PDFのリンクは除く
				if !strings.Contains(text, "PDF") {
					items[indexToMapKey[i]] = strings.Join(strings.Fields(s.Text()), " ")
				}
			}
		})

		dog := &model.Dog{}

		mapstructure.Decode(items, dog)
		dogs = append(dogs, *dog)

	})

	// Start scraping on https://www.aniwel-pref.okinawa
	c.Visit("https://www.aniwel-pref.okinawa/animals/accommodate/dogs")
	fmt.Println("Finish")
	fmt.Println(dogs)

}
