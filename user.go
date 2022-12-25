package githubscrape

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

// scrapes data of repos based on orgnizations name
func ScrapeUser(userID string) {
	total := TotalPagesUser(userID)

	// make chan of total available repos
	// 1 page has 30 repos
	reposUserCh := make(chan *Repo, total*30)
	var wg sync.WaitGroup

	wg.Add(total)
	for p := 1; p <= total; p++ {
		// generate url per page
		url := fmt.Sprintf("https://github.com/%s?tab=repositories&page=%d", userID, p)
		ProcessPage(false, url, userID, reposUserCh, &wg)
	}

	wg.Wait()
	close(reposUserCh)

	// create file after closing channel
	CreateFile(fmt.Sprintf("users/%s.json", userID), reposUserCh)
}

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
