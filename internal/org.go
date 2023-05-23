package internal

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/PuerkitoBio/goquery"
)

// get total pages of repositories
func TotalPages(orgName string, f *Filter) int {
	rootURL := fmt.Sprintf("https://github.com/orgs/%s/repositories?q=&type=%s&language=%s&sort=%s", orgName, f.Type, f.Lang, f.Sort)
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
