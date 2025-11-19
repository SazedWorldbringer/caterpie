package main

import "testing"

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
