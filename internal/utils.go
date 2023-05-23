package internal

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

const URL = "https://github.com"

func CreateFilter(typef, langf, sortf, entity string, org bool) string {
	if org {
		return fmt.Sprintf("https://github.com/orgs/%s/repositories?q=&type=%s&language=%s&sort=%s", entity, typef, langf, sortf)
	}
	return fmt.Sprintf("https://github.com/%s?tab=repositories&q=&type=%s&language=%s&sort=%s", entity, typef, langf, sortf)
}

// creates json file for per org
func CreateFile(path string, ch <-chan *Repo) {
	f, err := os.Create(path)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()
	var ll []*Repo
	for i := range ch {
		ll = append(ll, i)
	}
	b, _ := json.MarshalIndent(map[string]any{"count": len(ll), "repos": ll}, "", "  ")
	f.Write(b)
}

// cleans the string
func ClearString(s string) string {
	return strings.Trim(
		strings.ReplaceAll(
			s,
			"\n",
			"",
		),
		" ",
	)
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
		log.Println(err)
	}
	defer res.Body.Close()
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Println(err)
	}
	selection := doc.Find(id).Find("ul").Find("li")
	selection.Each(ProcessRepo(entity, ch))
}
