package githubscrape

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

// scrapes data of repos based on orgnizations name
func Scrape(orgName string) {
	total := TotalPages(orgName)
	// make chan of total available repos
	// 1 page has 30 repos
	reposCh := make(chan *Repo, total*30)
	var wg sync.WaitGroup

	wg.Add(total)
	for p := 1; p <= total; p++ {
		// generate url per page
		url := fmt.Sprintf("https://github.com/orgs/%s/repositories?page=%d", orgName, p)
		ProcessPage(true, url, orgName, reposCh, &wg)
	}

	wg.Wait()
	close(reposCh)

	// create file after closing channel
	CreateFile(fmt.Sprintf("orgs/%s.json", orgName), reposCh)
}

// get total pages of repositories
func TotalPages(orgName string) int {
	rootURL := fmt.Sprintf("https://github.com/orgs/%s/repositories", orgName)
	res, err := http.Get(rootURL)
	if err != nil {
		log.Println(err)
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Println(err)
	}

	pages, ok := doc.Find("#org-repositories").Find("div.pagination").Find("em.current").Attr("data-total-pages")
	if ok {
		pagesInt, _ := strconv.Atoi(pages)
		return pagesInt
	}
	return 1
}
