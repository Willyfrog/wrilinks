package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/url"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"urlshortener/readable"
)

var db *sql.DB

func initDB() error {
	var err error
	db, err = sql.Open("sqlite3", "urls.db")
	if err != nil {
		return err
	}

	// Create table if it doesn't exist
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS urls (
			short_code TEXT PRIMARY KEY,
			original_url TEXT NOT NULL
		)
	`)
	return err
}

func decodeReadableString(encoded string) (string, error) {
	// Remove hyphens from the input
	cleaned := strings.ReplaceAll(encoded, "-", "")
	
	// Query the database for the original URL
	var originalURL string
	err := db.QueryRow("SELECT original_url FROM urls WHERE short_code = ?", cleaned).Scan(&originalURL)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("no URL found for code: %s", encoded)
		}
		return "", err
	}

	return originalURL, nil
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

	fullURL := parsedURL.String()
	
	// Generate a readable code
	urlBytes := []byte(fullURL + fmt.Sprintf("%d", strings.Count(fullURL, "")))
	encoded := readable.GenerateReadableString(urlBytes)
	
	// Store in database
	_, err = db.Exec("INSERT INTO urls (short_code, original_url) VALUES (?, ?)", encoded, fullURL)
	if err != nil {
		return "", fmt.Errorf("failed to store URL: %v", err)
	}
	
	// Insert hyphens for readability
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
	if err := initDB(); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()
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
