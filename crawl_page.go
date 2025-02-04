package main

import (
	"fmt"
	"net/url"
)

func (cfg *config) crawlPage(rawCurrentURL string) {
	cfg.concurrencyControl <- struct{}{}
	defer func() { <-cfg.concurrencyControl }()
	defer cfg.wg.Done()

	currentURL, err := url.Parse(rawCurrentURL)
	if err != nil {
		fmt.Printf("Error - crawlPage: couldn't parse URL '%s': %v\n", rawCurrentURL, err)
		return
	}

	// only crawl the current website
	if currentURL.Hostname() != cfg.baseURL.Hostname() {
		return
	}

	normalizedCurrentURL, err := normalizeURL(rawCurrentURL)
	if err != nil {
		fmt.Printf("Error - normalizedURL: %v", err)
		return
	}

	isFirst := cfg.addPageVisit(normalizedCurrentURL)
	if !isFirst {
		return
	}

	fmt.Printf("crawling: %s\n", rawCurrentURL)

	htmlBody, err := getHTML(rawCurrentURL)
	if err != nil {
		fmt.Printf("Error getting HTML for %s: %v\n", normalizedCurrentURL, err)
		return
	}

	// Print the HTML to watch the crawler in real-time. Change this to a subset for larger HTML documents
	fmt.Println(htmlBody)

	// Get the URLs from the HTML response body:
	nextURLs, err := getURLsFromHTML(htmlBody, cfg.baseURL.String())
	if err != nil {
		fmt.Printf("Error - getURLsFromHTML: %v", err)
		return
	}

	for _, nextURL := range nextURLs {
		cfg.wg.Add(1)
		go cfg.crawlPage(nextURL)
	}
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