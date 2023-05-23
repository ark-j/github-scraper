package internal

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/PuerkitoBio/goquery"
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

// concurrently process per page for user or org
// of user then org should be false
func ProcessPage(isOrg bool, url string, entity string, ch chan<- *Repo, wg *sync.WaitGroup) {
	defer wg.Done()
	var id string
	switch isOrg {
	case false:
		id = "#user-repositories-list"
	case true:
		id = "#org-repositories"
	}
	res, err := http.Get(url)
	if err != nil {
		log.Println("ERROR", "msg=not able to get request", "error=", err)
	}
	defer res.Body.Close()
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Println("ERROR", "msg=not able to parse body", "error=", err)
	}
	// all the repos are in unorderd list
	selection := doc.Find(id).Find("ul").Find("li")
	selection.Each(ProcessRepo(entity, ch))
}

// process repo data for single repo found
func ProcessRepo(entity string, ch chan<- *Repo) func(i int, s *goquery.Selection) {
	return func(i int, s *goquery.Selection) {
		baseName := s.Find("a[itemprop='name codeRepository']")
		title := ClearString(baseName.Text())
		link, _ := baseName.Attr("href")
		description := ClearString(s.Find("p[itemprop='description']").Text())
		language := s.Find("span[itemprop='programmingLanguage']").Text()
		forks := ClearString(s.Find(fmt.Sprintf("a[href='/%s/%s/network/members']", entity, title)).Text())
		stars := ClearString(s.Find(fmt.Sprintf("a[href='/%s/%s/stargazers']", entity, title)).Text())
		ch <- &Repo{
			Title:       title,
			Link:        URL + link,
			Description: description,
			Language:    language,
			Forks:       forks,
			Stars:       stars,
		}
	}
}
