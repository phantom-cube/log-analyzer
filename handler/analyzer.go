package handler

import (
	"io"
	"log"
	"net/http"

	"log-analyzer/client"

	"github.com/gin-gonic/gin"
)

const maxLogSize = 10 * 1024 * 1024 // 10MB

type AnalyzeResponse struct {
	Analysis string `json:"analysis"`
	Error    string `json:"error,omitempty"`
}

func AnalyzeLog(c *gin.Context) {
	// Get uploaded file
	file, err := c.FormFile("logfile")
	if err != nil {
		c.JSON(http.StatusBadRequest, AnalyzeResponse{
			Error: "Failed to get log file: " + err.Error(),
		})
		return
	}

	// Validate file size
	if file.Size > maxLogSize {
		c.JSON(http.StatusBadRequest, AnalyzeResponse{
			Error: "File too large (max 10MB)",
		})
		return
	}

	// Open uploaded file
	uploadedFile, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, AnalyzeResponse{
			Error: "Failed to open log file: " + err.Error(),
		})
		return
	}
	defer uploadedFile.Close()

	// Read log content
	logContent, err := io.ReadAll(uploadedFile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, AnalyzeResponse{
			Error: "Failed to read log file: " + err.Error(),
		})
		return
	}

	if len(logContent) == 0 {
		c.JSON(http.StatusBadRequest, AnalyzeResponse{
			Error: "Log file is empty",
		})
		return
	}

	log.Printf("Analyzing log file: %s (size: %d bytes)", file.Filename, len(logContent))

	// Call Ollama API for analysis
	analysis, err := client.AnalyzeLogWithAI(string(logContent))
	if err != nil {
		log.Printf("Error analyzing log: %v", err)
		c.JSON(http.StatusInternalServerError, AnalyzeResponse{
			Error: "Failed to analyze log: " + err.Error(),
		})
		return
	}

	// Send response
	c.JSON(http.StatusOK, AnalyzeResponse{
		Analysis: analysis,
	})
}
