package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
)

func setupRoutes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/unlock", unlockHandler)
	mux.HandleFunc("/health", healthHandler)
	return mux
}

func unlockHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to read file from form", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Read the file content into a byte slice
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Failed to read file content", http.StatusInternalServerError)
		return
	}

	// Check if the file is a PDF
	if http.DetectContentType(fileBytes) != "application/pdf" {
		http.Error(w, "Uploaded file is not a PDF", http.StatusBadRequest)
		return
	}

	password := r.FormValue("password")
	if password == "" {
		http.Error(w, "Password not provided", http.StatusBadRequest)
		return
	}

	unlockedPdf, err := unlockPdf(fileBytes, password)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to unlock PDF: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment; filename=\"unlocked.pdf\"")
	if _, err := w.Write(unlockedPdf); err != nil {
		// Log the error, but the response is likely already sent.
		log.Printf("Failed to write unlocked PDF to response: %v", err)
	}
}

func unlockPdf(fileBytes []byte, password string) ([]byte, error) {
	reader := bytes.NewReader(fileBytes)
	writer := &bytes.Buffer{}
	config := model.NewDefaultConfiguration()
	config.UserPW = password

	err := api.Decrypt(reader, writer, config)
	if err != nil {
		return nil, err
	}

	return writer.Bytes(), nil
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("I'm healthy and sound!"))
}
