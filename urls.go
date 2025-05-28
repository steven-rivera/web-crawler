package main

import (
	"fmt"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

/*
Strip the scheme from url, convert all characters to lowercase,
and remove trailing slashes.

Ex: http://BLOG.example.com/path/ -> blog.example.com/path
*/
func normalizeURL(parsedURL *url.URL) string {
	normal := parsedURL.Host + parsedURL.EscapedPath()
	normal = strings.ToLower(normal)
	normal = strings.TrimSuffix(normal, "/")

	return normal
}

func getURLsFromHTML(htmlBody string, parsedURL *url.URL) ([]string, error) {
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
					// Convert relative urls (/path) to absolute urls (example.com/path)
					urls = append(urls, parsedURL.ResolveReference(href).String())
					break
				}
			}
		}
	}

	return urls, nil
}
