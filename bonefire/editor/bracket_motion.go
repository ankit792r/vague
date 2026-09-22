package editor

import (
	"vague/bonefire/text"
)

var bracketClose = map[byte]byte{
	')': '(',
	']': '[',
	'}': '{',
}

var bracketOpen = map[byte]byte{
	'(': ')',
	'[': ']',
	'{': '}',
}

func isBracket(b byte) bool {
	_, ok := bracketOpen[b]
	if ok {
		return true
	}
	_, ok = bracketClose[b]
	return ok
}

// moveMatchingBracket finds the matching bracket for % (Vim-style, same line first then buffer).
func moveMatchingBracket(t *text.Text, off text.Offset, count int) text.Offset {
	if count < 1 {
		count = 1
	}

	at := off
	for i := 0; i < count; i++ {
		next, ok := findMatchingBracket(t, at)
		if !ok {
			break
		}
		at = next
	}
	return at
}

func findMatchingBracket(t *text.Text, off text.Offset) (text.Offset, bool) {
	point := t.PointOf(off)
	line := t.Line(point.Line)

	// Prefer bracket at or after cursor on this line.
	col := point.Col
	if col >= len(line) {
		col = len(line) - 1
	}
	if col < 0 {
		return off, false
	}

	if !isBracket(line[col]) {
		found := -1
		for i := point.Col; i < len(line); i++ {
			if isBracket(line[i]) {
				found = i
				break
			}
		}
		if found < 0 {
			for i := point.Col - 1; i >= 0; i-- {
				if isBracket(line[i]) {
					found = i
					break
				}
			}
		}
		if found < 0 {
			return off, false
		}
		col = found
	}

	ch := line[col]
	if close, isOpen := bracketOpen[ch]; isOpen {
		if match, ok := scanForwardBalance(t, point.Line, col+1, ch, close); ok {
			return t.OffsetOf(text.Point{Line: match.line, Col: match.col}), true
		}
		return off, false
	}

	if open, isClose := bracketClose[ch]; isClose {
		if match, ok := scanBackwardBalance(t, point.Line, col-1, open, ch); ok {
			return t.OffsetOf(text.Point{Line: match.line, Col: match.col}), true
		}
		return off, false
	}

	return off, false
}

type lineCol struct {
	line int
	col  int
}

func scanForwardBalance(t *text.Text, startLine, startCol int, open, close byte) (lineCol, bool) {
	depth := 1
	for line := startLine; line < t.LineCount(); line++ {
		lineBytes := t.Line(line)
		col := 0
		if line == startLine {
			col = startCol
		}
		for col < len(lineBytes) {
			switch lineBytes[col] {
			case open:
				depth++
			case close:
				depth--
				if depth == 0 {
					return lineCol{line, col}, true
				}
			}
			col++
		}
	}
	return lineCol{}, false
}

func scanBackwardBalance(t *text.Text, startLine, startCol int, open, close byte) (lineCol, bool) {
	depth := 1
	for line := startLine; line >= 0; line-- {
		lineBytes := t.Line(line)
		col := len(lineBytes) - 1
		if line == startLine {
			col = startCol
		}
		for col >= 0 {
			switch lineBytes[col] {
			case close:
				depth++
			case open:
				depth--
				if depth == 0 {
					return lineCol{line, col}, true
				}
			}
			col--
		}
	}
	return lineCol{}, false
}
