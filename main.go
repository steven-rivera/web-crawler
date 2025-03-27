package main

import (
	"fmt"
	"os"
)

func main() {
	args := os.Args[1:]
	if len(args) < 1 {
		fmt.Println("no website provided")
		os.Exit(1)
	}
	if len(args) > 1 {
		fmt.Println("too many arguments provided")
		os.Exit(1)
	}

	rawBaseURL := args[0]
	fmt.Printf("starting crawl of: %s\n", rawBaseURL)
	
	html, err := getHTML(rawBaseURL)
	if err != nil {
		fmt.Printf("Error getting html: %s\n", err)
		os.Exit(1)
	}
	fmt.Print(html)
}
