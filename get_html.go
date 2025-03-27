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
		return "", fmt.Errorf("got http status code: %d", resp.StatusCode)
	}
	if contentType := resp.Header.Get("content-type");  !strings.HasPrefix(contentType, "text/html") {
		return "", fmt.Errorf("got content of type: '%s'", contentType)
	}

	html, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(html), nil
}