package main

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
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
	toVisitQueue  []string
	toVisitMutex  *sync.Mutex
	maxGoroutines int
	maxPages      int
	wg            *sync.WaitGroup
	startURL      *url.URL
	sameDomain    bool
	savePages     bool
}

func NewCrawler(startingURL string, maxGoroutines int, maxPages int, sameDomain bool, savePages bool) (*Crawler, error) {
	startingURLStruct, err := url.Parse(startingURL)
	if err != nil {
		return nil, err
	}

	return &Crawler{
		visited:       make(map[string]int),
		vistedMutex:   &sync.Mutex{},
		toVisitQueue:  []string{},
		toVisitMutex:  &sync.Mutex{},
		maxGoroutines: maxGoroutines,
		maxPages:      maxPages,
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
					nextURL := c.popleftURL()
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

	if c.pagesVisited() >= c.maxPages {
		// Skip if reached maxPages crawled
		return
	}

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

	// dont crawl again if already visited
	if firstVisit := c.addPageVisit(normalizedURL); !firstVisit {
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
		err := c.savePageToDisk(html, currURL)
		if err != nil {
			log.Printf(yellow("Error: %v"), err)
		}
	}
}

func (c *Crawler) addPageVisit(normalizedURL string) (firstVisit bool) {
	c.vistedMutex.Lock()
	defer c.vistedMutex.Unlock()

	if _, ok := c.visited[normalizedURL]; !ok {
		c.visited[normalizedURL] = 1
		return true
	}

	c.visited[normalizedURL] += 1
	return false
}

func (c *Crawler) pagesVisited() int {
	c.vistedMutex.Lock()
	defer c.vistedMutex.Unlock()

	return len(c.visited)
}

func (c *Crawler) savePageToDisk(html string, currURL *url.URL) error {
	hostDir := strings.ReplaceAll(currURL.Host, ".", "_")
	documentDir := filepath.Join(CORPUS_DIR, hostDir)

	_, err := os.Stat(documentDir)
	if errors.Is(err, fs.ErrNotExist) {
		err := os.Mkdir(documentDir, 0o750)
		if err != nil {
			return err
		}
	}

	hash := md5.Sum([]byte(html))
	fileName := fmt.Sprintf("%x.json", hash)
	filePath := filepath.Join(documentDir, fileName)

	type Document struct {
		Url     string `json:"url"`
		Content string `json:"content"`
	}

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "  ")

	return encoder.Encode(Document{
		Url:     currURL.String(),
		Content: html,
	})
}

func (c *Crawler) appendURL(url string) {
	c.toVisitMutex.Lock()
	defer c.toVisitMutex.Unlock()

	c.toVisitQueue = append(c.toVisitQueue, url)
}

func (c *Crawler) popleftURL() string {
	c.toVisitMutex.Lock()
	defer c.toVisitMutex.Unlock()

	size := len(c.toVisitQueue)
	if size == 0 {
		return ""
	}

	// DEPTH FIRST SEARCH
	// url := c.toVisit[size-1]
	// c.toVisit = c.toVisit[:size-1]

	// BREADTH FIRST SEARCH
	url := c.toVisitQueue[0]
	c.toVisitQueue = c.toVisitQueue[1:]
	return url
}
