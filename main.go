package main

import (
	"fmt"
	"os"
)

func main() {
	args := os.Args
	if len(args) < 2 {
		fmt.Println("no website provided")
		os.Exit(1)
	}

	if len(args) > 2 {
		fmt.Println("too many arguments provided")
		os.Exit(1)
	}
	rawBaseURL := os.Args[1]

	fmt.Printf("starting crawl of: %s\n", rawBaseURL)

	html, err := getHTML(rawBaseURL)
	if err != nil {
		fmt.Printf("error getting HTML")
		os.Exit(1)
	}

	fmt.Println(html)
}