package main

import (
	"net/url"
	"reflect"
	"testing"
)

func TestGetH1FromHTML(t *testing.T) {
	tests := []struct {
		name      string
		inputHTML string
		expected  string
	}{
		{
			name:      "get <h1> content",
			inputHTML: "<h1>Great Books on Leadership</h1>",
			expected:  "Great Books on Leadership",
		},
		{
			name:      "no <h1>",
			inputHTML: "It can be really hard to put your code out there and ask someone to give you advice on how to make it better. It’s easy to take that advice personally.",
			expected:  "",
		},
		{
			name: "many tags",
			inputHTML: `
<div class="ly-list-edit ly-dim listly-form-container listly-icon-select listly-hide-empty" title="Edit List Settings"></div>
<h1 class="ly-title" data-show-hide="show_list_title">Great Books on Communication</h1>
<p class="ly-author-link" data-show-hide="show_author">Listly by <a href="//list.ly/cuchullainn" class="ly-ext-link ly-dim">Conor Neill</a></p>
	`,
			expected: "Great Books on Communication",
		},
		// 		{
		// 			name: "multiple h1 tags",
		// 			inputHTML: `
		// 			<h1>Great Books on Leadership</h1>
		// 			<h1>Great Books on Communication</h1>
		// 			<h1>Great Books on Life</h1>
		// 			`,
		// 			expected: `
		// Great Books on Leadership
		// Great Books on Communication
		// Great Books on Life`,
		// 		},
		{
			name:      "idk",
			inputHTML: "<html><body><h1>Test Title</h1></body></html>",
			expected:  "Test Title",
		},
	}

	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := getH1FromHTML(tc.inputHTML)

			if actual != tc.expected {
				t.Errorf("Test %v - '%s' FAIL: \nexpected: %v, \nactual: %v", i, tc.name, tc.expected, actual)
			}
		})
	}
}

func TestGetFirstParagraphFromHTMLMainPriority(t *testing.T) {
	tests := []struct {
		name      string
		inputHTML string
		expected  string
	}{
		{
			name:      "get no content",
			inputHTML: "<h1>Great Books on Leadership</h1>",
			expected:  "",
		},
		{
			name:      "only p",
			inputHTML: "<p>It can be really hard to put your code out there and ask someone to give you advice on how to make it better. It’s easy to take that advice personally.<p>",
			expected:  "It can be really hard to put your code out there and ask someone to give you advice on how to make it better. It’s easy to take that advice personally.",
		},
		{
			name: "get main p",
			inputHTML: `<html><body>
		<p>Outside paragraph.</p>
		<main>
			<p>Main paragraph.</p>
		</main>
	</body></html>`,
			expected: "Main paragraph.",
		},
		{
			name: "many tags",
			inputHTML: `
<div class="ly-list-edit ly-dim listly-form-container listly-icon-select listly-hide-empty" title="Edit List Settings"></div>
<h1 class="ly-title" data-show-hide="show_list_title">Great Books on Communication</h1>
<p class="ly-author-link" data-show-hide="show_author">Listly by <a href="//list.ly/cuchullainn" class="ly-ext-link ly-dim">Conor Neill</a></p>
	`,
			expected: "Listly by Conor Neill",
		},
	}

	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := getFirstParagraphFromHTML(tc.inputHTML)

			if actual != tc.expected {
				t.Errorf("Test %v - '%s' FAIL: \nexpected: %v, \nactual: %v", i, tc.name, tc.expected, actual)
			}
		})
	}
}

func TestGetURLsFromHTML(t *testing.T) {
	tests := []struct {
		name      string
		inputURL  string
		inputHTML string
		expected  []string
	}{
		{
			name:      "single url, not relative",
			inputURL:  "https://blog.boot.dev",
			inputHTML: `<html><body><a href="https://blog.boot.dev"><span>Boot.dev</span></a></body></html>`,
			expected:  []string{"https://blog.boot.dev"},
		},
		{
			name:     "multiple urls, not relative",
			inputURL: "https://gobyexample.com/",
			inputHTML: `<div>
					<h2><a href="https://gobyexample.com/">Go by Example</a>: Variables</h2>
					<p class="next">
						Next example: <a href="https://gobyexample.com/constants" rel="next">Constants</a>.
					</p>
					<p class="footer">
						by <a href="https://markmcgranaghan.com">Mark McGranaghan</a> and <a href="https://eli.thegreenplace.net">Eli Bendersky</a> | <a href="https://github.com/mmcgrana/gobyexample">source</a> | <a href="https://github.com/mmcgrana/gobyexample#license">license</a>
					</p>
				</div>
			`,
			expected: []string{"https://gobyexample.com", "https://gobyexample.com/constants", "https://markmcgranaghan.com", "https://eli.thegreenplace.net", "https://github.com/mmcgrana/gobyexample", "https://github.com/mmcgrana/gobyexample#license"},
		},
		{
			name:      "single url, relative",
			inputURL:  "https://blog.boot.dev",
			inputHTML: `<html><body><a href="./"><span>Boot.dev</span></a></body></html>`,
			expected:  []string{"https://blog.boot.dev"},
		},
		{
			name:     "multiple urls, relative",
			inputURL: "https://gobyexample.com/",
			inputHTML: `<div>
					<h2><a href="./">Go by Example</a>: Variables</h2>
					<p class="next">
						Next example: <a href="/constants" rel="next">Constants</a>.
					</p>
				</div>
			`,
			expected: []string{"https://gobyexample.com", "https://gobyexample.com/constants"},
		},
	}

	for i, tc := range tests {
		baseURL, err := url.Parse(tc.inputURL)
		if err != nil {
			t.Errorf("couldn't parse input URL: %v", err)
		}

		actual, err := getURLsFromHTML(tc.inputHTML, baseURL)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if !reflect.DeepEqual(actual, tc.expected) {
			t.Errorf("Test %v - '%s' FAIL: \nexpected: %v, \nactual: %v", i, tc.name, tc.expected, actual)
		}
	}
}
