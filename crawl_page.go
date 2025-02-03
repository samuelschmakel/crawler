package main

import (
	"fmt"
	"net/url"
)

func crawlPage(rawBaseURL, rawCurrentURL string, pages map[string]int) {
	baseURL, err := url.Parse(rawBaseURL)
	if err != nil {
		fmt.Printf("Error - crawlPage: couldn't parse URL '%s': %v\n", rawBaseURL, err)
		return
	}

	currentURL, err := url.Parse(rawCurrentURL)
	if err != nil {
		fmt.Printf("Error - crawlPage: couldn't parse URL '%s': %v\n", rawCurrentURL, err)
		return
	}

	// only crawl the current website
	if currentURL.Hostname() != baseURL.Hostname() {
		return
	}

	normalizedCurrentURL, err := normalizeURL(rawCurrentURL)
	if err != nil {
		fmt.Printf("Error - normalizedURL: %v", err)
		return
	}

	// Check if the current URL exists in the map already (if it's already been crawled)
	_, exists := pages[normalizedCurrentURL]
	if exists{
		pages[normalizedCurrentURL] += 1
		return
	}

	// The current URL is new - create an entry to the map and set the count to 1
	pages[normalizedCurrentURL] = 1

	fmt.Printf("crawling: %s\n", rawCurrentURL)

	htmlBody, err := getHTML(rawCurrentURL)
	if err != nil {
		fmt.Printf("Error getting HTML for %s: %v\n", normalizedCurrentURL, err)
		return
	}

	// Print the HTML to watch the crawler in real-time. Change this to a subset for larger HTML documents
	fmt.Println(htmlBody)

	// Get the URLs from the HTML response body:
	nextURLs, err := getURLsFromHTML(htmlBody, rawBaseURL)
	if err != nil {
		fmt.Printf("Error - getURLsFromHTML: %v", err)
		return
	}

	for _, nextURL := range nextURLs {
		crawlPage(rawBaseURL, nextURL, pages)
	}
}