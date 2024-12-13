package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/btcsuite/btcutil/base58"
)

func decodeReadableString(encoded string) (string, error) {
	// Remove hyphens from the input
	cleaned := strings.ReplaceAll(encoded, "-", "")
	
	// Decode from base58
	decoded := base58.Decode(cleaned)
	if len(decoded) == 0 {
		return "", fmt.Errorf("invalid encoded string")
	}

	// Convert back to string and validate as URL
	urlStr := string(decoded)
	_, err := url.Parse(urlStr)
	if err != nil {
		return "", fmt.Errorf("decoded string is not a valid URL: %v", err)
	}

	return urlStr, nil
}

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
	urlFlag := flag.String("url", "", "URL to convert to readable string")
	decodeFlag := flag.String("decode", "", "Decode a shortened string back to URL")
	flag.Parse()

	if *urlFlag != "" && *decodeFlag != "" {
		log.Fatal("Please use either -url OR -decode, not both")
	}

	if *urlFlag == "" && *decodeFlag == "" {
		log.Fatal("Please provide either -url or -decode flag")
	}

	if *urlFlag != "" {
		// Encode mode
		result, err := createReadableString(*urlFlag)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(result)
	} else {
		// Decode mode
		result, err := decodeReadableString(*decodeFlag)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(result)
	}
}
