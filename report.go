package main

import (
	"cmp"
	"fmt"
	"slices"
	"strings"
)

type page struct {
	count int
	url   string
}

func printReport(pages map[string]int, baseURL string) {
	fmt.Println("=============================")
	fmt.Printf("REPORT for %s\n", baseURL)
	fmt.Println("=============================")

	sortedPages := sortPages(pages)

	for _, page := range sortedPages {
		fmt.Printf("Found %d internal links to %s\n", page.count, page.url)
	}
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
