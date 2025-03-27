package main

import (
	"fmt"
	"net/url"
)

func crawlPage(rawBaseURL, rawCurrentURL string, pages map[string]int) {
	baseURL, err1 := url.Parse(rawBaseURL)
	currentURL, err2 := url.Parse(rawCurrentURL)
	if err1 != nil || err2 != nil {
		fmt.Printf("Error - crawlPage: couldn't parse URL '%s'\n", rawBaseURL)
		return
	}

	// skip other domains
	if  baseURL.Hostname() != currentURL.Hostname(){
		return
	}

	normalizedCurrentURL, err := normalizeURL(rawCurrentURL)
	if err != nil {
		return
	}

	// dont crawl again if already visited
	if _, visited := pages[normalizedCurrentURL]; visited {
		pages[normalizedCurrentURL] += 1
		return
	}
	pages[normalizedCurrentURL] = 1

	fmt.Printf("Crawling: \"%s\"\n", rawCurrentURL)

	html, err := getHTML(rawCurrentURL)
	if err != nil {
		fmt.Printf("Error - getHTML: %v\n", err)
		return
	}

	urls, err := getURLsFromHTML(html, rawBaseURL)
	if err != nil {
		fmt.Printf("Error - getURLsFromHTML: %v\n", err)
		return
	}

	for _, url := range urls {
		crawlPage(rawBaseURL, url, pages)
	}
}
