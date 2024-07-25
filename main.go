package main

import (
	"context"
	"flag"
	"fmt"
	"githubscrape/internal"
	"log"
	"os"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := os.MkdirAll("./orgs", os.ModePerm); err != nil {
		log.Println(err)
		return
	}

	if err := os.MkdirAll("./users", os.ModePerm); err != nil {
		log.Println(err)
		return
	}

	typef := flag.String("type", "", "type filter for repositories\npossible values -> source, forks, archived, mirror, template")
	langf := flag.String("lang", "", "language filter for repositories\npossible values -> go, html, javascript, java, rust,\npython, typescript, css, haskell, shell, c++, c, ruby")
	sortf := flag.String("sort", "", "sort filter for repositories\npossible value -> name, stargazers\nleave empty for last updated sort")
	user := flag.String("user", "", "github username for scraping information")
	org := flag.String("org", "", "github orgname for scraping information")
	save := flag.String("save", "json", "saver either csv or json")
	flag.Parse()

	var format int

	f := &internal.Filter{Type: *typef, Lang: *langf, Sort: *sortf}
	logger := log.New(os.Stderr, "", log.Ltime|log.LstdFlags|log.Lshortfile)
	reqwest := internal.NewReqwest()
	if *save == "csv" {
		format = 1
	}
	scrapper := internal.NewScrapper(logger, reqwest, internal.SaveFormat(format))
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
