package main

import (
	"reflect"
	"strings"
	"testing"
)

func TestGetURLsFromHTML(t *testing.T) {
	tests:= []struct {
		name string
		inputURL string
		inputBody string
		expected []string
		errorContains string
	}{
		{
			name: "absolute and relative URLs",
			inputURL: "https://blog.boot.dev",
			inputBody: `
				<html>
					<body>
						<a href="/path/one">
							<span>Boot.dev</span>
						</a>
						<a href="https://other.com/path/one">
							<span>Boot.dev</span>
						</a>
					</body>
				</html>
				`,
			expected: []string{"https://blog.boot.dev/path/one", "https://other.com/path/one"},
		},
		{
			name: "href not at front of tag",
			inputURL: "https://google.com",
			inputBody: `
				<html>
					<body>
						<a class="link" href="/path/main">
							<span>Google</span>
						</a>
					</body>
				</html>
			`,
			expected: []string{"https://google.com/path/main"},
		},
		{
			name:     "no href",
			inputURL: "https://blog.boot.dev",
			inputBody: `
					<html>
						<body>
							<a>
								<span>Boot.dev></span>
							</a>
						</body>
					</html>
					`,
			expected: nil,
		},
		{
			name:     "bad HTML",
			inputURL: "https://blog.boot.dev",
			inputBody: `
					<html body>
						<a href="path/one">
							<span>Boot.dev></span>
						</a>
					</html body>
					`,
			expected: []string{"https://blog.boot.dev/path/one"},
		},
		{
			name:     "bad HTML",
			inputURL: "https://blog.boot.dev",
			inputBody: `
					<html body>
						<a href="path/one">
							<span>Boot.dev></span>
						</a>
					</html body>
					`,
			expected: []string{"https://blog.boot.dev/path/one"},
		},
		{
			name:     "invalid href URL",
			inputURL: "https://blog.boot.dev",
			inputBody: `
					<html>
						<body>
							<a href=":\\invalidURL">
								<span>Boot.dev</span>
							</a>
						</body>
					</html>
					`,
			expected: nil,
		},
		{
			name:     "handle invalid base URL",
			inputURL: `:\\invalidBaseURL`,
			inputBody: `
					<html>
						<body>
							<a href="/path">
								<span>Boot.dev</span>
							</a>
						</body>
					</html>
					`,
			expected:      nil,
			errorContains: "couldn't parse base URL",
		},
	}

	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := getURLsFromHTML(tc.inputBody, tc.inputURL)
			if err != nil && !strings.Contains(err.Error(), tc.errorContains) {
				t.Errorf("Test %v - '%s' FAIL: unexpected error: %v", i, tc.name, err)
				return
			} else if err != nil && tc.errorContains == "" {
				t.Errorf("Test %v - '%s' FAIL: unexpected error: %v", i, tc.name, err)
				return
			} else if err == nil && tc.errorContains != "" {
				t.Errorf("Test %v - '%s' FAIL: expected error containing '%v', got none.", i, tc.name, tc.errorContains)
				return
			}

			if !reflect.DeepEqual(actual, tc.expected) {
				t.Errorf("Test %v - %s FAIL: expected URL: %v, actual: %v", i, tc.name, tc.expected, actual)
			}
		})
	}
}