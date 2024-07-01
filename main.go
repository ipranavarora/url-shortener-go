package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type URL struct {
	ID           string    `json:"id"`
	OriginalURL  string    `json:"original_url"`
	ShortURL     string    `json:"short_url"`
	CreationDate time.Time `json:"creation_date"`
}

// in-memory database
var urlDB = make(map[string]URL)

func generateShortUrl(OriginalURL string) string {
	hasher := md5.New()
	hasher.Write([]byte(OriginalURL)) // it converts the originalURL into a byte slice
	data := hasher.Sum(nil)
	hash := hex.EncodeToString(data)
	return hash[:8]
}

func createURL(originalURL string) string {
	shortUrl := generateShortUrl(originalURL)
	id := shortUrl // use the shorturl as id for simplicity
	urlDB[id] = URL{
		ID:           id,
		OriginalURL:  originalURL,
		ShortURL:     shortUrl,
		CreationDate: time.Now(),
	}
	return shortUrl
}

func getURL(id string) (URL, error) {
	url, ok := urlDB[id]
	if !ok {
		return URL{}, errors.New("URL not found")
	}
	return url, nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Home")
}

func shortUrlHandler(w http.ResponseWriter, r *http.Request) {
	var data struct {
		URL string `json:"url"`
	}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	shortURL_ := createURL(data.URL)

	response := struct {
		ShortURL string `json:"short_url"`
	}{
		ShortURL: shortURL_,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func redirectUrlHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/redirect/")
	url, err := getURL(id)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusNotFound)
		return
	}
	http.Redirect(w, r, url.OriginalURL, http.StatusFound)
}

func main() {
	// Register the handler function to handle all requests to the root url ("/")
	http.HandleFunc("/", handler)
	http.HandleFunc("/shorten", shortUrlHandler)
	http.HandleFunc("/redirect/", redirectUrlHandler) // Added the redirect handler

	// Start the http server on port 8080
	fmt.Println("Starting the server on port 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error on starting server", err)
	}
}
