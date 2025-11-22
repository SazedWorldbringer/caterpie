package main

import (
	"encoding/csv"
	"os"
	"strings"
)

func writeCSVReport(pages map[string]PageData, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// write header
	writer.Write([]string{"page_url", "h1", "first_paragraph", "outgoing_link_urls", "image_urls"})

	for _, page := range pages {
		outLinks := strings.Join(page.OutgoingLinks, ";")
		imgURLs := strings.Join(page.ImageURLs, ";")

		// write record
		writer.Write([]string{page.URL, page.H1, page.FirstParagraph, outLinks, imgURLs})
	}

	return nil
}
