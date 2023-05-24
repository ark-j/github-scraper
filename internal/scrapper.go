package internal

import (
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"os"
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
	var fpath, pageFrameID string
	switch isOrg {
	case true:
		total = sc.TotalPagesOrg(ctx, entityID)
		fpath = fmt.Sprintf("orgs/%s.json", entityID)
		pageFrameID = "#org-repositories"
		urlCreate = func(s string, p int) string {
			return fmt.Sprintf(OrgURL, s, p, f.Type, f.Lang, f.Sort)
		}
	case false:
		total = sc.TotalPagesUser(ctx, entityID)
		fpath = fmt.Sprintf("users/%s.json", entityID)
		pageFrameID = "#user-repositories-list"
		urlCreate = func(s string, p int) string {
			return fmt.Sprintf(UserURL, s, p, f.Type, f.Lang, f.Sort)
		}
	}

	// make chan of total available repos
	// 1 page has 30 repos
	reposCh := make(chan *Repo)
	var wg sync.WaitGroup
	wg.Add(total)
	for p := 1; p <= total; p++ {
		// generate url per page
		URL := urlCreate(entityID, p)
		go sc.ProcessPage(ctx, reposCh, &wg, pageFrameID, URL, entityID)
	}

	go func() {
		wg.Wait()
		close(reposCh)
	}()

	// create file after closing channel
	CreateFile(fpath, reposCh)
}

// TotalPagesUser method return count of pages present in User account
func (sc *Scrapper) TotalPagesUser(ctx context.Context, userID string) int {
	counter := 1
	for {
		// start from root url
		// by increasing counter
		// so we can crawl next pages until none
		rootURL := fmt.Sprintf("https://github.com/%s?tab=repositories&page=%d", userID, counter)
		doc, err := sc.request.Source(ctx, rootURL)
		if err != nil {
			log.Println("msg=not able to parse body", "error=", err)
			break
		}
		_, ok := doc.Find("#user-repositories-list").Find("div.paginate-container").Find("a.next_page").Attr("href")
		if !ok {
			break
		}
		counter += 1
	}
	sc.log.Printf("INFO user_id=%s pages=%d\n", userID, counter)
	return counter
}

// TotalPagesOrg method return count of pages present in orgnization account
func (sc *Scrapper) TotalPagesOrg(ctx context.Context, orgName string) int {
	pageCount := 1
	rootURL := fmt.Sprintf("https://github.com/orgs/%s/repositories", orgName)
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

// ProcessPage method finds all attribute of all repo present in one page and send it to repo chan
func (sc *Scrapper) ProcessPage(ctx context.Context, ch chan<- *Repo, wg *sync.WaitGroup, pageFrameID, URL, entity string) {
	defer wg.Done()
	doc, err := sc.request.Source(ctx, URL)
	if err != nil {
		sc.log.Println("ERROR", "error=", err)
	}
	// all the repos are in unorderd list
	doc.Find(pageFrameID).
		Find("ul").
		Find("li").
		Each(func(i int, s *goquery.Selection) {
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
		})
}

func (sc *Scrapper) SaveCSV(fpath string, ch <-chan *Repo) error {
	f, err := os.Create(fpath)
	if err != nil {
		return err
	}
	cw := csv.NewWriter(f)
	for line := range ch {
		err := cw.Write([]string{line.Title, line.Link, line.Description, line.Language, line.Forks, line.Stars})
		if err != nil {
			sc.log.Printf("ERROR msg=error while adding csv record  err=%v\n", err)
		}
	}
	cw.Flush()
	return f.Close()
}
