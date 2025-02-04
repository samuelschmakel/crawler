package main

import (
	"fmt"
	"net/url"
	"os"
	"sync"
)

type config struct {
	pages map[string]int
	baseURL *url.URL
	mu *sync.Mutex
	concurrencyControl chan struct{}
	wg *sync.WaitGroup
}

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
	pages := make(map[string]int)

	rawBaseURL := os.Args[1]
	baseURL, err := url.Parse(rawBaseURL)
	if err != nil {
		fmt.Printf("couldn't parse base URL: %v\n", err)
		os.Exit(1)
	}

	const maxConcurrent = 5
	concurrencyControl := make(chan struct{}, maxConcurrent)

	cfg := config{
		pages: pages,
		baseURL: baseURL,
		mu: &sync.Mutex{},
		concurrencyControl: concurrencyControl,
		wg: &sync.WaitGroup{},
	}

	fmt.Printf("starting crawl of: %s\n", rawBaseURL)

	cfg.wg.Add(1)
	go cfg.crawlPage(rawBaseURL)

	cfg.wg.Wait()

	for normalizedURL, count := range pages {
		fmt.Printf("%d - %s\n", count, normalizedURL)
	}
}