package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	var startURL string
	var maxGoroutines int
	var sameDomain bool

	flag.StringVar(&startURL, "startURL", "", "the URL used to start the crawl")
	flag.IntVar(&maxGoroutines, "maxGoroutines", 3, "max number of goroutines to spawn")
	flag.BoolVar(&sameDomain, "sameDomain", false, "limit crawling to pages with same domain as startURL")

	flag.Parse()

	if startURL == "" {
		fmt.Fprint(os.Stderr, red("-startURL is required\n\n"))
		fmt.Fprintf(os.Stderr, "Usage:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\n")
		os.Exit(1)
	}

	crawler, err := NewCrawler(startURL, maxGoroutines, sameDomain)
	if err != nil {
		fmt.Fprintf(os.Stderr, red("NewCrawler: %s"), err)
	}

	log.Printf(green(`--- Starting crawl at "%s" ---`), startURL)
	crawler.StartCrawl()

	printReport(crawler.visited, startURL)
}
