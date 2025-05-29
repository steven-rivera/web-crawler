package main

import (
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type Crawler struct {
	visited       map[string]int
	vistedMutex   *sync.Mutex
	toVisit       []string
	toVisitMutex  *sync.Mutex
	maxGoroutines int
	wg            *sync.WaitGroup
	startURL      *url.URL
	sameDomain    bool
	savePages     bool
}

func NewCrawler(startingURL string, maxGoroutines int, sameDomain bool, savePages bool) (*Crawler, error) {
	startingURLStruct, err := url.Parse(startingURL)
	if err != nil {
		return nil, err
	}

	return &Crawler{
		visited:       make(map[string]int),
		vistedMutex:   &sync.Mutex{},
		toVisit:       []string{},
		toVisitMutex:  &sync.Mutex{},
		maxGoroutines: maxGoroutines,
		wg:            &sync.WaitGroup{},
		startURL:      startingURLStruct,
		sameDomain:    sameDomain,
		savePages:     savePages,
	}, nil
}

func (c *Crawler) StartCrawl() {
	ch := make(chan struct{}, c.maxGoroutines)

	for i := range c.maxGoroutines {
		go func(id int, ch chan struct{}) {
			for {
				select {
				case <-ch:
					return
				default:
					nextURL := c.popURL()
					if nextURL == "" {
						// Wait for other goroutines to add URLs to stack
						time.Sleep(time.Second)
						continue
					}

					c.crawlPage(nextURL, id)
				}

			}
		}(i+1, ch)
	}

	c.wg.Add(1)
	c.appendURL(c.startURL.String())
	c.wg.Wait()

	// Signal all goroutines to return once done crawling
	for range c.maxGoroutines {
		ch <- struct{}{}
	}
	close(ch)
}

func (c *Crawler) crawlPage(rawCurrURL string, id int) {
	defer c.wg.Done()

	currURL, err := url.Parse(rawCurrURL)
	if err != nil {
		// Skip invalid URL
		return
	}

	if c.sameDomain && c.startURL.Hostname() != currURL.Hostname() {
		// Skip if only crawling start domain
		return
	}

	normalizedURL := normalizeURL(currURL)
	defer c.addPageVisit(normalizedURL)

	// dont crawl again if already visited
	if c.visitedPage(normalizedURL) {
		return
	}

	log.Printf(grey(`Goroutine %d crawling: %s`), id, rawCurrURL)
	html, err := getHTML(currURL)
	if err != nil {
		log.Printf(yellow("Error: %v"), err)
		return
	}

	urls, err := getURLsFromHTML(html, currURL)
	if err != nil {
		log.Printf(yellow("Error: %v"), err)
		return
	}

	for _, url := range urls {
		c.wg.Add(1)
		c.appendURL(url)
	}

	if c.savePages {
		err := c.savePageToDisk(html, normalizedURL)
		if err != nil {
			log.Printf(yellow("Error: %v"), err)
		}
	}
}

func (c *Crawler) visitedPage(normalizedURL string) bool {
	c.vistedMutex.Lock()
	defer c.vistedMutex.Unlock()

	_, ok := c.visited[normalizedURL]
	return ok
}

func (c *Crawler) addPageVisit(normalizedURL string) {
	c.vistedMutex.Lock()
	defer c.vistedMutex.Unlock()

	c.visited[normalizedURL] += 1
}

func (c *Crawler) savePageToDisk(html string, fileName string) error {
	escapedFileNmae := strings.ReplaceAll(fileName, "/", "_")
	path := filepath.Join(SAVE_PAGES_DIR, escapedFileNmae)

	return os.WriteFile(path, []byte(html), 0666)
}

func (c *Crawler) appendURL(url string) {
	c.toVisitMutex.Lock()
	defer c.toVisitMutex.Unlock()

	c.toVisit = append(c.toVisit, url)
}

func (c *Crawler) popURL() string {
	c.toVisitMutex.Lock()
	defer c.toVisitMutex.Unlock()

	size := len(c.toVisit)
	if size == 0 {
		return ""
	}

	url := c.toVisit[size-1]
	c.toVisit = c.toVisit[:size-1]
	return url
}
