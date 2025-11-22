# Caterpie!
## Crawl the web with Caterpie!
![Caterpie Banner](./010-caterpie-pokemon-pixel-art.png)

## Features

- Only crawls URLs within the domain, avoids going "off-site"
- Concurrently fetches multiple pages at once, worker-style limit via a buffered channel
- Extracts structured page data 
  - H1
  - First paragraph
  - Outgoing links
  - Image URLs
- Has a max crawl page limit, to prevent expanding infinitely over very large sites, save bandwidth, time, and avoids hitting servers too hard

## Installation

### Requirements

- Go 1.22+
- Linux/macOS/Windows

```sh
go install github.com/SazedWorldbringer/caterpie@latest
```

### Run

```sh
  Usage: caterpie URL maxConcurrency maxPages

  caterpie https://example.com 10 10
```

This will generate a report.csv in pwd

## Development setup

### Clone the repo
```sh
git clone https://github.com/SazedWorldbringer/caterpie
cd caterpie
```

### Run tests
```sh
go test ./...
```

### Run
```sh
go run . https://example.com 10 10
```

### Build
```sh
go build .
```
