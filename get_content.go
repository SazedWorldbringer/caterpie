package main

import (
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func getH1FromHTML(html string) string {
	reader := strings.NewReader(html)

	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		log.Fatal(err)
	}

	output := ""

	doc.Find("h1").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the title
		content := s.Contents().Text()
		output += content
	})

	return output
}

func getFirstParagraphFromHTML(html string) string {
	reader := strings.NewReader(html)

	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		log.Fatal(err)
	}

	// if main exists, get the first p in it
	if doc.Find("main").Length() > 0 {
		return doc.Find("main").Find("p").First().Contents().Text()
	}

	// fallback to find the first p, if main doesn't exist
	return doc.Find("p").First().Contents().Text()
}
