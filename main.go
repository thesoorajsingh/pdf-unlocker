package main

import (
	"log"
	"net/http"
)

func main() {
	mux := setupRoutes()
	log.Println("Server starting on port 8080...")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}