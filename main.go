package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"log-analyzer/handler"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using default configuration")
	}

	// Setup routes
	http.HandleFunc("/", handleHome)
	http.HandleFunc("/analyze", handler.AnalyzeLog)
	http.HandleFunc("/health", handleHealth)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server starting on port %s...\n", port)
	fmt.Println("Available endpoints:")
	fmt.Println("  - GET  /        - Home page")
	fmt.Println("  - POST /analyze - Analyze log file")
	fmt.Println("  - GET  /health  - Health check")

	// Create server with extended timeouts for LLM processing
	server := &http.Server{
		Addr:         ":" + port,
		ReadTimeout:  15 * time.Minute,
		WriteTimeout: 15 * time.Minute,
		IdleTimeout:  2 * time.Minute,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	html := `<!DOCTYPE html>
<html>
<head>
    <title>AI Log Analyzer</title>
    <style>
        body { font-family: Arial, sans-serif; max-width: 800px; margin: 50px auto; padding: 20px; }
        h1 { color: #333; }
        .upload-form { background: #f5f5f5; padding: 20px; border-radius: 8px; margin: 20px 0; }
        button { background: #4CAF50; color: white; padding: 10px 20px; border: none; cursor: pointer; border-radius: 4px; }
        button:hover { background: #45a049; }
        #result { margin-top: 20px; padding: 15px; background: #fff; border: 1px solid #ddd; border-radius: 4px; }
        .loading { display: none; color: #666; }

        /* Markdown rendering styles */
        #resultContent { line-height: 1.6; }
        #resultContent h1 { font-size: 1.8em; margin-top: 20px; margin-bottom: 10px; border-bottom: 2px solid #333; }
        #resultContent h2 { font-size: 1.5em; margin-top: 18px; margin-bottom: 8px; color: #444; }
        #resultContent h3 { font-size: 1.3em; margin-top: 15px; margin-bottom: 6px; color: #555; }
        #resultContent strong { color: #d9534f; font-weight: bold; }
        #resultContent em { font-style: italic; color: #555; }
        #resultContent code { background: #f4f4f4; padding: 2px 6px; border-radius: 3px; font-family: monospace; }
        #resultContent pre { background: #f4f4f4; padding: 10px; border-radius: 4px; overflow-x: auto; }
        #resultContent ul, #resultContent ol { margin: 10px 0; padding-left: 25px; }
        #resultContent li { margin: 5px 0; }
        #resultContent p { margin: 10px 0; }
        #resultContent blockquote { border-left: 4px solid #ddd; padding-left: 15px; margin: 10px 0; color: #666; }
    </style>
</head>
<body>
    <h1>AI Log Analyzer</h1>
    <p>Upload your log file for AI-powered analysis</p>

    <div class="upload-form">
        <h3>Upload Log File</h3>
        <form id="uploadForm">
            <input type="file" id="logFile" accept=".log,.txt" required>
            <br><br>
            <button type="submit">Analyze Log</button>
            <span class="loading" id="loading">Analyzing...</span>
        </form>
    </div>

    <div id="result" style="display:none;">
        <h3>Analysis Result:</h3>
        <div id="resultContent"></div>
    </div>

    <script>
        // Simple Markdown to HTML converter
        function markdownToHtml(markdown) {
            let html = markdown;

            // Escape HTML tags first
            html = html.replace(/&/g, '&amp;')
                       .replace(/</g, '&lt;')
                       .replace(/>/g, '&gt;');

            // Headers (must be at start of line)
            html = html.replace(/^### (.*$)/gm, '<h3>$1</h3>');
            html = html.replace(/^## (.*$)/gm, '<h2>$1</h2>');
            html = html.replace(/^# (.*$)/gm, '<h1>$1</h1>');

            // Bold and Italic
            html = html.replace(/\*\*\*(.+?)\*\*\*/g, '<strong><em>$1</em></strong>');
            html = html.replace(/\*\*(.+?)\*\*/g, '<strong>$1</strong>');
            html = html.replace(/\*(.+?)\*/g, '<em>$1</em>');
            html = html.replace(/___(.+?)___/g, '<strong><em>$1</em></strong>');
            html = html.replace(/__(.+?)__/g, '<strong>$1</strong>');
            html = html.replace(/_(.+?)_/g, '<em>$1</em>');

            // Inline code
            html = html.replace(/` + "`" + `(.+?)` + "`" + `/g, '<code>$1</code>');

            // Lists (simple version)
            html = html.replace(/^\d+\.\s+(.+)$/gm, '<li>$1</li>');
            html = html.replace(/^[-*+]\s+(.+)$/gm, '<li>$1</li>');

            // Wrap consecutive <li> tags in <ul>
            html = html.replace(/(<li>.*<\/li>\n?)+/g, function(match) {
                return '<ul>' + match + '</ul>';
            });

            // Blockquotes
            html = html.replace(/^&gt;\s+(.+)$/gm, '<blockquote>$1</blockquote>');

            // Line breaks and paragraphs
            html = html.replace(/\n\n/g, '</p><p>');
            html = html.replace(/\n/g, '<br>');

            // Wrap in paragraph if not already wrapped
            if (!html.startsWith('<')) {
                html = '<p>' + html + '</p>';
            }

            return html;
        }

        document.getElementById('uploadForm').addEventListener('submit', async (e) => {
            e.preventDefault();

            const fileInput = document.getElementById('logFile');
            const loading = document.getElementById('loading');
            const result = document.getElementById('result');
            const resultContent = document.getElementById('resultContent');

            if (!fileInput.files[0]) {
                alert('Please select a file');
                return;
            }

            const formData = new FormData();
            formData.append('logfile', fileInput.files[0]);

            loading.style.display = 'inline';
            result.style.display = 'none';

            try {
                const response = await fetch('/analyze', {
                    method: 'POST',
                    body: formData
                });

                const data = await response.json();

                if (response.ok) {
                    // Render Markdown to HTML
                    const htmlContent = markdownToHtml(data.analysis);
                    resultContent.innerHTML = htmlContent;
                    result.style.display = 'block';
                } else {
                    alert('Error: ' + (data.error || 'Unknown error'));
                }
            } catch (error) {
                alert('Error: ' + error.message);
            } finally {
                loading.style.display = 'none';
            }
        });
    </script>
</body>
</html>`

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status":"healthy"}`))
}
