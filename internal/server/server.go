package server

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/ghchinoy/montage/ui"
)

// StartServer spins up the Montage web service and SPA static file server.
func StartServer(port int) error {
	mux := http.NewServeMux()

	// Health check endpoint (for Cloud Run & k8s)
	mux.HandleFunc("/health", handleHealth)

	// API Endpoints
	mux.HandleFunc("/api/assess", handleAssess)
	mux.HandleFunc("/api/export", handleExportZip)

	// UI & SPA Static Handler
	distFS, hasEmbed := ui.DistFS()
	if hasEmbed {
		fileServer := http.FileServer(http.FS(distFS))
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			path := strings.TrimPrefix(r.URL.Path, "/")
			if path == "" {
				path = "index.html"
			}

			// Try to open file directly
			f, err := distFS.Open(path)
			if err == nil {
				f.Close()
				fileServer.ServeHTTP(w, r)
				return
			}

			// SPA fallback: return index.html
			indexFile, err := distFS.Open("index.html")
			if err != nil {
				http.NotFound(w, r)
				return
			}
			defer indexFile.Close()

			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			_, _ = io.Copy(w, indexFile)
		})
	} else {
		// Fallback when UI has not been built
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			fmt.Fprintf(w, `<!doctype html>
<html>
<head><title>Montage API Service</title><style>body{font-family:sans-serif;background:#090d16;color:#f8fafc;padding:3rem;text-align:center;}</style></head>
<body>
  <h1>🎞️ Montage API Service</h1>
  <p>The backend service is running, but UI assets were not found.</p>
  <p>Run <code>npm run build</code> in <code>ui/</code> or rebuild with <code>go build -tags embedui</code>.</p>
  <p><a href="/health" style="color:#10b981;">Check /health endpoint</a></p>
</body>
</html>`)
		})
	}

	addr := fmt.Sprintf(":%d", port)
	fmt.Printf("🚀 Montage Web Server listening on http://0.0.0.0%s\n", addr)
	if hasEmbed {
		fmt.Println("✓ Serving embedded Lit WebComponent production assets")
	} else {
		fmt.Println("ℹ️  Running in API mode (no embedded UI assets)")
	}

	return http.ListenAndServe(addr, mux)
}

// ResolvePort returns the port to listen on, respecting $PORT (Cloud Run standard).
func ResolvePort(flagPort int) int {
	if envPort := os.Getenv("PORT"); envPort != "" {
		var p int
		if _, err := fmt.Sscanf(envPort, "%d", &p); err == nil && p > 0 {
			return p
		}
	}
	if flagPort > 0 {
		return flagPort
	}
	return 8080
}
