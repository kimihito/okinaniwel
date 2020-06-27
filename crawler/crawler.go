package crawler

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/queue"
	"github.com/kimihito/okinaniwel/model"
	"github.com/mitchellh/mapstructure"
)

var url = "https://www.aniwel-pref.okinawa/animals/accommodate/dogs"
var dateRe = regexp.MustCompile(`[0-9]{1}\/[0-9]{1,2}\s*(\W*)（\W*）`)

var dogs = []model.Dog{}

var indexToMapKey = []string{"id", "addedOn", "deadlineOn", "place", "color", "sex", "size", "age", "collar", "note", "image"}

func Run() []model.Dog {
	// Instantiate default collector
	c := colly.NewCollector(
		// Visit only domains: aniwel-pref.okinawa
		colly.AllowedDomains("www.aniwel-pref.okinawa"),
	)

	detailCollector := c.Clone()

	q, _ := queue.New(2, &queue.InMemoryQueueStorage{MaxSize: 10000})

	// On every a element which has href attribute call callback
	c.OnHTML("#main .lists > a", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		// Print link
		fmt.Printf("Link found: %q -> %s\n", strings.Join(strings.Fields(e.Text), " "), link)
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

		items := map[string]string{"breed": strings.TrimSpace(title), "url": e.Request.URL.String()}
		e.DOM.Find(".box").First().Find("tr>td").Each(func(i int, s *goquery.Selection) {
			if s.HasClass("photo") {
				if src, ok := s.Find("img").First().Attr("src"); ok {
					items["image"] = e.Request.URL.Scheme + "://" + e.Request.URL.Host + src
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
	// c.Visit(url)
	for i := 1; i < 5; i++ {
		q.AddURL(fmt.Sprintf("%s/page:%d", url, i))
	}
	q.Run(c)
	fmt.Println("Crwaling is complete")
	return dogs
}
