package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

func getHTML(rawURL string) (string, error) {
	resp, err := http.Get(rawURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return "", fmt.Errorf(`"%s" returned status code: %d`, rawURL, resp.StatusCode)
	}
	if contentType := resp.Header.Get("content-type"); !strings.HasPrefix(contentType, "text/html") {
		return "", fmt.Errorf(`"%s" returned content of type: "%s"`, rawURL, contentType)
	}

	html, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(html), nil
}
