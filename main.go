package main

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/gorilla/mux"
)

// OpenAI API endpoint
const openAIAPIURL = "https://api.openai.com/v1/chat/completions"

// APIKey is your OpenAI API key that should be kept secret
var apiKey = os.Getenv("OPENAI_API_KEY")

// ProxyHandler handles requests, adds auth header, and forwards the request to OpenAI
func ProxyHandler(w http.ResponseWriter, r *http.Request) {
	// received request
	log.Printf("Received request to %s", r.URL.Path)
	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read request body", http.StatusBadRequest)
		return
	}

	// Forward the request to OpenAI's API
	proxyRequest(w, body)
}

// proxyRequest sends the modified request to OpenAI's API with the correct headers
func proxyRequest(w http.ResponseWriter, body []byte) {
	// Create the OpenAI API URL
	openAIURL, _ := url.Parse(openAIAPIURL)

	// Create a new HTTP request with the original body and the OpenAI URL
	req, err := http.NewRequest("POST", openAIURL.String(), io.NopCloser(bytes.NewReader(body)))
	if err != nil {
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		return
	}

	log.Printf("Forwarding request to %s", openAIURL.String())
	// Add headers
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	// Forward the request to OpenAI
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Error forwarding request", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Read the response from OpenAI and return it to the client
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Error reading OpenAI response", http.StatusInternalServerError)
		return
	}

	// Set the response headers
	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Set(key, value)
		}
	}

	w.WriteHeader(resp.StatusCode)
	w.Write(responseBody)
}

func main() {
	// Load the OpenAI API key from the environment
	if apiKey == "" {
		log.Fatal("API key not set in environment variable OPENAI_API_KEY")
	}

	// Create a new router
	r := mux.NewRouter()

	// Define the proxy endpoint
	r.HandleFunc("/chat/completions", ProxyHandler).Methods("POST")

	// Start the server
	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
