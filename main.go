package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/btcsuite/btcutil/base58"
)

func createReadableString(inputURL string) (string, error) {
	// Parse and validate the URL
	parsedURL, err := url.Parse(inputURL)
	if err != nil {
		return "", fmt.Errorf("invalid URL: %v", err)
	}

	// Ensure scheme is present
	if parsedURL.Scheme == "" {
		parsedURL.Scheme = "https"
	}

	// Convert URL to bytes
	urlBytes := []byte(parsedURL.String())
	
	// Encode to base58
	encoded := base58.Encode(urlBytes)
	
	// Take first 12 characters for a shorter, but still unique string
	if len(encoded) > 12 {
		encoded = encoded[:12]
	}
	
	// Insert a hyphen every 4 characters for better readability
	var result strings.Builder
	for i, char := range encoded {
		if i > 0 && i%4 == 0 {
			result.WriteRune('-')
		}
		result.WriteRune(char)
	}
	
	return result.String(), nil
}

func main() {
	// Parse command line flags
	url := flag.String("url", "", "URL to convert to readable string")
	flag.Parse()

	if *url == "" {
		log.Fatal("Please provide a URL using the -url flag")
	}

	// Convert URL to readable string
	result, err := createReadableString(*url)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(result)
}
