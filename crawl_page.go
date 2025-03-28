package main

import (
	"fmt"
	"net/url"
	"sync"
)

type crawler struct {
	pages              map[string]int
	baseURL            *url.URL
	mu                 *sync.Mutex
	concurrencyControl chan struct{}
	wg                 *sync.WaitGroup
	maxPages           int
}

func newCrawler(rawBaseURL string, maxGoroutines, maxPages int) (*crawler, error) {
	baseURL, err := url.Parse(rawBaseURL)
	if err != nil {
		return nil, fmt.Errorf("couldn't parse base URL: %v", err)
	}

	return &crawler{
		pages:              make(map[string]int),
		baseURL:            baseURL,
		mu:                 &sync.Mutex{},
		concurrencyControl: make(chan struct{}, maxGoroutines),
		wg:                 &sync.WaitGroup{},
		maxPages:           maxPages,
	}, nil
}

func (c *crawler) crawlPage(rawCurrentURL string) {
	c.concurrencyControl <- struct{}{}
	defer func() {
		<-c.concurrencyControl
		c.wg.Done()
	}()

	
	if c.pagesVisited() >= c.maxPages {
		return
	}

	currentURL, err := url.Parse(rawCurrentURL)
	if err != nil {
		fmt.Printf("Error - crawlPage: couldn't parse URL '%s': %s\n", rawCurrentURL, err)
		return
	}

	// skip other domains
	if c.baseURL.Hostname() != currentURL.Hostname() {
		return
	}

	normalizedCurrentURL, err := normalizeURL(rawCurrentURL)
	if err != nil {
		return
	}

	// dont crawl again if already visited
	if isFirst := c.addPageVisit(normalizedCurrentURL); !isFirst {
		return
	}

	fmt.Printf("Crawling: \"%s\"\n", rawCurrentURL)

	html, err := getHTML(rawCurrentURL)
	if err != nil {
		fmt.Printf("Error - getHTML: %v\n", err)
		return
	}

	urls, err := getURLsFromHTML(html, c.baseURL.String())
	if err != nil {
		fmt.Printf("Error - getURLsFromHTML: %v\n", err)
		return
	}

	for _, url := range urls {
		c.wg.Add(1)
		go c.crawlPage(url)
	}
}

func (c *crawler) addPageVisit(normalizedURL string) (isFirst bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, visited := c.pages[normalizedURL]; visited {
		c.pages[normalizedURL] += 1
		return false
	}
	c.pages[normalizedURL] = 1
	return true
}

func (c *crawler) pagesVisited() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.pages)
}
