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
