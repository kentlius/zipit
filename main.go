package main

import (
	"log"
	"net/http"
	"os"
	"time"
)

var (
	httpClient = &http.Client{
		Timeout: 30 * time.Second,
	}
	encryptionKey []byte
)

func main() {
	keyStr := os.Getenv("ENCRYPTION_KEY")
	if keyStr == "" {
		log.Fatal("ENCRYPTION_KEY environment variable is required")
	}
	if len(keyStr) != 32 {
		log.Fatal("ENCRYPTION_KEY must be exactly 32 characters long")
	}
	encryptionKey = []byte(keyStr)

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
