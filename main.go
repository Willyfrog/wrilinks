package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"urlshortener/readable"
)

type PageData struct {
	ShortURL    string
	OriginalURL string
}

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

func handleIndex(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	tmpl.Execute(w, &PageData{})
}

func handleShorten(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	inputURL := r.FormValue("url")
	shortURL, err := createReadableString(inputURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	tmpl.Execute(w, &PageData{ShortURL: shortURL})
}

func handleDecode(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	code := r.FormValue("code")
	originalURL, err := decodeReadableString(code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	tmpl.Execute(w, &PageData{OriginalURL: originalURL})
}

func main() {
	if err := initDB(); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	// Set up HTTP routes
	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/shorten", handleShorten)
	http.HandleFunc("/decode", handleDecode)

	// Start the server
	fmt.Println("Server starting on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
