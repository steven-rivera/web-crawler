package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	var startURL string
	var maxGoroutines int
	var maxPages int

	flag.StringVar(&startURL, "url", "", "the URL used to start the crawl")
	flag.IntVar(&maxGoroutines, "maxGoroutines", 1, "max number of goroutines to spawn")
	flag.IntVar(&maxPages, "maxPages", 1, "max number of pages to crawl")

	flag.Parse()

	if startURL == "" {
		fmt.Fprintln(os.Stderr, "Error: -startURL is required")
		flag.PrintDefaults()
		os.Exit(1)
	}

	crawler := NewCrawler(startURL, maxGoroutines)

	fmt.Printf("---Starting crawl of \"%s\"---\n", startURL)
	crawler.StartCrawl()

	printReport(crawler.visited, startURL)
}
