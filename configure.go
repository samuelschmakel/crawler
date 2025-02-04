package main

import (
	"fmt"
	"net/url"
	"sync"
)

type config struct {
	pages map[string]int
	baseURL *url.URL
	mu *sync.Mutex
	concurrencyControl chan struct{}
	wg *sync.WaitGroup
}

func (cfg *config) addPageVisit(normalizedURL string) (isFirst bool) {
	// Check if the current URL exists in the map already (if it's already been crawled)
	cfg.mu.Lock()
	defer cfg.mu.Unlock()

	_, exists := cfg.pages[normalizedURL]
	if exists{
		cfg.pages[normalizedURL] += 1
		return false
	}

	// The current URL is new - create an entry to the map and set the count to 1
	cfg.pages[normalizedURL] = 1
	return true
}

func configure(rawBaseURL string, maxConcurrency int) (*config, error) {
	baseURL, err := url.Parse(rawBaseURL)
	if err != nil {
		return nil, fmt.Errorf("couldn't parse base URL: %v", err)
	}

	return &config{
		pages: make(map[string]int),
		baseURL: baseURL,
		mu: &sync.Mutex{},
		concurrencyControl: make(chan struct{}, maxConcurrency),
		wg: &sync.WaitGroup{},
	}, nil
}