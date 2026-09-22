package editor

import (
	"unicode/utf8"

	"vague/bonefire/text"
)

func isBlankLine(line []byte) bool {
	for _, b := range line {
		if b != ' ' && b != '\t' {
			return false
		}
	}
	return true
}

func paragraphStartLine(t *text.Text, line int) int {
	line = clampInt(line, 0, t.LineCount()-1)
	if isBlankLine(t.Line(line)) {
		return line
	}

	for line > 0 && !isBlankLine(t.Line(line-1)) {
		line--
	}
	return line
}

func paragraphEndLine(t *text.Text, line int) int {
	line = clampInt(line, 0, t.LineCount()-1)
	if isBlankLine(t.Line(line)) {
		return line
	}

	last := t.LineCount() - 1
	for line < last && !isBlankLine(t.Line(line+1)) {
		line++
	}
	return line
}

func moveParagraphForward(t *text.Text, off text.Offset, count int) text.Offset {
	if count < 1 {
		count = 1
	}

	line := t.PointOf(off).Line
	for i := 0; i < count; i++ {
		end := paragraphEndLine(t, line)
		next := end + 1
		for next < t.LineCount() && isBlankLine(t.Line(next)) {
			next++
		}
		if next >= t.LineCount() {
			line = t.LineCount() - 1
			break
		}
		line = next
	}

	return firstNonBlankOnLine(t, line)
}

func moveParagraphBackward(t *text.Text, off text.Offset, count int) text.Offset {
	if count < 1 {
		count = 1
	}

	line := t.PointOf(off).Line
	for i := 0; i < count; i++ {
		start := paragraphStartLine(t, line)
		if line > start {
			line = start
			continue
		}

		prev := start - 1
		for prev > 0 && isBlankLine(t.Line(prev)) {
			prev--
		}
		if start == 0 {
			line = 0
			break
		}
		line = paragraphStartLine(t, prev)
	}

	return firstNonBlankOnLine(t, line)
}

func isSentenceEnd(line []byte, col int) bool {
	if col >= len(line) {
		return false
	}

	switch line[col] {
	case '.', '!', '?':
	default:
		return false
	}

	next := col + 1
	if next >= len(line) {
		return true
	}

	return line[next] == ' ' || line[next] == '\t'
}

func skipSentenceGap(line []byte, col int) int {
	for col < len(line) && (line[col] == ' ' || line[col] == '\t') {
		col++
	}
	return col
}

func moveSentenceForward(t *text.Text, off text.Offset, count int) text.Offset {
	if count < 1 {
		count = 1
	}

	at := off
	for i := 0; i < count; i++ {
		next, ok := nextSentenceStart(t, at)
		if !ok {
			break
		}
		at = next
	}
	return at
}

func nextSentenceStart(t *text.Text, off text.Offset) (text.Offset, bool) {
	point := t.PointOf(off)

	for line := point.Line; line < t.LineCount(); line++ {
		lineBytes := t.Line(line)
		col := 0
		if line == point.Line {
			col = point.Col
			if col < len(lineBytes) {
				col++
			}
		}

		for col < len(lineBytes) {
			if isSentenceEnd(lineBytes, col) {
				after := skipSentenceGap(lineBytes, col+1)
				if after < len(lineBytes) {
					return t.OffsetOf(text.Point{Line: line, Col: after}), true
				}

				nextLine := line + 1
				if nextLine < t.LineCount() {
					return firstNonBlankOnLine(t, nextLine), true
				}
				return t.LineEnd(line), true
			}

			_, size := utf8.DecodeRune(lineBytes[col:])
			if size == 0 {
				col++
			} else {
				col += size
			}
		}
	}

	return off, false
}

func moveSentenceBackward(t *text.Text, off text.Offset, count int) text.Offset {
	if count < 1 {
		count = 1
	}

	at := off
	for i := 0; i < count; i++ {
		prev, ok := previousSentenceStart(t, at)
		if !ok {
			break
		}
		at = prev
	}
	return at
}

func previousSentenceStart(t *text.Text, off text.Offset) (text.Offset, bool) {
	if off <= 0 {
		return 0, false
	}

	point := t.PointOf(off)

	for line := point.Line; line >= 0; line-- {
		lineBytes := t.Line(line)
		endCol := len(lineBytes)
		if line == point.Line {
			endCol = point.Col
		}

		for col := endCol - 1; col >= 0; {
			if isSentenceEnd(lineBytes, col) {
				after := skipSentenceGap(lineBytes, col+1)
				start := t.OffsetOf(text.Point{Line: line, Col: after})
				if start < off {
					return start, true
				}
			}
			if col == 0 {
				break
			}
			_, size := utf8.DecodeLastRune(lineBytes[:col])
			if size <= 0 {
				col--
			} else {
				col -= size
			}
		}
	}

	lineStart := t.LineStart(point.Line)
	if off > lineStart {
		return lineStart, true
	}
	if point.Line > 0 {
		return firstNonBlankOnLine(t, point.Line-1), true
	}
	return 0, false
}
