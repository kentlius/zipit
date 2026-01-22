package main

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"net/http"
	"path"
)

func zipHandler(w http.ResponseWriter, r *http.Request) {
	encryptedPayload := r.URL.Path[1:]
	if encryptedPayload == "" {
		http.Error(w, "Missing encrypted path", http.StatusBadRequest)
		return
	}

	targetURL, err := decryptURL(encryptedPayload, encryptionKey)
	if err != nil {
		log.Printf("Decryption Error: %v", err)
		http.Error(w, "Invalid Link", http.StatusForbidden)
		return
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
	if filename == "" || filename == "." {
		filename = "file.bin"
	}

	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", "attachment; filename=\"download.zip\"")

	zipWriter := zip.NewWriter(w)
	defer zipWriter.Close()

	zipEntry, err := zipWriter.Create(filename)
	if err != nil {
		log.Printf("Zip creation failed: %v", err)
		return
	}

	if _, err := io.Copy(zipEntry, resp.Body); err != nil {
		log.Printf("Stream interrupted: %v", err)
		return
	}
}
