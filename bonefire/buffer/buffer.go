package buffer

import (
	"time"
	"vague/bonefire/text"
)

type Buffer struct {
	ID   uint64
	Name string
	Path string

	Text *text.Text
	// History *text.UndoTree

	NoEOL bool
	// LineEnding LineEnding
	Binary   bool
	ReadOnly bool

	startedEmpty bool
	onDisk       bool
	diskSize     int64
	diskModTime  time.Time
	savedSeq     int
}

// NewScratch creates a buffer with no backing file.
func NewScratch(id uint64, name string) *Buffer {
	initial := []byte("This is scratch buffer\nModified contents are not saved.")
	buf := &Buffer{
		ID:   id,
		Name: name,
		Text: text.New(initial),
		// History: text.NewUndoTree(),
	}
	// buf.History.SetInitial(initial)
	return buf
}
