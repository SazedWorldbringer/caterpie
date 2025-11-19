package main

import (
	"fmt"
	"os"
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
}
