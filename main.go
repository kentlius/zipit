package main

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"time"
)

var httpClient = &http.Client{
	Timeout: 30 * time.Second,
}

func main() {
	http.HandleFunc("/", zipHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func zipHandler(w http.ResponseWriter, r *http.Request) {
	targetURL := r.URL.Path[1:]

	if targetURL == "" {
		http.Error(w, "Please provide a URL in the path, e.g., /example.com/file.pdf", http.StatusBadRequest)
		return
	}

	if !strings.HasPrefix(targetURL, "http") {
		targetURL = "https://" + targetURL
	}

	if strings.HasPrefix(targetURL, "https:/") && !strings.HasPrefix(targetURL, "https://") {
		targetURL = strings.Replace(targetURL, "https:/", "https://", 1)
	}

	if r.URL.RawQuery != "" {
		targetURL += "?" + r.URL.RawQuery
	}

	log.Printf("Downloading: %s", targetURL)

	resp, err := httpClient.Get(targetURL)
	if err != nil {
		http.Error(w, "Failed to download file: "+err.Error(), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, fmt.Sprintf("Upstream returned status: %d", resp.StatusCode), http.StatusBadGateway)
		return
	}

	filename := path.Base(targetURL)
	if filename == "" || filename == "." || filename == "/" {
		filename = "downloaded-file"
	}

	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s.zip\"", filename))

	zipWriter := zip.NewWriter(w)
	defer zipWriter.Close()

	zipEntry, err := zipWriter.Create(filename)
	if err != nil {
		log.Printf("Failed to create zip entry: %v", err)
		return
	}

	if _, err := io.Copy(zipEntry, resp.Body); err != nil {
		log.Printf("Stream interrupted: %v", err)
		return
	}
}
