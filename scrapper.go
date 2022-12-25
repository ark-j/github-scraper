package githubscrape

import (
	"fmt"
	"sync"
)

// scrapes data of repos based on orgnizations name
func Scrape(isOrg bool, entityID string) {
	var urlCreate func(s string, p int) string
	var total int
	var fpath string
	switch {
	case isOrg:
		total = TotalPages(entityID)
		fpath = fmt.Sprintf("orgs/%s.json", entityID)
		urlCreate = func(s string, p int) string {
			return fmt.Sprintf("https://github.com/orgs/%s/repositories?page=%d", s, p)
		}
	case !isOrg:
		total = TotalPagesUser(entityID)
		fpath = fmt.Sprintf("users/%s.json", entityID)
		urlCreate = func(s string, p int) string {
			return fmt.Sprintf("https://github.com/%s?tab=repositories&page=%d", s, p)
		}
	}

	// make chan of total available repos
	// 1 page has 30 repos
	reposCh := make(chan *Repo, total*30)
	var wg sync.WaitGroup

	wg.Add(total)
	for p := 1; p <= total; p++ {
		// generate url per page
		url := urlCreate(entityID, p)
		ProcessPage(isOrg, url, entityID, reposCh, &wg)
	}

	wg.Wait()
	close(reposCh)

	// create file after closing channel
	CreateFile(fpath, reposCh)
}
