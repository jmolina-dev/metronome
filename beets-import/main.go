package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type AppState struct {
	mu sync.Mutex
}

type apiResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func (s *AppState) taskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondWithError(w, http.StatusMethodNotAllowed, "Only POST method is allowed")
		return
	}

	if !s.mu.TryLock() {
		respondWithError(w, http.StatusConflict, "A task is already in progress")
		return
	}

	go s.runTask()

	respondWithJSON(w, http.StatusAccepted, apiResponse{Status: "ok", Message: "Import and scan task started"})
}

func (s *AppState) runTask() {
	defer s.mu.Unlock()

	log.Println("ℹ️  Starting beet import  🫜")
	cmd := exec.Command("beet", "import", "/app/downloads")
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("❌  Beets import failed: %v\nOutput: %s", err, string(output))
		return
	}
	log.Printf("✅  Beets import finished successfully  🫜")

	navidromeScanURL := os.Getenv("NAVIDROME_SCAN_URL")
	if navidromeScanURL == "" {
		log.Println("⚠️  NAVIDROME_SCAN_URL not set. Skipping task.")
		return
	}

	log.Println("ℹ️  Executing Navidrome scan")
	httpClient := &http.Client{Timeout: 30 * time.Second}
	req, err := http.NewRequest("POST", navidromeScanURL, nil)
	if err != nil {
		log.Printf("❌  Failed to create Navidrome request: %v", err)
		return
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		log.Printf("❌  Navidrome scan request failed: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		log.Printf("❌  Navidrome API returned non-success status: %s", resp.Status)
		return
	}
	log.Println("✅  Navidrome scan launched successfully")
}

func main() {
	state := &AppState{}
	mux := http.NewServeMux()
	mux.HandleFunc("/task/start", state.taskHandler)

	port := os.Getenv("BEETS_PORT")
	if port == "" {
		port = "8081"
	}
	listenAddr := ":" + port
	server := &http.Server{Addr: listenAddr, Handler: mux}

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Printf("🗣️  Beets-webhook listening on %s", listenAddr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("❌  Server failed to start: %v", err)
		}
	}()

	<-stopChan

	log.Println("ℹ️  Shutting down server")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("❌  Server shutdown failed: %v", err)
	}
	log.Println("🛑  Server stopped")
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
