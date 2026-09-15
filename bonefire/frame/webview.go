package frame

import (
	"embed"
	"encoding/json"
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

// Create new web view frame
func (f *Frame) BuildWebView() error {
	w := webview.New(true)
	defer w.Destroy()

	w.Init(`
	window.hostEvent = (() => {
		const listeners = new Set();
		return {
			subscribe(fn) {
				listeners.add(fn);
				return () => listeners.delete(fn);
			},
			_emit(event, payload) {
				const msg = { event, payload };
				listeners.forEach(fn => fn(msg));
			}
		};
	})();
	`)

	go func() {
		for note := range f.client.Notifications() {
			var payload any
			_ = json.Unmarshal(note.Params, &payload)
			emitHostEvent(w, note.Method, payload) // "redraw", "quit", etc.
		}
	}()

	if err := w.Bind("hostRequest", f.HostRequest); err != nil {
		return fmt.Errorf("Host Request Binding Failed: %w", err)
	}

	w.SetTitle("Vague")
	w.SetSize(1200, 800, webview.HintNone)
	w.Navigate("http://localhost:5173")
	w.Run()

	return nil
}

//go:embed	output
var uiOutput embed.FS

func loadStaticUI() (string, error) {
	dist, err := fs.Sub(uiOutput, "output")
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

func emitHostEvent(w webview.WebView, event string, payload any) {
	data, err := json.Marshal(payload)
	if err != nil {
		return
	}
	js := fmt.Sprintf(
		`window.hostEvent && window.hostEvent._emit(%q, %s)`,
		event,
		string(data),
	)
	w.Dispatch(func() { w.Eval(js) })
}
