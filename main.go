package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
)

const (
	CORPUS_DIR             = "CORPUS"
	REPORT_FILE_NAME       = "CRAWL_REPORT.txt"
	DEFAULT_MAX_GOROUTINES = 3
	DEFAULT_MAX_PAGES      = 1000
)

func main() {
	var startURL string
	var maxGoroutines int
	var maxPages int
	var sameDomain bool
	var savePages bool
	var deletePrevPages bool

	flag.StringVar(&startURL, "startURL", "", "the URL used to start the crawl")
	flag.IntVar(&maxGoroutines, "maxGoroutines", DEFAULT_MAX_GOROUTINES, "max number of goroutines to spawn")
	flag.IntVar(&maxPages, "maxPages", DEFAULT_MAX_PAGES, "stop crawl of N pages visited")
	flag.BoolVar(&sameDomain, "sameDomain", false, "limit crawling to pages with same domain as startURL")
	flag.BoolVar(&savePages, "savePages", false, fmt.Sprintf("save crawled pages to ./%s/", CORPUS_DIR))
	flag.BoolVar(&deletePrevPages, "deletePrevPages", false, fmt.Sprintf("delete ./%s/ directory from previous crawl if exists", CORPUS_DIR))

	flag.Parse()

	if deletePrevPages {
		os.RemoveAll(CORPUS_DIR)
	}

	err := os.Mkdir(CORPUS_DIR, 0o750)
	if err != nil && !errors.Is(err, os.ErrExist) {
		fmt.Fprint(os.Stderr, red("unable to create ./%s/ directory"), CORPUS_DIR)
		os.Exit(1)
	}

	if startURL == "" {
		fmt.Fprint(os.Stderr, red("-startURL is required\n\n"))
		fmt.Fprintf(os.Stderr, "Usage:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\n")
		os.Exit(1)
	}

	crawler, err := NewCrawler(startURL, maxGoroutines, maxPages, sameDomain, savePages)
	if err != nil {
		fmt.Fprintf(os.Stderr, red("NewCrawler: %s"), err)
	}

	log.Printf(green(`--- Starting crawl at "%s" ---`), startURL)
	crawler.StartCrawl()

	log.Print(grey("Generating report..."))
	err = createReport(crawler.visited, startURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, red("Failed creating report: %s"), err)
		os.Exit(1)
	}
	log.Printf(green("Successfully created report %s"), REPORT_FILE_NAME)
}
