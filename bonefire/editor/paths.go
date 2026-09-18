package editor

import (
	"os"
	"path/filepath"
)

func initialWorkDir(clientDir string, initialPaths ...string) string {
	if len(initialPaths) > 0 && initialPaths[0] != "" {
		if abs, err := resolvePathAgainst(defaultWorkDir(clientDir), initialPaths[0]); err == nil {
			return filepath.Dir(abs)
		}
	}

	return defaultWorkDir(clientDir)
}

func defaultWorkDir(clientDir string) string {
	if clientDir != "" {
		if abs, err := filepath.Abs(clientDir); err == nil && abs != "/" {
			return abs
		}
	}

	if home, err := os.UserHomeDir(); err == nil {
		return home
	}

	return "/"
}

func resolvePathAgainst(baseDir, path string) (string, error) {
	if path == "" {
		return "", nil
	}

	if filepath.IsAbs(path) {
		return filepath.Clean(path), nil
	}

	if baseDir == "" {
		baseDir = defaultWorkDir("")
	}

	return filepath.Abs(filepath.Join(baseDir, path))
}

func (e *Editor) resolvePath(frameID uint64, path string) (string, error) {
	frame, ok := e.Frames[frameID]
	if !ok {
		return "", errNotFound("frame", frameID)
	}

	return resolvePathAgainst(frame.WorkDir, path)
}
