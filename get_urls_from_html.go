package main

import (
	"fmt"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

func getURLsFromHTML(htmlBody, rawBaseURL string) ([]string, error) {
	baseURL, err := url.Parse(rawBaseURL)
	if err != nil {
		return nil, fmt.Errorf("couldn't parse base URL: %v", err)
	}

	doc, err := html.Parse(strings.NewReader(htmlBody))
	if err != nil {
		return nil, fmt.Errorf("couldn't parse HTML: %v", err)
	}

	urls := make([]string, 0)
	for desc := range doc.Descendants() {
		if desc.Type == html.ElementNode && desc.Data == "a" {
			for _, attr := range desc.Attr {
				if attr.Key == "href" {
					href, err := url.Parse(attr.Val)
					if err != nil {
						break
					}
					resolvedURL := baseURL.ResolveReference(href)
					urls = append(urls, resolvedURL.String())
					break
				}
			}
		}
	}

	return urls, nil
}
