package main

import (
	"fmt"
	"os"
)

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
	fmt.Printf("---Starting crawl of \"%s\"---\n", rawBaseURL)

	pages := make(map[string]int)
	crawlPage(rawBaseURL, rawBaseURL, pages)
	fmt.Println("---Done crawling---")
	for k, v := range pages {
		fmt.Println(k, v)
	}
}
