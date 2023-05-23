package internal

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

const (
	OrgURL  = "https://github.com/orgs/%s/repositories?page=%d&q=&type=%s&language=%s&sort=%s"
	UserURL = "https://github.com/%s?tab=repositories&page=%d&q=&type=%s&language=%s&sort=%s"
	BaseURL = "https://github.com"
)

type Scrapper struct {
	log     *log.Logger
	request *Reqwest
}

func NewScrapper(log *log.Logger, request *Reqwest) *Scrapper {
	return &Scrapper{
		log:     log,
		request: request,
	}
}

// Scrape scrapes the data of repos based on entity id provided,
// it first require you to tell if entity is orgnization or not,
// after that it gets the total pages and process perpage. finally
// appending all data to repos chan
func (sc *Scrapper) Scrape(ctx context.Context, isOrg bool, entityID string, f *Filter) {
	var urlCreate func(s string, p int) string
	var total int
	var fpath string
	switch {
	case isOrg:
		total = sc.TotalPagesOrg(ctx, entityID, f)
		fpath = fmt.Sprintf("orgs/%s.json", entityID)
		urlCreate = func(s string, p int) string {
			return fmt.Sprintf("https://github.com/orgs/%s/repositories?page=%d&q=&type=%s&language=%s&sort=%s", s, p, f.Type, f.Lang, f.Sort)
		}
	case !isOrg:
		total = sc.TotalPagesUser(ctx, entityID, f)
		fpath = fmt.Sprintf("users/%s.json", entityID)
		urlCreate = func(s string, p int) string {
			return fmt.Sprintf("https://github.com/%s?tab=repositories&page=%d&q=&type=%s&language=%s&sort=%s", s, p, f.Type, f.Lang, f.Sort)
		}
	}

	// make chan of total available repos
	// 1 page has 30 repos
	reposCh := make(chan *Repo)
	var wg sync.WaitGroup
	wg.Add(total)
	for p := 1; p <= total; p++ {
		// generate url per page
		url := urlCreate(entityID, p)
		go sc.ProcessPage(ctx, isOrg, url, entityID, reposCh, &wg)
	}

	go func() {
		wg.Wait()
		close(reposCh)
	}()

	// create file after closing channel
	CreateFile(fpath, reposCh)
}

// get total pages of repositories
func (sc *Scrapper) TotalPagesUser(ctx context.Context, userID string, f *Filter) int {
	counter := 1
	stopper := true
	for stopper {
		// start from root url
		// by increasing counter
		// so we can crawl next pages until none
		rootURL := fmt.Sprintf("https://github.com/%s?tab=repositories&page=%d&q=&type=%s&language=%s&sort=%s", userID, counter, f.Type, f.Lang, f.Sort)
		doc, err := sc.request.Source(ctx, rootURL)
		if err != nil {
			log.Println("msg=not able to parse body", "error=", err)
			break
		}
		_, ok := doc.Find("#user-repositories-list").Find("div.paginate-container").Find("a.next_page").Attr("href")
		if ok {
			stopper = true
			counter += 1
		} else {
			stopper = false
		}
	}
	sc.log.Printf("INFO user_id=%s pages=%d\n", userID, counter)
	return counter
}

// get total pages of repositories
func (sc *Scrapper) TotalPagesOrg(ctx context.Context, orgName string, f *Filter) int {
	pageCount := 1
	rootURL := fmt.Sprintf("https://github.com/orgs/%s/repositories?q=&type=%s&language=%s&sort=%s", orgName, f.Type, f.Lang, f.Sort)
	doc, err := sc.request.Source(ctx, rootURL)
	if err != nil {
		sc.log.Println("msg=not able to parse body", "error=", err)
	}
	pages, ok := doc.Find("#org-repositories").Find("div.pagination").Find("em.current").Attr("data-total-pages")
	if ok {
		pageCount, _ = strconv.Atoi(pages)
	}
	sc.log.Printf("INFO org_name=%s pages=%d\n", orgName, pageCount)
	return pageCount
}

// concurrently process per page for user or org
// of user then org should be false
func (sc *Scrapper) ProcessPage(ctx context.Context, isOrg bool, URL string, entity string, ch chan<- *Repo, wg *sync.WaitGroup) {
	defer wg.Done()
	var id string
	switch isOrg {
	case false:
		id = "#user-repositories-list"
	case true:
		id = "#org-repositories"
	}
	doc, err := sc.request.Source(ctx, URL)
	if err != nil {
		sc.log.Println("ERROR", "error=", err)
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
			Link:        BaseURL + link,
			Description: description,
			Language:    language,
			Forks:       forks,
			Stars:       stars,
		}
	}
}
