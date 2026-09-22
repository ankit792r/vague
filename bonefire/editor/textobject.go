package editor

import (
	"vague/bonefire/buffer"
	frame "vague/bonefire/frame"
	"vague/bonefire/text"
	"vague/bonefire/window"
)

func textObjectRange(t *text.Text, win *window.Window, keys string) (from, to text.Offset, ok bool) {
	at := windowCursor(win)
	switch keys {
	case "iw":
		from = wordStartAt(t, at)
		to = wordEndExclusive(t, from)
	case "aw":
		from = wordStartAt(t, at)
		to = wordEndExclusive(t, from)
		for to < t.Len() {
			cls, _, n := charAt(t, to)
			if cls != wcWhitespace {
				break
			}
			to = n
		}
	case "iW":
		from, to = bigWordInner(t, at)
	case "aW":
		from, to = bigWordAround(t, at)
	case "ip":
		line := t.PointOf(at).Line
		startLine := paragraphStartLine(t, line)
		endLine := paragraphEndLine(t, line)
		from = t.LineStart(startLine)
		to = t.LineEnd(endLine)
	case "ap":
		from, to, _ = textObjectRange(t, win, "ip")
		if to < t.Len() {
			to++
		}
	default:
		return 0, 0, false
	}
	if to <= from {
		return 0, 0, false
	}
	return from, to, true
}

func bigWordInner(t *text.Text, at text.Offset) (text.Offset, text.Offset) {
	point := t.PointOf(at)
	line := t.Line(point.Line)
	if len(line) == 0 {
		return t.LineStart(point.Line), t.LineStart(point.Line)
	}

	startCol := point.Col
	for startCol > 0 && line[startCol-1] != ' ' && line[startCol-1] != '\t' {
		startCol--
	}
	endCol := point.Col
	for endCol < len(line) && line[endCol] != ' ' && line[endCol] != '\t' {
		endCol++
	}
	return t.OffsetOf(text.Point{Line: point.Line, Col: startCol}),
		t.OffsetOf(text.Point{Line: point.Line, Col: endCol})
}

func bigWordAround(t *text.Text, at text.Offset) (text.Offset, text.Offset) {
	from, to := bigWordInner(t, at)
	line := t.PointOf(at).Line
	lineBytes := t.Line(line)
	endCol := t.PointOf(to).Col
	for endCol < len(lineBytes) && (lineBytes[endCol] == ' ' || lineBytes[endCol] == '\t') {
		endCol++
	}
	return from, t.OffsetOf(text.Point{Line: line, Col: endCol})
}

func (e *Editor) applyOperatorTextObject(
	frame *frame.Frame,
	win *window.Window,
	buf *buffer.Buffer,
	op opKind,
	keys string,
	count int,
) error {
	if count < 1 {
		count = 1
	}
	for i := 0; i < count; i++ {
		from, to, ok := textObjectRange(buf.Text, win, keys)
		if !ok {
			break
		}
		landing := clampToLine(buf.Text, from)
		if err := e.applyOperatorRange(frame, win, buf, op, from, to, false, landing); err != nil {
			return err
		}
		if op == opChange {
			break
		}
	}
	return nil
}
