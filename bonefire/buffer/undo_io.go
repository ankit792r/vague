package buffer

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"

	"vague/bonefire/text"
)

func undoCachePath(filePath string) (string, error) {
	dir, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256([]byte(filePath))
	name := hex.EncodeToString(sum[:]) + ".json"
	return filepath.Join(dir, "vague", "undo", name), nil
}

// LoadUndoHistory reads a persisted undo tree for path, if present.
func LoadUndoHistory(filePath string, initial []byte) *text.UndoTree {
	path, err := undoCachePath(filePath)
	if err != nil {
		tree := text.NewUndoTree()
		tree.SetInitial(initial)
		return tree
	}
	data, err := os.ReadFile(path)
	if err != nil {
		tree := text.NewUndoTree()
		tree.SetInitial(initial)
		return tree
	}
	tree, err := text.UnmarshalUndo(data)
	if err != nil {
		tree = text.NewUndoTree()
		tree.SetInitial(initial)
		return tree
	}
	replayed, _, ok := tree.ReplayToSeq(tree.Current().Seq)
	if !ok || !bytes.Equal(replayed, initial) {
		tree = text.NewUndoTree()
		tree.SetInitial(initial)
		return tree
	}
	return tree
}

// SaveUndoHistory writes the buffer undo tree for its path.
func SaveUndoHistory(b *Buffer) error {
	if b.Path == "" {
		return nil
	}
	path, err := undoCachePath(b.Path)
	if err != nil {
		return err
	}
	data, err := b.History.MarshalUndo()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

// RestoreUndoSeq replays the buffer text to an undo sequence.
func (b *Buffer) RestoreUndoSeq(seq int) (text.Offset, bool) {
	out, cur, ok := b.History.ReplayToSeq(seq)
	if !ok {
		return 0, false
	}
	b.Text.SetBytes(out)
	return cur, true
}
