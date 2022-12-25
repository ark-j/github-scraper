package main

import (
	"flag"
	"fmt"
	"githubscrape"
	"os"
)

func main() {
	os.MkdirAll("./orgs", os.ModePerm)
	os.MkdirAll("./users", os.ModePerm)

	user := flag.String("user", "", "github username for scraping information")
	org := flag.String("org", "", "github orgname for scraping information")
	flag.Parse()

	if *org == "" && *user == "" {
		fmt.Println("please provide any user or org")
	}
	if *user != "" {
		githubscrape.Scrape(false, *user)
	}
	if *org != "" {
		githubscrape.Scrape(true, *org)
	}
}
