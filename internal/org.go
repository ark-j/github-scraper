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
	pageCount := 1
	rootURL := fmt.Sprintf("https://github.com/orgs/%s/repositories?q=&type=%s&language=%s&sort=%s", orgName, f.Type, f.Lang, f.Sort)
	res, err := http.Get(rootURL)
	if err != nil {
		log.Println("msg=not able to get request", "error=", err)
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Println("msg=not able to parse body", "error=", err)
	}

	pages, ok := doc.Find("#org-repositories").Find("div.pagination").Find("em.current").Attr("data-total-pages")
	if ok {
		pageCount, _ = strconv.Atoi(pages)
	}
	log.Printf("INFO org_name=%s msg=page count %d\n", orgName, pageCount)
	return pageCount
}
