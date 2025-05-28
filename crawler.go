package main

import (
	"log"
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
}

func NewCrawler(startingURL string, maxGoroutines int) *Crawler {
	return &Crawler{
		visited:       make(map[string]int),
		vistedMutex:   &sync.Mutex{},
		toVisit:       []string{startingURL},
		toVisitMutex:  &sync.Mutex{},
		maxGoroutines: maxGoroutines,
		wg:            &sync.WaitGroup{},
	}
}

func (c *Crawler) StartCrawl() {
	for i := range c.maxGoroutines {
		go func(id int) {
			for {
				nextURL := c.popURL()
				if nextURL == "" {
					// Wait for other goroutines to add URLs to stack
					time.Sleep(time.Second)
					continue
				}

				c.crawlPage(nextURL, id)
			}
		}(i)
	}

	c.wg.Add(1)
	c.wg.Wait()
}

func (c *Crawler) crawlPage(currURL string, id int) {
	defer c.wg.Done()

	if c.pagesVisited() >= 10 {
		return
	}

	log.Printf(`Goroutine %d: "%s"`, id, currURL)

	// // skip other domains
	// if c.startURLStruct.Hostname() != currentURL.Hostname() {
	// 	return
	// }

	currURLNormalized, err := normalizeURL(currURL)
	if err != nil {
		log.Printf("Error: %s", err)
		return
	}

	// dont crawl again if already visited
	if visited := c.addPageVisit(currURLNormalized); visited {
		return
	}

	html, err := getHTML(currURL)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	urls, err := getURLsFromHTML(html, currURL)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	for _, url := range urls {
		c.wg.Add(1)
		c.appendURL(url)
	}
}

func (c *Crawler) addPageVisit(normalizedURL string) (visited bool) {
	c.vistedMutex.Lock()
	defer c.vistedMutex.Unlock()

	if _, visited := c.visited[normalizedURL]; visited {
		c.visited[normalizedURL] += 1
		return true
	}
	c.visited[normalizedURL] = 1
	return false
}

func (c *Crawler) pagesVisited() int {
	c.vistedMutex.Lock()
	defer c.vistedMutex.Unlock()
	return len(c.visited)
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
