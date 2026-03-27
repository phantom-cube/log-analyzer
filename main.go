package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"log-analyzer/handler"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using default configuration")
	}

	// Set gin mode based on environment
	ginMode := os.Getenv("GIN_MODE")
	if ginMode == "" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create gin router
	r := gin.Default()

	// Setup routes
	r.GET("/", handleHome)
	r.POST("/analyze", handler.AnalyzeLog)
	r.GET("/health", handleHealth)

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

	// Run server
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}

func handleHome(c *gin.Context) {
	html := `<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>AI 日志分析器 - 专业日志分析工具</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }

        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            min-height: 100vh;
            padding: 20px;
            line-height: 1.6;
        }

        .container {
            max-width: 900px;
            margin: 0 auto;
        }

        /* Header */
        .header {
            text-align: center;
            color: white;
            margin-bottom: 40px;
            animation: fadeInDown 0.6s ease-out;
        }

        .header h1 {
            font-size: 2.5rem;
            font-weight: 700;
            margin-bottom: 10px;
            text-shadow: 2px 2px 4px rgba(0,0,0,0.2);
        }

        .header p {
            font-size: 1.1rem;
            opacity: 0.95;
        }

        /* Upload Card */
        .upload-card {
            background: white;
            border-radius: 16px;
            padding: 40px;
            box-shadow: 0 20px 60px rgba(0,0,0,0.3);
            margin-bottom: 30px;
            animation: fadeInUp 0.6s ease-out;
        }

        .upload-card h2 {
            color: #333;
            margin-bottom: 10px;
            font-size: 1.5rem;
        }

        .file-info {
            color: #666;
            font-size: 0.9rem;
            margin-bottom: 25px;
            display: flex;
            align-items: center;
            gap: 8px;
        }

        .file-info .icon {
            color: #667eea;
        }

        /* Drag and Drop Area */
        .drop-zone {
            border: 2px dashed #cbd5e0;
            border-radius: 12px;
            padding: 50px 20px;
            text-align: center;
            transition: all 0.3s ease;
            cursor: pointer;
            background: #f7fafc;
            position: relative;
        }

        .drop-zone:hover,
        .drop-zone.drag-over {
            border-color: #667eea;
            background: #edf2f7;
            transform: translateY(-2px);
        }

        .drop-zone.drag-over {
            border-style: solid;
            background: #e6f0ff;
        }

        .upload-icon {
            font-size: 3rem;
            color: #667eea;
            margin-bottom: 15px;
        }

        .drop-zone p {
            color: #4a5568;
            font-size: 1rem;
            margin-bottom: 8px;
        }

        .drop-zone .file-types {
            color: #718096;
            font-size: 0.85rem;
        }

        .file-input {
            display: none;
        }

        .selected-file {
            margin-top: 20px;
            padding: 15px;
            background: #f0f9ff;
            border-radius: 8px;
            display: none;
            align-items: center;
            gap: 10px;
        }

        .selected-file.show {
            display: flex;
        }

        .selected-file .file-icon {
            font-size: 1.5rem;
        }

        .selected-file .file-details {
            flex: 1;
        }

        .selected-file .file-name {
            color: #333;
            font-weight: 500;
        }

        .selected-file .file-size {
            color: #666;
            font-size: 0.85rem;
        }

        .selected-file .remove-btn {
            background: #ef4444;
            color: white;
            border: none;
            padding: 6px 12px;
            border-radius: 6px;
            cursor: pointer;
            font-size: 0.85rem;
            transition: background 0.3s;
        }

        .selected-file .remove-btn:hover {
            background: #dc2626;
        }

        /* Analyze Button */
        .btn-analyze {
            width: 100%;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            border: none;
            padding: 16px 32px;
            font-size: 1.1rem;
            font-weight: 600;
            border-radius: 10px;
            cursor: pointer;
            margin-top: 25px;
            transition: all 0.3s ease;
            box-shadow: 0 4px 15px rgba(102, 126, 234, 0.4);
        }

        .btn-analyze:hover:not(:disabled) {
            transform: translateY(-2px);
            box-shadow: 0 6px 20px rgba(102, 126, 234, 0.6);
        }

        .btn-analyze:disabled {
            opacity: 0.6;
            cursor: not-allowed;
        }

        /* Loading Spinner */
        .loading-container {
            display: none;
            text-align: center;
            padding: 30px;
        }

        .loading-container.show {
            display: block;
        }

        .spinner {
            width: 50px;
            height: 50px;
            margin: 0 auto 20px;
            border: 4px solid #f3f3f3;
            border-top: 4px solid #667eea;
            border-radius: 50%;
            animation: spin 1s linear infinite;
        }

        .loading-text {
            color: #667eea;
            font-size: 1rem;
            font-weight: 500;
        }

        .loading-subtext {
            color: #718096;
            font-size: 0.9rem;
            margin-top: 8px;
        }

        /* Result Card */
        .result-card {
            background: white;
            border-radius: 16px;
            padding: 40px;
            box-shadow: 0 20px 60px rgba(0,0,0,0.3);
            display: none;
            animation: fadeInUp 0.6s ease-out;
        }

        .result-card.show {
            display: block;
        }

        .result-card h2 {
            color: #333;
            margin-bottom: 20px;
            font-size: 1.5rem;
            display: flex;
            align-items: center;
            gap: 10px;
        }

        .result-card .success-icon {
            color: #10b981;
            font-size: 1.8rem;
        }

        /* Error Alert */
        .error-alert {
            background: #fee;
            border-left: 4px solid #ef4444;
            padding: 15px 20px;
            border-radius: 8px;
            margin-bottom: 20px;
            display: none;
        }

        .error-alert.show {
            display: block;
            animation: shake 0.4s ease-out;
        }

        .error-alert strong {
            color: #dc2626;
        }

        .error-alert p {
            color: #991b1b;
            margin-top: 5px;
        }

        /* Markdown rendering styles */
        #resultContent {
            line-height: 1.8;
            color: #333;
        }

        #resultContent h1 {
            font-size: 1.8em;
            margin-top: 25px;
            margin-bottom: 15px;
            padding-bottom: 10px;
            border-bottom: 3px solid #667eea;
            color: #1a202c;
        }

        #resultContent h2 {
            font-size: 1.5em;
            margin-top: 22px;
            margin-bottom: 12px;
            color: #2d3748;
        }

        #resultContent h3 {
            font-size: 1.3em;
            margin-top: 18px;
            margin-bottom: 10px;
            color: #4a5568;
        }

        #resultContent strong {
            color: #dc2626;
            font-weight: 600;
        }

        #resultContent em {
            font-style: italic;
            color: #4b5563;
        }

        #resultContent code {
            background: #f3f4f6;
            padding: 3px 8px;
            border-radius: 4px;
            font-family: 'Monaco', 'Courier New', monospace;
            font-size: 0.9em;
            color: #be185d;
        }

        #resultContent pre {
            background: #1f2937;
            color: #f3f4f6;
            padding: 15px;
            border-radius: 8px;
            overflow-x: auto;
            margin: 15px 0;
        }

        #resultContent ul,
        #resultContent ol {
            margin: 15px 0;
            padding-left: 30px;
        }

        #resultContent li {
            margin: 8px 0;
        }

        #resultContent p {
            margin: 12px 0;
        }

        #resultContent blockquote {
            border-left: 4px solid #667eea;
            padding-left: 20px;
            margin: 15px 0;
            color: #4a5568;
            font-style: italic;
        }

        /* Animations */
        @keyframes fadeInDown {
            from {
                opacity: 0;
                transform: translateY(-20px);
            }
            to {
                opacity: 1;
                transform: translateY(0);
            }
        }

        @keyframes fadeInUp {
            from {
                opacity: 0;
                transform: translateY(20px);
            }
            to {
                opacity: 1;
                transform: translateY(0);
            }
        }

        @keyframes spin {
            0% { transform: rotate(0deg); }
            100% { transform: rotate(360deg); }
        }

        @keyframes shake {
            0%, 100% { transform: translateX(0); }
            25% { transform: translateX(-10px); }
            75% { transform: translateX(10px); }
        }

        /* Responsive Design */
        @media (max-width: 768px) {
            .header h1 {
                font-size: 2rem;
            }

            .upload-card,
            .result-card {
                padding: 25px;
            }

            .drop-zone {
                padding: 35px 15px;
            }
        }
    </style>
</head>
<body>
    <div class="container">
        <!-- Header -->
        <div class="header">
            <h1>🤖 AI 日志分析器</h1>
            <p>基于AI的专业日志文件分析工具</p>
        </div>

        <!-- Error Alert -->
        <div class="error-alert" id="errorAlert">
            <strong>⚠️ 错误</strong>
            <p id="errorMessage"></p>
        </div>

        <!-- Upload Card -->
        <div class="upload-card">
            <h2>📄 上传日志文件</h2>
            <div class="file-info">
                <span class="icon">ℹ️</span>
                <span>支持格式: .log, .txt | 最大文件大小: <strong>10 MB</strong></span>
            </div>

            <form id="uploadForm">
                <div class="drop-zone" id="dropZone">
                    <div class="upload-icon">📁</div>
                    <p><strong>将日志文件拖放到此处</strong></p>
                    <p>或</p>
                    <p style="color: #667eea; font-weight: 500; margin-top: 10px;">点击浏览文件</p>
                    <p class="file-types">接受 .log 和 .txt 文件</p>
                </div>
                <input type="file" id="logFile" class="file-input" accept=".log,.txt" required>

                <div class="selected-file" id="selectedFile">
                    <span class="file-icon">📄</span>
                    <div class="file-details">
                        <div class="file-name" id="fileName"></div>
                        <div class="file-size" id="fileSize"></div>
                    </div>
                    <button type="button" class="remove-btn" id="removeBtn">✕ 移除</button>
                </div>

                <button type="submit" class="btn-analyze" id="analyzeBtn" disabled>
                    🚀 开始分析日志
                </button>
            </form>

            <div class="loading-container" id="loadingContainer">
                <div class="spinner"></div>
                <div class="loading-text">正在分析您的日志文件...</div>
                <div class="loading-subtext">根据文件大小，这可能需要 2-5 分钟</div>
            </div>
        </div>

        <!-- Result Card -->
        <div class="result-card" id="resultCard">
            <h2>
                <span class="success-icon">✅</span>
                分析结果
            </h2>
            <div id="resultContent"></div>
        </div>
    </div>

    <script>
        // Constants
        const MAX_FILE_SIZE = 10 * 1024 * 1024; // 10MB in bytes

        // DOM Elements
        const dropZone = document.getElementById('dropZone');
        const fileInput = document.getElementById('logFile');
        const uploadForm = document.getElementById('uploadForm');
        const selectedFile = document.getElementById('selectedFile');
        const fileName = document.getElementById('fileName');
        const fileSize = document.getElementById('fileSize');
        const removeBtn = document.getElementById('removeBtn');
        const analyzeBtn = document.getElementById('analyzeBtn');
        const loadingContainer = document.getElementById('loadingContainer');
        const resultCard = document.getElementById('resultCard');
        const resultContent = document.getElementById('resultContent');
        const errorAlert = document.getElementById('errorAlert');
        const errorMessage = document.getElementById('errorMessage');

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

        // Format file size
        function formatFileSize(bytes) {
            if (bytes < 1024) return bytes + ' B';
            if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(2) + ' KB';
            return (bytes / (1024 * 1024)).toFixed(2) + ' MB';
        }

        // Show error
        function showError(message) {
            errorMessage.textContent = message;
            errorAlert.classList.add('show');
            setTimeout(() => {
                errorAlert.classList.remove('show');
            }, 5000);
        }

        // Validate file
        function validateFile(file) {
            // Check file type
            const validTypes = ['.log', '.txt'];
            const fileExt = '.' + file.name.split('.').pop().toLowerCase();
            if (!validTypes.includes(fileExt)) {
                showError('文件类型无效。请上传 .log 或 .txt 文件。');
                return false;
            }

            // Check file size
            if (file.size > MAX_FILE_SIZE) {
                showError('文件太大。最大文件大小为 10 MB，您的文件大小为 ' + formatFileSize(file.size) + '。');
                return false;
            }

            // Check if file is empty
            if (file.size === 0) {
                showError('文件为空。请上传包含内容的文件。');
                return false;
            }

            return true;
        }

        // Handle file selection
        function handleFileSelect(file) {
            if (!validateFile(file)) {
                fileInput.value = '';
                return;
            }

            // Display selected file info
            fileName.textContent = file.name;
            fileSize.textContent = formatFileSize(file.size);
            selectedFile.classList.add('show');
            analyzeBtn.disabled = false;
        }

        // Click on drop zone to select file
        dropZone.addEventListener('click', () => {
            fileInput.click();
        });

        // File input change
        fileInput.addEventListener('change', (e) => {
            if (e.target.files.length > 0) {
                handleFileSelect(e.target.files[0]);
            }
        });

        // Remove file
        removeBtn.addEventListener('click', (e) => {
            e.stopPropagation();
            fileInput.value = '';
            selectedFile.classList.remove('show');
            analyzeBtn.disabled = true;
        });

        // Drag and drop events
        ['dragenter', 'dragover', 'dragleave', 'drop'].forEach(eventName => {
            dropZone.addEventListener(eventName, (e) => {
                e.preventDefault();
                e.stopPropagation();
            });
        });

        ['dragenter', 'dragover'].forEach(eventName => {
            dropZone.addEventListener(eventName, () => {
                dropZone.classList.add('drag-over');
            });
        });

        ['dragleave', 'drop'].forEach(eventName => {
            dropZone.addEventListener(eventName, () => {
                dropZone.classList.remove('drag-over');
            });
        });

        dropZone.addEventListener('drop', (e) => {
            const files = e.dataTransfer.files;
            if (files.length > 0) {
                fileInput.files = files;
                handleFileSelect(files[0]);
            }
        });

        // Form submission
        uploadForm.addEventListener('submit', async (e) => {
            e.preventDefault();

            if (!fileInput.files[0]) {
                showError('请选择文件');
                return;
            }

            const formData = new FormData();
            formData.append('logfile', fileInput.files[0]);

            // Show loading, hide results
            uploadForm.style.display = 'none';
            loadingContainer.classList.add('show');
            resultCard.classList.remove('show');

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
                    resultCard.classList.add('show');
                } else {
                    showError(data.error || '分析失败，请重试。');
                    uploadForm.style.display = 'block';
                }
            } catch (error) {
                showError('网络错误: ' + error.message + '。请检查网络连接后重试。');
                uploadForm.style.display = 'block';
            } finally {
                loadingContainer.classList.remove('show');
            }
        });
    </script>
</body>
</html>`

	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}

func handleHealth(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "healthy",
	})
}
