package githubscrape

import (
	"fmt"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

// get total pages of repositories
func TotalPagesUser(userID string) int {
	counter := 1
	stopper := true
	for stopper {
		// start from root url
		// by increasing counter
		// so we can crawl next pages until none
		url := fmt.Sprintf("https://github.com/%s?tab=repositories&page=%d", userID, counter)
		res, err := http.Get(url)
		if err != nil {
			log.Println(err)
		}
		defer res.Body.Close()
		doc, err := goquery.NewDocumentFromReader(res.Body)
		if err != nil {
			log.Println(err)
		}
		_, ok := doc.Find("#user-repositories-list").Find("div.paginate-container").Find("a.next_page").Attr("href")
		if ok {
			stopper = true
			counter += 1
		} else {
			stopper = false
		}
	}
	return counter
}
