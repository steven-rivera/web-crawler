package main

import (
	"fmt"
	"os"
	"strconv"
)

func main() {
	args := os.Args[1:]

	if len(args) != 3 {
		switch len(args) {
		case 0:
			fmt.Println("no baseURL provided")
		case 1:
			fmt.Println("maxGoroutines not provided")
		case 2:
			fmt.Println("maxPages not provided")
		default:
			fmt.Println("too many arguments provided")
		}
		fmt.Println("usage: crawler <baseURL> <maxGoroutines> <maxPages>")
		os.Exit(1)
	}

	rawBaseURL := args[0]
	maxGoroutines, err := strconv.Atoi(args[1])
	if err != nil {
		fmt.Println("maxGoroutines must be an integer")
		os.Exit(1)
	}
	maxPages, err := strconv.Atoi(args[2])
	if err != nil {
		fmt.Println("maxPages must be an integer")
		os.Exit(1)
	}

	c, err := newCrawler(rawBaseURL, maxGoroutines, maxPages)
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
