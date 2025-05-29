package main

import "fmt"

const (
	RED    = "\x1b[31m"
	GREEN  = "\x1b[32m"
	YELLOW = "\x1b[33m"
	GREY   = "\x1b[90m"
	RESET  = "\x1b[0m"
)

func red(text string) string {
	return fmt.Sprintf("%s%s%s", RED, text, RESET)
}

func green(text string) string {
	return fmt.Sprintf("%s%s%s", GREEN, text, RESET)
}

func yellow(text string) string {
	return fmt.Sprintf("%s%s%s", YELLOW, text, RESET)
}

func grey(text string) string {
	return fmt.Sprintf("%s%s%s", GREY, text, RESET)
}
