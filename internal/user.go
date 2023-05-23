package internal

import (
	"fmt"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

// get total pages of repositories
func TotalPagesUser(userID string, f *Filter) int {
	counter := 1
	stopper := true
	for stopper {
		// start from root url
		// by increasing counter
		// so we can crawl next pages until none
		url := fmt.Sprintf("https://github.com/%s?tab=repositories&page=%d&q=&type=%s&language=%s&sort=%s", userID, counter, f.Type, f.Lang, f.Sort)
		res, err := http.Get(url)
		if err != nil {
			log.Println("msg=not able to get request", "error=", err)
		}
		defer res.Body.Close()
		doc, err := goquery.NewDocumentFromReader(res.Body)
		if err != nil {
			log.Println("msg=not able to parse body", "error=", err)
		}
		_, ok := doc.Find("#user-repositories-list").Find("div.paginate-container").Find("a.next_page").Attr("href")
		if ok {
			stopper = true
			counter += 1
		} else {
			stopper = false
		}
	}
	log.Printf("INFO user_id=%s msg=page count %d\n", userID, counter)
	return counter
}
