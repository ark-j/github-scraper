package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"githubscrape/internal"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	os.MkdirAll("./orgs", os.ModePerm)
	os.MkdirAll("./users", os.ModePerm)

	typef := flag.String("type", "", "type filter for repositories\npossible values -> source, forks, archived, mirror, template")
	langf := flag.String("lang", "", "language filter for repositories\npossible values -> go, html, javascript, java, rust,\npython, typescript, css, haskell, shell, c++, c, ruby")
	sortf := flag.String("sort", "", "sort filter for repositories\npossible value -> name, stargazers\nleave empty for last updated sort")
	user := flag.String("user", "", "github username for scraping information")
	org := flag.String("org", "", "github orgname for scraping information")
	flag.Parse()

	f := &internal.Filter{Type: *typef, Lang: *langf, Sort: *sortf}
	logger := log.New(os.Stderr, "", log.Ltime|log.LstdFlags|log.Lshortfile)
	reqwest := internal.NewReqwest()
	scrapper := internal.NewScrapper(logger, reqwest)
	if *org == "" && *user == "" {
		fmt.Println("please provide any user or org")
	}
	if *user != "" {
		scrapper.Scrape(ctx, false, *user, f)
	}
	if *org != "" {
		scrapper.Scrape(ctx, true, *org, f)
	}
}
