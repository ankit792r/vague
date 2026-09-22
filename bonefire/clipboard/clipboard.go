package clipboard

import "sync"

var (
	mu              sync.RWMutex
	primaryContent  string
	clipboardContent string
)

// ReadPrimary returns the star-register / primary selection (stub: in-memory).
func ReadPrimary() (string, bool) {
	mu.RLock()
	defer mu.RUnlock()
	if primaryContent == "" {
		return "", false
	}
	return primaryContent, true
}

// ReadClipboard returns the plus-register / system clipboard (stub: in-memory).
func ReadClipboard() (string, bool) {
	mu.RLock()
	defer mu.RUnlock()
	if clipboardContent == "" {
		return "", false
	}
	return clipboardContent, true
}

// WritePrimary stores content for the * register.
func WritePrimary(text string) {
	mu.Lock()
	primaryContent = text
	mu.Unlock()
}

// WriteClipboard stores content for the + register.
func WriteClipboard(text string) {
	mu.Lock()
	clipboardContent = text
	mu.Unlock()
}

// WriteOSC52 emits clipboard data for terminal hosts (stub for webview).
func WriteOSC52(_ []byte) error {
	return nil
}
