package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

type config struct {
	pages              map[string]PageData
	baseURL            *url.URL
	mu                 *sync.Mutex
	concurrencyControl chan struct{}
	wg                 *sync.WaitGroup
}

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
	rawBase := os.Args[1]
	fmt.Printf("Starting crawl of: %s\n", rawBase)

	parsedBase, err := url.Parse(rawBase)
	if err != nil {
		fmt.Println(err)
	}

	cfg := &config{
		pages:              make(map[string]PageData),
		baseURL:            parsedBase,
		mu:                 &sync.Mutex{},
		concurrencyControl: make(chan struct{}, 10),
		wg:                 &sync.WaitGroup{},
	}

	cfg.wg.Add(1)
	go cfg.crawlPage(rawBase)
	cfg.wg.Wait()

	cfg.mu.Lock()
	finalPages := cfg.pages
	cfg.mu.Unlock()

	for _, page := range finalPages {
		fmt.Println("URL:", page.URL)
		fmt.Println("H1:", page.H1)
		fmt.Println("First Paragraph:", page.FirstParagraph)
		fmt.Println("Outgoing links:")
		for i, url := range page.OutgoingLinks {
			fmt.Println(i, "-", url)
		}
		fmt.Println("ImageURLs:")
		for i, url := range page.ImageURLs {
			fmt.Println(i, "-", url)
		}
		fmt.Println("---")
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

// func crawlPage(rawBaseURL, rawCurrentURL string, pages map[string]int) {
func (cfg *config) crawlPage(rawCurrentURL string) {
	// signal done
	defer cfg.wg.Done()

	// check if rawCurrentURL is on the same domain as rawBaseURL
	// return if not
	currentURL, err := url.Parse(rawCurrentURL)
	if err != nil {
		fmt.Println(err)
	}

	if strings.Compare(cfg.baseURL.Hostname(), currentURL.Hostname()) != 0 {
		return
	}

	// check if normalizeURL(rawCurrentURL) is in pages
	normalized, err := normalizeURL(currentURL.String())
	if err != nil {
		fmt.Println(err)
	}

	cfg.mu.Lock()
	if _, exists := cfg.pages[normalized]; exists {
		cfg.mu.Unlock()
		return
	}
	cfg.pages[normalized] = PageData{URL: currentURL.String()}
	cfg.mu.Unlock()

	// fmt.Println("current: ", normalized)

	// get html for currentURL and store html in pages[currentURLstr]
	// limit concurrent access
	cfg.concurrencyControl <- struct{}{}
	html, err := getHTML(currentURL.String())
	<-cfg.concurrencyControl
	// if error in getting html, store final page data and return
	if err != nil {
		cfg.mu.Lock()
		cfg.pages[normalized] = PageData{
			URL:            currentURL.String(),
			H1:             "",
			FirstParagraph: "",
		}
		cfg.mu.Unlock()
		fmt.Println(err)
		return
	}

	// parse page
	pageData := extractPageData(html, currentURL)

	cfg.mu.Lock()
	cfg.pages[normalized] = pageData
	cfg.mu.Unlock()

	// get all outgoing links from current html
	for _, link := range pageData.OutgoingLinks {
		// recursively crawl each outgoing link
		cfg.wg.Add(1)
		go cfg.crawlPage(link)
	}
}
