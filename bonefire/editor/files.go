package editor

import (
	"fmt"
	"os"
	"path/filepath"

	"errors"

	"vague/bonefire/buffer"
)

func (e *Editor) LoadBufferPath(path string) (*buffer.Buffer, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	return e.LoadBufferPathAt("", abs)
}

// LoadBufferPathAt resolves path against workDir and loads or reuses a buffer.
func (e *Editor) LoadBufferPathAt(workDir, path string) (*buffer.Buffer, error) {
	abs, err := resolvePathAgainst(workDir, path)
	if err != nil {
		return nil, err
	}

	if existing := e.FindBuffer(abs); existing != nil {
		e.currentBuffer = existing.ID
		return existing, nil
	}

	loaded, err := buffer.Load(abs)
	if err != nil {
		return nil, err
	}

	loaded.ID = e.nextBufferID
	e.nextBufferID++
	e.Buffers[loaded.ID] = loaded
	e.currentBuffer = loaded.ID

	return loaded, nil
}

func resolvePathAgainst(baseDir, path string) (string, error) {
	if path == "" {
		return "", fmt.Errorf("empty path")
	}

	if filepath.IsAbs(path) {
		return filepath.Clean(path), nil
	}

	if baseDir == "" {
		if home, err := os.UserHomeDir(); err == nil {
			baseDir = home
		} else {
			baseDir = "/"
		}
	}

	return filepath.Abs(filepath.Join(baseDir, path))
}

// BufferInfo builds the execute result for buffer commands.
func BufferInfo(buf *buffer.Buffer, current bool) map[string]any {
	return map[string]any{
		"id":       buf.ID,
		"name":     buf.Name,
		"path":     buf.Path,
		"modified": buf.Modified(),
		"readonly": buf.ReadOnly,
		"current":  current,
		"lines":    buf.Text.LineCount(),
		"bytes":    int(buf.Text.Len()),
	}
}

// OpenFileError wraps an open failure for the wire.
func OpenFileError(path string, err error) error {
	return fmt.Errorf("edit %q: %w", path, err)
}

// WriteFileError wraps a save failure for the wire.
func WriteFileError(path string, err error) error {
	if errors.Is(err, buffer.ErrNoFileName) {
		return fmt.Errorf("No file name")
	}

	if path == "" {
		return fmt.Errorf("write: %w", err)
	}

	return fmt.Errorf("write %q: %w", path, err)
}
