package main

import (
	"fmt"
	"os"
	"strconv"
)

func main() {
	args := os.Args
	if len(args) < 2 {
		fmt.Println("no website provided")
		os.Exit(1)
	}

	if len(args) > 5 {
		fmt.Println("too many arguments provided")
		os.Exit(1)
	}

	rawBaseURL := args[1]

	// Default values for optional arguments
	maxConcurrency := 5
	maxPages := 10

	if len(args) > 2 {
		if val, err := strconv.Atoi(args[2]); err == nil {
			maxConcurrency = val
		} else {
			fmt.Println("Error: maxConcurrency must be an integer")
			os.Exit(1)
		}
	}

	if len(args) > 3 {
		if val, err := strconv.Atoi(args[3]); err == nil {
			maxPages = val
		} else {
			fmt.Println("Error: maxPages must be an integer")
			os.Exit(1)
		}
	}

	cfg, err := configure(rawBaseURL, maxConcurrency, maxPages)
	if err != nil {
		fmt.Printf("couldn't configure: %v", err)
		return
	}

	fmt.Printf("starting crawl of: %s\n", rawBaseURL)

	cfg.wg.Add(1)

	go cfg.crawlPage(rawBaseURL)
	cfg.wg.Wait()

	printReport(cfg.pages, rawBaseURL)
}