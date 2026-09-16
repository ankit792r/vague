package buffer

// LineEnding is how a file separates lines on disk.
//
// Line endings are normalised away at load and reapplied at save, so the
// text layer only ever sees "\n".
type LineEnding int

const (
	LF LineEnding = iota
	CRLF
)

func (l LineEnding) String() string {
	if l == CRLF {
		return "dos"
	}
	return "unix"
}
