package buffer

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"vague/bonefire/text"
)

// Load reads a file from disk into a new buffer.
func Load(path string) (*Buffer, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(abs)
	if err != nil {
		if os.IsNotExist(err) {
			return newEmptyFile(abs), nil
		}
		return nil, err
	}

	info, err := os.Stat(abs)
	if err != nil {
		return nil, err
	}

	normalized := normalizeNewlines(data)
	lineEnding := detectLineEnding(data)
	noEOL := len(data) > 0 && data[len(data)-1] != '\n'

	buf := &Buffer{
		Name:         filepath.Base(abs),
		Path:         abs,
		Text:         text.New(normalized),
		History:      text.NewUndoTree(),
		NoEOL:        noEOL,
		LineEnding:   lineEnding,
		Binary:       bytes.IndexByte(data, 0) >= 0,
		ReadOnly:     info.Mode()&0200 == 0,
		onDisk:       true,
		diskSize:     info.Size(),
		diskModTime:  info.ModTime(),
		startedEmpty: len(data) == 0,
	}

	return buf, nil
}

// Reload replaces buffer contents from its path on disk.
func (b *Buffer) Reload() error {
	if b.Path == "" {
		return ErrNoFileName
	}

	loaded, err := Load(b.Path)
	if err != nil {
		return err
	}

	b.Text.SetBytes(loaded.Text.Bytes())
	b.History = text.NewUndoTree()
	b.NoEOL = loaded.NoEOL
	b.LineEnding = loaded.LineEnding
	b.Binary = loaded.Binary
	b.ReadOnly = loaded.ReadOnly
	b.onDisk = loaded.onDisk
	b.diskSize = loaded.diskSize
	b.diskModTime = loaded.diskModTime
	b.startedEmpty = loaded.startedEmpty
	b.savedSeq = b.History.Seq()

	return nil
}

// Save writes the buffer to its path.
func (b *Buffer) Save() error {
	if b.Path == "" {
		return ErrNoFileName
	}

	return b.saveTo(b.Path, false)
}

// SaveForce writes the buffer even if the file changed on disk.
func (b *Buffer) SaveForce() error {
	if b.Path == "" {
		return ErrNoFileName
	}

	return b.saveTo(b.Path, true)
}

// SaveAs writes the buffer to path and makes that its backing file.
func (b *Buffer) SaveAs(path string) error {
	return b.saveAs(path, false)
}

// SaveAsForce writes to path even if the file changed on disk.
func (b *Buffer) SaveAsForce(path string) error {
	return b.saveAs(path, true)
}

func (b *Buffer) saveAs(path string, force bool) error {
	abs, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	if err := b.saveTo(abs, force); err != nil {
		return err
	}

	b.Path = abs
	b.Name = filepath.Base(abs)
	return nil
}

func (b *Buffer) saveTo(path string, force bool) error {
	if b.ReadOnly {
		return ErrReadOnly
	}

	if b.onDisk && !force {
		info, err := os.Stat(path)
		if err != nil {
			if !os.IsNotExist(err) {
				return err
			}
		} else if info.Size() != b.diskSize || !info.ModTime().Equal(b.diskModTime) {
			return ErrFileChanged
		}
	}

	data := b.encodeForDisk()

	dir := filepath.Dir(path)
	tmp, err := os.CreateTemp(dir, ".vague-save-*")
	if err != nil {
		return fmt.Errorf("create temp file: %w", err)
	}

	tmpName := tmp.Name()
	cleanup := true
	defer func() {
		if cleanup {
			_ = os.Remove(tmpName)
		}
	}()

	if _, err := tmp.Write(data); err != nil {
		_ = tmp.Close()
		return fmt.Errorf("write temp file: %w", err)
	}

	if err := tmp.Sync(); err != nil {
		_ = tmp.Close()
		return fmt.Errorf("sync temp file: %w", err)
	}

	if err := tmp.Close(); err != nil {
		return fmt.Errorf("close temp file: %w", err)
	}

	if err := os.Rename(tmpName, path); err != nil {
		return fmt.Errorf("rename temp file: %w", err)
	}

	cleanup = false

	info, err := os.Stat(path)
	if err != nil {
		return err
	}

	b.onDisk = true
	b.diskSize = info.Size()
	b.diskModTime = info.ModTime()
	b.savedSeq = b.History.Seq()
	return nil
}

// Modified reports whether the buffer has unsaved edits.
func (b *Buffer) Modified() bool {
	return b.History.Seq() != b.savedSeq
}

// ChangedOnDisk reports whether the file on disk differs from when we last read/saved.
func (b *Buffer) ChangedOnDisk() bool {
	if b.Path == "" {
		return false
	}
	info, err := os.Stat(b.Path)
	if err != nil {
		return false
	}
	return info.Size() != b.diskSize || !info.ModTime().Equal(b.diskModTime)
}

func (b *Buffer) encodeForDisk() []byte {
	if b.startedEmpty && b.Text.Len() == 0 {
		return nil
	}

	lines := make([][]byte, b.Text.LineCount())
	for i := 0; i < b.Text.LineCount(); i++ {
		lines[i] = bytes.Clone(b.Text.Line(i))
	}

	sep := []byte("\n")
	if b.LineEnding == CRLF {
		sep = []byte("\r\n")
	}

	out := bytes.Join(lines, sep)
	if b.NoEOL {
		return out
	}

	last := b.Text.LineCount() - 1
	if last >= 0 && len(b.Text.Line(last)) == 0 {
		return out
	}

	return append(out, sep...)
}

func normalizeNewlines(data []byte) []byte {
	if len(data) == 0 {
		return data
	}

	data = bytes.ReplaceAll(data, []byte("\r\n"), []byte("\n"))
	return bytes.ReplaceAll(data, []byte("\r"), []byte("\n"))
}

func detectLineEnding(data []byte) LineEnding {
	if bytes.Contains(data, []byte("\r\n")) {
		return CRLF
	}

	return LF
}

func newEmptyFile(abs string) *Buffer {
	return &Buffer{
		Name:         filepath.Base(abs),
		Path:         abs,
		Text:         text.New(nil),
		History:      text.NewUndoTree(),
		startedEmpty: true,
	}
}
