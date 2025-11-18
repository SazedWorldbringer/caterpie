package main

import "testing"

// tests if normalizeURL using a URL and checks if it returns a normalized URL
func TestNormalizeURL(t *testing.T) {
	tests := []struct {
		name     string
		inputURL string
		expected string
	}{
		{
			name:     "remove scheme",
			inputURL: "https://example.com",
			expected: "example.com",
		},
		{
			name:     "remove uppercase and www",
			inputURL: "HTTP://WWW.EXAMPLE.COM/PATH",
			expected: "example.com/path",
		},
		{
			name:     "trailing slash",
			inputURL: "http://www.example.com/path/",
			expected: "example.com/path",
		},
		{
			name:     "remove default port 80",
			inputURL: "http://www.example.com:80/path",
			expected: "example.com/path",
		},
		{
			name:     "remove default port 443",
			inputURL: "https://www.example.com:443/thing",
			expected: "example.com/thing",
		},
		{
			name:     "remove query and fragment",
			inputURL: "http://www.example.com/path?x=10#section1",
			expected: "example.com/path",
		},
		{
			name:     "redundant path segments",
			inputURL: "http://www.example.com/a/b/../c/./d.html",
			expected: "example.com/a/c/d.html",
		},
		{
			name:     "handle double slashes",
			inputURL: "http://example.com/a//b///c",
			expected: "example.com/a/b/c",
		},
	}

	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := normalizeURL(tc.inputURL)
			if err != nil {
				t.Errorf("Test %v - '%s' FAIL: unexpected error: %v", i, tc.name, err)
			}
			if actual != tc.expected {
				t.Errorf("Test %v - '%s' FAIL: expected URL: %v, actual: %v", i, tc.name, tc.expected, actual)
			}
		})
	}
}
