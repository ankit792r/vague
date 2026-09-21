package editor

import (
	"unicode"
	"unicode/utf8"

	"vague/bonefire/text"
)

type wordClass int

const (
	wcWhitespace wordClass = iota
	wcKeyword
	wcPunct
)

func classifyRune(r rune) wordClass {
	switch {
	case r == ' ' || r == '\t' || r == '\n':
		return wcWhitespace
	case r == '_' || unicode.IsLetter(r) || unicode.IsDigit(r):
		return wcKeyword
	default:
		return wcPunct
	}
}

// charAt returns the word class at off, the offset of that character, and the
// offset just past it.
func charAt(t *text.Text, off text.Offset) (wordClass, text.Offset, text.Offset) {
	if off >= t.Len() {
		return wcWhitespace, off, off
	}

	point := t.PointOf(off)
	line := t.Line(point.Line)
	if point.Col >= len(line) {
		if point.Line+1 < t.LineCount() {
			next := t.LineStart(point.Line + 1)
			return wcWhitespace, off, next
		}
		return wcWhitespace, off, off
	}

	r, size := utf8.DecodeRune(line[point.Col:])
	next := t.OffsetOf(text.Point{Line: point.Line, Col: point.Col + size})
	return classifyRune(r), off, next
}

func moveWordForward(t *text.Text, off text.Offset, count int) text.Offset {
	for i := 0; i < count; i++ {
		next := moveWordForwardOnce(t, off)
		if next == off {
			break
		}
		off = next
	}
	return off
}

// wordEndExclusive is the end of the word at off, without trailing whitespace.
func wordEndExclusive(t *text.Text, off text.Offset) text.Offset {
	if off >= t.Len() {
		return off
	}

	cls, _, _ := charAt(t, off)
	if cls == wcWhitespace {
		return off
	}

	end := off
	for end < t.Len() {
		c, _, n := charAt(t, end)
		if c != cls {
			break
		}
		end = n
	}

	return end
}

func moveWordForwardOnce(t *text.Text, off text.Offset) text.Offset {
	start := off
	if off >= t.Len() {
		return off
	}

	cls, _, next := charAt(t, off)
	if cls == wcKeyword || cls == wcPunct {
		off = next
		for off < t.Len() {
			c, _, n := charAt(t, off)
			if c != cls {
				break
			}
			off = n
		}
	}

	for off < t.Len() {
		c, _, n := charAt(t, off)
		if c != wcWhitespace {
			break
		}
		off = n
	}

	if off >= t.Len() {
		return start
	}

	return off
}

func moveWordBack(t *text.Text, off text.Offset, count int) text.Offset {
	for i := 0; i < count; i++ {
		next := moveWordBackOnce(t, off)
		if next == off {
			break
		}
		off = next
	}
	return off
}

func moveWordBackOnce(t *text.Text, off text.Offset) text.Offset {
	if off <= 0 {
		return 0
	}

	off = back(t, off)

	for off > 0 {
		cls, _, _ := charAt(t, off)
		if cls != wcWhitespace {
			break
		}
		off = back(t, off)
	}

	if off <= 0 {
		return 0
	}

	cls, _, _ := charAt(t, off)
	for off > 0 {
		prev := back(t, off)
		c, _, _ := charAt(t, prev)
		if c != cls {
			break
		}
		off = prev
	}

	return off
}

func moveWordEnd(t *text.Text, off text.Offset, count int) text.Offset {
	for i := 0; i < count; i++ {
		next := moveWordEndOnce(t, off)
		if next == off {
			break
		}
		off = next
	}
	return off
}

func moveWordEndOnce(t *text.Text, off text.Offset) text.Offset {
	start := off
	if off >= t.Len() {
		return lastBufferChar(t)
	}

	cls, _, next := charAt(t, off)
	if cls == wcWhitespace {
		off = next
		for off < t.Len() {
			c, _, n := charAt(t, off)
			if c != wcWhitespace {
				break
			}
			off = n
		}
		if off >= t.Len() {
			return start
		}
		cls, _, _ = charAt(t, off)
	} else {
		c, _, _ := charAt(t, next)
		if c != cls || next >= t.Len() {
			off = next
			for off < t.Len() {
				c2, _, n2 := charAt(t, off)
				if c2 != wcWhitespace {
					break
				}
				off = n2
			}
			if off >= t.Len() {
				return lastBufferChar(t)
			}
			cls, _, _ = charAt(t, off)
		}
	}

	end := off
	for end < t.Len() {
		c, _, n := charAt(t, end)
		if c != cls {
			break
		}
		last := end
		end = n
		if end >= t.Len() {
			return last
		}
		c2, _, _ := charAt(t, end)
		if c2 != cls {
			return last
		}
	}

	return end
}

func lastBufferChar(t *text.Text) text.Offset {
	lastLine := t.LineCount() - 1
	line := t.Line(lastLine)
	if len(line) == 0 {
		return t.LineStart(lastLine)
	}

	return t.OffsetOf(text.Point{Line: lastLine, Col: lastColumn(line)})
}
