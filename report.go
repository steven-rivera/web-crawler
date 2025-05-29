package main

import (
	"cmp"
	"fmt"
	"os"
	"slices"
	"strings"
)

type page struct {
	count int
	url   string
}

func createReport(pages map[string]int, baseURL string) error {
	file, err := os.Create(REPORT_FILE_NAME)
	if err != nil {
		return err
	}
	defer file.Close()

	title := fmt.Sprintf("REPORT for crawl starting at %s\n", baseURL)

	fmt.Fprint(file, strings.Repeat("=", len(title)), "\n")
	fmt.Fprint(file, title)
	fmt.Fprint(file, strings.Repeat("=", len(title)), "\n\n")

	sortedPages := sortPages(pages)
	for _, page := range sortedPages {
		fmt.Fprintf(file, "Found %d links to %s\n", page.count, page.url)
	}
	return nil
}

func sortPages(pages map[string]int) []page {
	pagesSlice := make([]page, 0, len(pages))
	for url, count := range pages {
		pagesSlice = append(pagesSlice, page{count: count, url: url})
	}

	slices.SortFunc(pagesSlice, func(a, b page) int {
		if n := cmp.Compare(a.count, b.count); n != 0 {
			return -n
		}
		return strings.Compare(a.url, b.url)
	})
	return pagesSlice
}
