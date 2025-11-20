package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Error: no website provided")
		os.Exit(1)
	}
	if len(os.Args) > 2 {
		fmt.Println("Error: too many arguments provided")
		os.Exit(1)
	}

	// exactly 1 argument, base url
	BASE_URL := os.Args[1]
	fmt.Printf("Starting crawl of: %s\n", BASE_URL)

	pages := map[string]int{}

	crawlPage(BASE_URL, BASE_URL, pages)

	for page, value := range pages {
		fmt.Printf("%s - %d\n", page, value)
	}
}

func getHTML(rawURL string) (string, error) {
	// http client
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	// create request
	URL, err := url.Parse(rawURL)
	if err != nil {
		fmt.Println(err)
	}
	URLStr := URL.String()
	method := "GET"

	req, err := http.NewRequest(method, URLStr, nil)
	if err != nil {
		fmt.Println("Error creating request: ", err)
		return "", err
	}

	req.Header.Add("Content-Type", "text/html")
	req.Header.Add("User-Agent", "Caterpie/1.0")

	// make request
	resp, err := client.Do(req)
	// error in making the request - timeout, connection refused, etc
	if err != nil {
		fmt.Println("Error making request: ", err)
		return "", err
	}
	defer resp.Body.Close()

	// exit if error code
	if resp.StatusCode > 400 {
		return "", fmt.Errorf("%d", resp.StatusCode)
	}

	contentHeader := resp.Header.Get("Content-Type")

	// exit if content type mismatch
	if !strings.Contains(contentHeader, "text/html") {
		return "", fmt.Errorf("content type mismatch, %v", resp.Header.Get("Content-Type"))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body: ", err)
		return "", err
	}

	return string(body), nil
}

func crawlPage(rawBaseURL, rawCurrentURL string, pages map[string]int) {
	// check if rawCurrentURL is on the same domain as rawBaseURL
	// return if not
	baseURL, err := url.Parse(rawBaseURL)
	if err != nil {
		fmt.Println(err)
	}
	currentURL, err := url.Parse(rawCurrentURL)
	if err != nil {
		fmt.Println(err)
	}

	if strings.Compare(baseURL.Hostname(), currentURL.Hostname()) != 0 {
		return
	}

	// check if normalizeURL(rawCurrentURL) is in pages
	// if yes, increment pages[normalizeURL(rawCurrentURL)]
	// if no, make pages[normalizeURL(rawCurrentURL)] = 1
	currentURLstr, err := normalizeURL(currentURL.String())
	if err != nil {
		fmt.Println(err)
	}

	_, ok := pages[currentURLstr]
	if ok {
		pages[currentURLstr]++
		return
	} else {
		pages[currentURLstr] = 1
	}

	// get html for rawCurrentURL and print
	html, err := getHTML(currentURL.String())
	if err != nil {
		fmt.Println(err)
	}

	// get all outgoing links from current html
	outgoingLinks := extractPageData(html, baseURL.String()).OutgoingLinks
	for _, link := range outgoingLinks {
		// recursively crawl each outgoing link
		crawlPage(rawBaseURL, link, pages)
	}
}
