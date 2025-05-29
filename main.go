package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
)

const (
	SAVE_PAGES_DIR   = "PAGES"
	REPORT_FILE_NAME = "CRAWL_REPORT.txt"
)

func main() {
	var startURL string
	var maxGoroutines int
	var sameDomain bool
	var savePages bool
	var deletePrevPages bool

	flag.StringVar(&startURL, "startURL", "", "the URL used to start the crawl")
	flag.IntVar(&maxGoroutines, "maxGoroutines", 3, "max number of goroutines to spawn")
	flag.BoolVar(&sameDomain, "sameDomain", false, "limit crawling to pages with same domain as startURL")
	flag.BoolVar(&savePages, "savePages", false, fmt.Sprintf("save crawled pages to ./%s/", SAVE_PAGES_DIR))
	flag.BoolVar(&deletePrevPages, "deletePrevPages", false, fmt.Sprintf("delete pages from previous crawl in ./%s/", SAVE_PAGES_DIR))

	flag.Parse()

	if deletePrevPages {
		os.RemoveAll(SAVE_PAGES_DIR)
	}

	err := os.Mkdir(SAVE_PAGES_DIR, 0750)
	if err != nil && !errors.Is(err, os.ErrExist) {
		fmt.Fprint(os.Stderr, red("unable to create ./%s/ directory"), SAVE_PAGES_DIR)
		os.Exit(1)
	}

	if startURL == "" {
		fmt.Fprint(os.Stderr, red("-startURL is required\n\n"))
		fmt.Fprintf(os.Stderr, "Usage:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\n")
		os.Exit(1)
	}

	crawler, err := NewCrawler(startURL, maxGoroutines, sameDomain, savePages)
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
