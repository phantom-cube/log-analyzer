package handler

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"log-analyzer/client"
)

const maxLogSize = 10 * 1024 * 1024 // 10MB

type AnalyzeResponse struct {
	Analysis string `json:"analysis"`
	Error    string `json:"error,omitempty"`
}

func AnalyzeLog(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse multipart form
	if err := r.ParseMultipartForm(maxLogSize); err != nil {
		sendError(w, "Failed to parse form: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Get uploaded file
	file, header, err := r.FormFile("logfile")
	if err != nil {
		sendError(w, "Failed to get log file: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Validate file size
	if header.Size > maxLogSize {
		sendError(w, "File too large (max 10MB)", http.StatusBadRequest)
		return
	}

	// Read log content
	logContent, err := io.ReadAll(file)
	if err != nil {
		sendError(w, "Failed to read log file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if len(logContent) == 0 {
		sendError(w, "Log file is empty", http.StatusBadRequest)
		return
	}

	log.Printf("Analyzing log file: %s (size: %d bytes)", header.Filename, len(logContent))

	// Call OpenAI API for analysis
	analysis, err := client.AnalyzeLogWithAI(string(logContent))
	if err != nil {
		log.Printf("Error analyzing log: %v", err)
		sendError(w, "Failed to analyze log: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Send response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(AnalyzeResponse{
		Analysis: analysis,
	})
}

func sendError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(AnalyzeResponse{
		Error: message,
	})
}
