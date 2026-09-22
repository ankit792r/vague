package editor

import "vague/bonefire/text"

func (e *Editor) resolveTextObjectKeys(keys string) (combined string, wait bool) {
	if e.pendingTextObject != "" {
		combined = e.pendingTextObject + keys
		e.pendingTextObject = ""
		return combined, false
	}
	if (keys == "i" || keys == "a") && e.wantsTextObjectKey() {
		e.pendingTextObject = keys
		return "", true
	}
	return keys, false
}

func (e *Editor) wantsTextObjectKey() bool {
	return e.pendingOp != opNone ||
		e.pendingCaseChange != caseChangeNone ||
		e.pendingFormat
}

func quotedObjectRange(t *text.Text, at text.Offset, quote byte, inner bool) (from, to text.Offset, ok bool) {
	point := t.PointOf(at)
	line := t.Line(point.Line)
	col := point.Col
	if col >= len(line) {
		if len(line) > 0 {
			col = len(line) - 1
		}
	}

	openCol := -1
	for c := col; c >= 0; c-- {
		if line[c] == quote {
			openCol = c
			break
		}
	}
	if openCol < 0 {
		for c := col + 1; c < len(line); c++ {
			if line[c] == quote {
				openCol = c
				break
			}
		}
	}
	if openCol < 0 {
		return 0, 0, false
	}

	closeCol := -1
	for c := openCol + 1; c < len(line); c++ {
		if line[c] == quote {
			closeCol = c
			break
		}
	}
	if closeCol < 0 {
		return 0, 0, false
	}

	from = t.OffsetOf(text.Point{Line: point.Line, Col: openCol + 1})
	to = t.OffsetOf(text.Point{Line: point.Line, Col: closeCol})
	if inner {
		if to <= from {
			return 0, 0, false
		}
		return from, to, true
	}
	from = t.OffsetOf(text.Point{Line: point.Line, Col: openCol})
	to = t.OffsetOf(text.Point{Line: point.Line, Col: closeCol + 1})
	if to <= from {
		return 0, 0, false
	}
	return from, to, true
}

func pairObjectRange(t *text.Text, at text.Offset, open, close byte, inner bool) (from, to text.Offset, ok bool) {
	point := t.PointOf(at)
	line := t.Line(point.Line)
	col := point.Col
	if col >= len(line) && len(line) > 0 {
		col = len(line) - 1
	}

	openCol := -1
	for c := col; c >= 0; c-- {
		if line[c] == open {
			openCol = c
			break
		}
		if line[c] == close {
			break
		}
	}
	if openCol < 0 {
		return 0, 0, false
	}

	depth := 1
	closeCol := -1
	for c := openCol + 1; c < len(line); c++ {
		switch line[c] {
		case open:
			depth++
		case close:
			depth--
			if depth == 0 {
				closeCol = c
				break
			}
		}
	}
	if closeCol < 0 {
		return 0, 0, false
	}

	if inner {
		from = t.OffsetOf(text.Point{Line: point.Line, Col: openCol + 1})
		to = t.OffsetOf(text.Point{Line: point.Line, Col: closeCol})
	} else {
		from = t.OffsetOf(text.Point{Line: point.Line, Col: openCol})
		to = t.OffsetOf(text.Point{Line: point.Line, Col: closeCol + 1})
	}
	if to <= from {
		return 0, 0, false
	}
	return from, to, true
}

func tagObjectRange(t *text.Text, at text.Offset, inner bool) (from, to text.Offset, ok bool) {
	point := t.PointOf(at)
	line := point.Line
	openLine := line
	openCol := -1

	for l := line; l >= 0; l-- {
		lineBytes := t.Line(l)
		start := len(lineBytes) - 1
		if l == line {
			start = point.Col
		}
		for c := start; c >= 0; c-- {
			if c+1 < len(lineBytes) && lineBytes[c] == '<' && lineBytes[c+1] != '/' {
				openLine = l
				openCol = c
				goto foundOpen
			}
		}
	}
	return 0, 0, false

foundOpen:
	openBytes := t.Line(openLine)
	endTag := openCol
	for endTag < len(openBytes) && openBytes[endTag] != '>' {
		endTag++
	}
	if endTag >= len(openBytes) {
		return 0, 0, false
	}
	tagName := openBytes[openCol+1 : endTag]
	if len(tagName) == 0 {
		return 0, 0, false
	}

	closeNeedle := []byte("</" + string(tagName) + ">")
	searchFrom := t.OffsetOf(text.Point{Line: openLine, Col: endTag + 1})
	closeAt := indexOfBytes(t, searchFrom, closeNeedle)
	if closeAt < 0 {
		return 0, 0, false
	}

	from = t.OffsetOf(text.Point{Line: openLine, Col: openCol})
	to = closeAt + text.Offset(len(closeNeedle))
	if inner {
		from = t.OffsetOf(text.Point{Line: openLine, Col: endTag + 1})
		to = closeAt
	}
	if to <= from {
		return 0, 0, false
	}
	return from, to, true
}

func indexOfBytes(t *text.Text, from text.Offset, needle []byte) text.Offset {
	if len(needle) == 0 {
		return -1
	}
	hay := t.Slice(from, t.Len())
	idx := bytesIndex(hay, needle)
	if idx < 0 {
		return -1
	}
	return from + text.Offset(idx)
}

func bytesIndex(hay, needle []byte) int {
	for i := 0; i+len(needle) <= len(hay); i++ {
		match := true
		for j := range needle {
			if hay[i+j] != needle[j] {
				match = false
				break
			}
		}
		if match {
			return i
		}
	}
	return -1
}

func sentenceObjectRange(t *text.Text, at text.Offset, around bool) (from, to text.Offset, ok bool) {
	start, okStart := previousSentenceStart(t, at)
	if !okStart {
		start = t.LineStart(t.PointOf(at).Line)
	}
	end, okEnd := nextSentenceStart(t, at)
	if !okEnd {
		end = t.LineEnd(t.PointOf(at).Line)
	}
	if end <= start {
		return 0, 0, false
	}
	from, to = start, end
	if around {
		for to < t.Len() {
			b := t.Slice(to, to+1)[0]
			if b != ' ' && b != '\t' && b != '\n' {
				break
			}
			to++
		}
	}
	if to <= from {
		return 0, 0, false
	}
	return from, to, true
}
