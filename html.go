package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func getHTML(parsedURL *url.URL) (string, error) {
	resp, err := http.Get(parsedURL.String())
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return "", fmt.Errorf(`"%s" returned status code: %d`, parsedURL.String(), resp.StatusCode)
	}
	if contentType := resp.Header.Get("content-type"); !strings.HasPrefix(contentType, "text/html") {
		return "", fmt.Errorf(`"%s" returned content of type: "%s"`, parsedURL.String(), contentType)
	}

	html, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(html), nil
}
