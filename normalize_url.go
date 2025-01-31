package main

import (
	"fmt"
	"net/url"
	"strings"
)

func normalizeURL(URL string) (string, error) {
	parsed, err := url.Parse(URL)
	if err != nil {
		return "", fmt.Errorf("couldn't parse URL: %w", err)
	}

	parsedPath := parsed.Host + parsed.Path
	parsedPath = strings.ToLower(parsedPath)
	parsedPath = strings.TrimSuffix(parsedPath, "/")

	return parsedPath, nil
}