package main

import (
	"fmt"
	"net/url"
	"sort"
)

type Page struct {
	URL  string
	Count int
}

func printReport(pages map[string]int, baseURL string) {
	fmt.Printf(`
=============================
  REPORT for %s
=============================
`, baseURL)

	// Sort map into slice of structs
	sortedPages := sortPages(pages)

	parsed, err := url.Parse(baseURL)
	if err != nil {
		fmt.Println("couldn't parse baseURL to print")
		return
	}
	protocol := parsed.Scheme

	for _, page := range sortedPages {
		fmt.Printf("Found %d internal links to %s://%s\n", page.Count, protocol, page.URL)
	}

}

func sortPages(pages map[string]int) []Page {
	pageStructs := []Page{}
	for url, count := range pages {
		pageStructs = append(pageStructs, Page{URL: url, Count: count})
	}

	sort.Slice(pageStructs, func(i, j int) bool {
		if pageStructs[i].Count == pageStructs[j].Count {
			return pageStructs[i].URL < pageStructs[j].URL
		}
		return pageStructs[i].Count > pageStructs[j].Count
	})

	return pageStructs
}