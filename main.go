package main

import (
	"fmt"
	"os"
)

const maxGoroutines = 5

func main() {
	args := os.Args[1:]
	if len(args) < 1 {
		fmt.Println("no website provided")
		os.Exit(1)
	}
	if len(args) > 1 {
		fmt.Println("too many arguments provided")
		os.Exit(1)
	}

	rawBaseURL := args[0]

	c, err := newCrawler(rawBaseURL, maxGoroutines)
	if err != nil {
		fmt.Printf("Error - newCrawler: %v", err)
		os.Exit(1)
	}

	fmt.Printf("---Starting crawl of \"%s\"---\n", rawBaseURL)
	c.wg.Add(1)
	go c.crawlPage(rawBaseURL)
	c.wg.Wait()

	fmt.Println("---Done crawling---")
	for normalizedURL, count := range c.pages {
		fmt.Println(normalizedURL, count)
	}
}
