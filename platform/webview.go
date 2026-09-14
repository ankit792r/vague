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
	"vague/backbone/process"

	"github.com/abemedia/go-webview"
	_ "github.com/abemedia/go-webview/embedded"
)

// Create new web view frame
func (f *Frame) BuildWebView() error {
	w := webview.New(true)
	defer func() {
		if err := f.client.FrameDetach(f.ctx); err != nil {
			fmt.Errorf("Error Detach: %w", err)
		}
		w.Destroy()
	}()

	if err := w.Bind("ipcBinding", f.IpcBinding); err != nil {
		return fmt.Errorf("Binding Ipc Failed: %w", err)
	}

	res, err := f.client.FrameAttach(f.ctx, process.AttachParams{})
	if err != nil {
		return err
	}

	fmt.Printf("Got res: %s\n", res)

	w.SetTitle("Vague")
	w.SetSize(1200, 800, webview.HintNone)
	w.Navigate("http://localhost:5173")
	w.Run()

	return nil
}

//go:embed frontend/dist/*
var uiDist embed.FS

func loadStaticUI() (string, error) {
	dist, err := fs.Sub(uiDist, "frontend/dist")
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
