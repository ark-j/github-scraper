package main

import "githubscrape"

func main() {
	// TODO: add github repo make changes -> push -> track
	//
	// TODO: refactor code remove duplicate
	//
	// TODO: add command line flag for parsing filter, orgnization, user, at once

	githubscrape.Scrape("uber")
	githubscrape.Scrape("vuejs")
	githubscrape.Scrape("slidevjs")
	githubscrape.Scrape("go-chi")
	githubscrape.Scrape("redis")
	githubscrape.ScrapeUser("spf13")
	githubscrape.ScrapeUser("graydon")
	githubscrape.ScrapeUser("Sajmani")
	githubscrape.Scrape("google")
}
