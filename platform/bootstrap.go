package platform

import (
	"embed"
	"fmt"
	"io/fs"
	"net"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/abemedia/go-webview"
	_ "github.com/abemedia/go-webview/embedded"
)

//go:embed frontend/dist/*
var frontend embed.FS

func BootUI() error {
	// _, err := loadUI()
	// if err != nil {
	// 	return err
	// }

	w := webview.New(true)
	w.SetTitle("Vague")
	w.SetSize(1200, 800, webview.HintNone)
	// w.Navigate("http://" + addr)
	w.Navigate("http://localhost:5173")
	w.Run()

	return nil
}

func loadUI() (string, error) {
	dist, err := fs.Sub(frontend, "frontend/dist")
	if err != nil {
		return "", err
	}

	// Find an available local port.
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return "", err
	}

	addr := listener.Addr().String()

	// Serve embedded files.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestPath := strings.TrimPrefix(
			path.Clean(r.URL.Path),
			"/",
		)

		if requestPath == "" || requestPath == "." {
			requestPath = "index.html"
		}

		// Try requested file.
		if _, err := fs.Stat(dist, requestPath); err == nil {
			http.FileServer(http.FS(dist)).ServeHTTP(w, r)
			return
		}

		// SPA fallback.
		index, err := fs.ReadFile(dist, "index.html")
		if err != nil {
			http.Error(
				w,
				"index.html not found",
				http.StatusInternalServerError,
			)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(index)
	})

	server := &http.Server{
		Handler: handler,
	}

	go func() {
		if err := server.Serve(listener); err != nil &&
			err != http.ErrServerClosed {
			fmt.Fprintln(os.Stderr, err)
		}
	}()

	return addr, nil
}
