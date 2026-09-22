package shada

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// FileMark is a persisted uppercase file mark (A–Z).
type FileMark struct {
	Path string `json:"path"`
	Line int    `json:"line"`
	Col  int    `json:"col"`
}

// Store holds session data written to disk (shada subset).
type Store struct {
	FileMarks map[string]FileMark `json:"file_marks"`
}

func defaultStore() *Store {
	return &Store{FileMarks: make(map[string]FileMark)}
}

// Path returns the default shada file location.
func Path() (string, error) {
	dir, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "vague", "shada.json"), nil
}

// Load reads the shada file or returns an empty store.
func Load() (*Store, error) {
	path, err := Path()
	if err != nil {
		return defaultStore(), nil
	}
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return defaultStore(), nil
		}
		return nil, err
	}
	var s Store
	if err := json.Unmarshal(data, &s); err != nil {
		return defaultStore(), nil
	}
	if s.FileMarks == nil {
		s.FileMarks = make(map[string]FileMark)
	}
	return &s, nil
}

// Save writes the store to disk.
func (s *Store) Save() error {
	if s == nil {
		return nil
	}
	path, err := Path()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}
