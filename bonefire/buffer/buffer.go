package buffer

import "errors"

var (
	// ErrNoFileName is returned when saving a buffer that has no path.
	ErrNoFileName = errors.New("buffer has no file name")

	// ErrReadOnly is returned when editing a buffer marked read-only.
	ErrReadOnly = errors.New("buffer is read-only")

	// ErrFileChanged means the file changed on disk since it was read.
	// Callers decide whether to overwrite; :w! is the escape hatch.
	ErrFileChanged = errors.New("file changed on disk since it was read")
)

// A Buffer is a file's contents in memory, together with its edit history
// and everything needed to write it back faithfully.
type Buffer struct {
	ID   uint64
	Name string

	// Path is empty for scratch buffers.
	Path string

	Text string // FIXME: this is temp, only until we dont have text

	// Text    *text.Text
	// History *text.UndoTree

	// NoEOL records that the file did not end with a newline. Preserving
	// this is what keeps saving an untouched file byte-identical instead of
	// producing a spurious one-line diff.
	NoEOL      bool
	LineEnding LineEnding
	Binary     bool
	ReadOnly   bool

	// startedEmpty records that the buffer began with no content at all,
	// either from a zero-byte file or from a file that does not exist yet.
	startedEmpty bool

	// Marks are the 'a-'z positions.
	// Marks map[rune]*text.Marker

	// LastCursor is where the cursor was when this buffer stopped being
	// displayed, so showing it again resumes in the same place instead of
	// jumping to the top.
	// LastCursor text.Offset

	// Options Options

	// Disk state, for detecting changes made behind the editor's back.
	onDisk bool
	// diskMtime time.Time
	diskSize int64

	// savedSeq is the undo sequence at the last successful write.
	savedSeq int
}

// NewScratch creates a buffer with no backing file.
func NewScratch(id uint64, name string) *Buffer {
	return &Buffer{
		ID:   id,
		Name: name,
		Text: "This is scratch buffer\nModified contents are not saved.",
		// Text:    text.New(nil),
		// History: text.NewUndoTree(),
		// Marks:   make(map[rune]*text.Marker),
		// Options: DefaultOptions(),
	}
}
