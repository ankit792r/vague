// Package text implements buffer text storage and edit primitives.
package text

import (
	"bytes"
	"sort"
)

// Offset is a byte offset from the start of the buffer.
type Offset int

// Point is a buffer coordinate. Line is 0-based. Col is a byte offset within
// the line, not a rune index and not a display column.
type Point struct {
	Line int
	Col  int
}

// Compare orders points in buffer order.
func (p Point) Compare(q Point) int {
	switch {
	case p.Line != q.Line:
		return sign(p.Line - q.Line)
	default:
		return sign(p.Col - q.Col)
	}
}

// Delta describes a completed edit.
type Delta struct {
	Start  Offset
	OldEnd Offset
	NewEnd Offset

	StartPoint  Point
	OldEndPoint Point
	NewEndPoint Point

	Removed []byte
}

// Text stores buffer contents as an array of lines.
type Text struct {
	lines  [][]byte
	starts []Offset

	markers []*Marker
}

// New builds a Text from raw bytes. b is copied and not retained.
func New(b []byte) *Text {
	t := &Text{}
	t.SetBytes(b)
	return t
}

// SetBytes replaces the entire contents.
func (t *Text) SetBytes(b []byte) {
	t.lines = bytes.Split(bytes.Clone(b), []byte{'\n'})
	t.starts = nil
}

// LineCount returns the number of lines, always at least 1.
func (t *Text) LineCount() int { return len(t.lines) }

// Line returns the contents of line i without a terminating newline.
func (t *Text) Line(i int) []byte { return t.lines[i] }

// LineStart returns the offset of the first byte of line i.
func (t *Text) LineStart(i int) Offset {
	t.ensureStarts()
	return t.starts[i]
}

// LineEnd returns the offset just past the last byte of line i.
func (t *Text) LineEnd(i int) Offset {
	return t.LineStart(i) + Offset(len(t.lines[i]))
}

// Len returns the total size of the buffer in bytes.
func (t *Text) Len() Offset {
	last := len(t.lines) - 1
	return t.LineStart(last) + Offset(len(t.lines[last]))
}

// Bytes returns a copy of the whole buffer.
func (t *Text) Bytes() []byte {
	return bytes.Join(t.lines, []byte{'\n'})
}

// PointOf converts a byte offset to a line/column coordinate.
func (t *Text) PointOf(o Offset) Point {
	o = t.Clamp(o)
	t.ensureStarts()

	i := sort.Search(len(t.starts), func(i int) bool {
		return t.starts[i] > o
	}) - 1

	return Point{Line: i, Col: int(o - t.starts[i])}
}

// OffsetOf converts a line/column coordinate to a byte offset.
func (t *Text) OffsetOf(p Point) Offset {
	line := clamp(p.Line, 0, len(t.lines)-1)
	return t.LineStart(line) + Offset(clamp(p.Col, 0, len(t.lines[line])))
}

// Clamp restricts an offset to the buffer.
func (t *Text) Clamp(o Offset) Offset {
	if o < 0 {
		return 0
	}
	if end := t.Len(); o > end {
		return end
	}
	return o
}

// Replace substitutes ins for the bytes in [a, b) and returns what changed.
func (t *Text) Replace(a, b Offset, ins []byte) Delta {
	a, b = t.span(a, b)

	pa, pb := t.PointOf(a), t.PointOf(b)
	removed := t.Slice(a, b)

	head := t.lines[pa.Line][:pa.Col]
	tail := t.lines[pb.Line][pb.Col:]

	segs := bytes.Split(ins, []byte{'\n'})
	repl := make([][]byte, len(segs))

	if len(segs) == 1 {
		repl[0] = concat(head, segs[0], tail)
	} else {
		last := len(segs) - 1
		repl[0] = concat(head, segs[0])
		for i := 1; i < last; i++ {
			repl[i] = bytes.Clone(segs[i])
		}
		repl[last] = concat(segs[last], tail)
	}

	t.lines = splice(t.lines, pa.Line, pb.Line+1, repl)
	t.invalidateFrom(pa.Line)

	newEnd := a + Offset(len(ins))
	t.adjustMarkers(a, b, newEnd)

	return Delta{
		Start:       a,
		OldEnd:      b,
		NewEnd:      newEnd,
		StartPoint:  pa,
		OldEndPoint: pb,
		NewEndPoint: t.PointOf(newEnd),
		Removed:     removed,
	}
}

// Insert adds text at an offset.
func (t *Text) Insert(at Offset, ins []byte) Delta {
	return t.Replace(at, at, ins)
}

// Delete removes the bytes in [a, b).
func (t *Text) Delete(a, b Offset) Delta {
	return t.Replace(a, b, nil)
}

// Slice returns a copy of the bytes in [a, b).
func (t *Text) Slice(a, b Offset) []byte {
	a, b = t.span(a, b)
	if a == b {
		return nil
	}

	pa, pb := t.PointOf(a), t.PointOf(b)
	if pa.Line == pb.Line {
		return bytes.Clone(t.lines[pa.Line][pa.Col:pb.Col])
	}

	out := make([]byte, 0, b-a)
	out = append(out, t.lines[pa.Line][pa.Col:]...)
	for i := pa.Line + 1; i < pb.Line; i++ {
		out = append(out, '\n')
		out = append(out, t.lines[i]...)
	}
	out = append(out, '\n')
	return append(out, t.lines[pb.Line][:pb.Col]...)
}

func (t *Text) span(a, b Offset) (Offset, Offset) {
	if a > b {
		a, b = b, a
	}
	return t.Clamp(a), t.Clamp(b)
}

func (t *Text) ensureStarts() {
	if len(t.starts) == 0 {
		t.starts = append(t.starts, 0)
	}
	for i := len(t.starts); i < len(t.lines); i++ {
		t.starts = append(t.starts, t.starts[i-1]+Offset(len(t.lines[i-1]))+1)
	}
}

func (t *Text) invalidateFrom(line int) {
	if line < 1 {
		line = 1
	}
	if line < len(t.starts) {
		t.starts = t.starts[:line]
	}
}

func concat(parts ...[]byte) []byte {
	n := 0
	for _, p := range parts {
		n += len(p)
	}

	out := make([]byte, 0, n)
	for _, p := range parts {
		out = append(out, p...)
	}

	return out
}

func splice(dst [][]byte, from, to int, repl [][]byte) [][]byte {
	if len(repl) == to-from {
		copy(dst[from:to], repl)
		return dst
	}

	tail := append([][]byte(nil), dst[to:]...)
	out := append(dst[:from], repl...)
	return append(out, tail...)
}

func clamp(v, lo, hi int) int {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

func sign(n int) int {
	switch {
	case n < 0:
		return -1
	case n > 0:
		return 1
	default:
		return 0
	}
}
