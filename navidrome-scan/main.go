package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// apiResponse defines the standard JSON structure for API responses.
type apiResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func scanHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondWithError(w, http.StatusMethodNotAllowed, "Only POST method is allowed")
		return
	}

	// Create a dedicated HTTP client with a reasonable timeout
	httpClient := &http.Client{Timeout: 30 * time.Second}

	// Build the Navidrome webhook URL form env variables
	scanURL := fmt.Sprintf("%s/rest/startScan.view?u=%s&p=%s&v=1.16.1&c=go-api&f=json",
		os.Getenv("NAVIDROME_API_URL"),
		os.Getenv("NAVIDROME_USER"),
		os.Getenv("NAVIDROME_PASS"),
	)

	log.Println("Received request, executing Navidrome scan 🔎")
	navidromeReq, err := http.NewRequest("GET", scanURL, nil)
	if err != nil {
		log.Printf("❌ Failed to create Navidrome request: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to create internal request")
		return
	}

	// Execute the request to the Navidrome API
	resp, err := httpClient.Do(navidromeReq)
	if err != nil {
		log.Printf("❌ Navidrome scan request failed: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to execute Navidrome scan")
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	log.Printf("ℹ️  Navidrome scan triggered. Status; %s, Response: %s", resp.Status, string(body))

	respondWithJSON(w, http.StatusOK, apiResponse{Status: "ok", Message: "Navidrome scan executed successfully."})
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/scan", scanHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	listenAddr := ":" + port
	server := &http.Server{Addr: listenAddr, Handler: mux}

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Printf("🗣️  Navidrome scan service listening on %s", listenAddr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("❌ Server failed to start: %v", err)
		}
	}()

	<-stopChan

	log.Println("ℹ️  Shutting down service")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("❌ Server shutdown failed: %v", err)
	}
	log.Println("✅ Server gracefully stopped.")
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, apiResponse{Status: "error", Message: message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
