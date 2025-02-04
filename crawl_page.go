package main

import (
	"fmt"
	"net/url"
)

func (cfg *config) crawlPage(rawCurrentURL string) {
	cfg.concurrencyControl <- struct{}{}
	defer func() { 
		<-cfg.concurrencyControl 
		cfg.wg.Done()
		}()

	if cfg.pagesLen() >=  cfg.maxPages{
		return
	}

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