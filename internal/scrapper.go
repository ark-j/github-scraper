package internal

import (
	"fmt"
	"log"
	"sync"
)

const (
	PerPageRepos = 30
	OrgURL       = "https://github.com/orgs/%s/repositories?page=%d&q=&type=%s&language=%s&sort=%s"
	UserURL      = "https://github.com/%s?tab=repositories&page=%d&q=&type=%s&language=%s&sort=%s"
)

type Scrapper struct {
	log     *log.Logger
	request *Reqwest
	reposCh chan *Repo
}

func NewScrapper(log *log.Logger, request *Reqwest) *Scrapper {
	return &Scrapper{
		log:     log,
		request: request,
		reposCh: make(chan *Repo, PerPageRepos),
	}
}

// Scrape scrapes the data of repos based on entity id provided,
// it first require you to tell if entity is orgnization or not,
// after that it gets the total pages and process perpage. finally
// appending all data to repos chan
func Scrape(isOrg bool, entityID string, f *Filter) {
	var urlCreate func(s string, p int) string
	var total int
	var fpath string
	switch {
	case isOrg:
		total = TotalPages(entityID, f)
		fpath = fmt.Sprintf("orgs/%s.json", entityID)
		urlCreate = func(s string, p int) string {
			return fmt.Sprintf("https://github.com/orgs/%s/repositories?page=%d&q=&type=%s&language=%s&sort=%s", s, p, f.Type, f.Lang, f.Sort)
		}
	case !isOrg:
		total = TotalPagesUser(entityID, f)
		fpath = fmt.Sprintf("users/%s.json", entityID)
		urlCreate = func(s string, p int) string {
			return fmt.Sprintf("https://github.com/%s?tab=repositories&page=%d&q=&type=%s&language=%s&sort=%s", s, p, f.Type, f.Lang, f.Sort)
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
