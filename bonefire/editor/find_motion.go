package editor

import (
	"unicode/utf8"

	"vague/bonefire/buffer"
	frame "vague/bonefire/frame"
	"vague/bonefire/text"
	"vague/bonefire/window"
)

type charFindKind int

const (
	charFindNone charFindKind = iota
	charFindF
	charFindBigF
	charFindT
	charFindBigT
)

func oppositeCharFind(kind charFindKind) charFindKind {
	switch kind {
	case charFindF:
		return charFindBigF
	case charFindBigF:
		return charFindF
	case charFindT:
		return charFindBigT
	case charFindBigT:
		return charFindT
	default:
		return charFindNone
	}
}

func runeFromKey(keys string) (rune, bool) {
	if keys == "" || keys[0] == '<' {
		return 0, false
	}

	r, _ := utf8.DecodeRuneInString(keys)
	if r == utf8.RuneError {
		return 0, false
	}

	return r, true
}

// findCharOnLine moves within the current line. count is how many matches to step.
func findCharOnLine(t *text.Text, at text.Offset, kind charFindKind, target rune, count int) (text.Offset, bool) {
	if count < 1 {
		count = 1
	}

	point := t.PointOf(at)
	line := t.Line(point.Line)
	if len(line) == 0 {
		return at, false
	}

	switch kind {
	case charFindF, charFindT:
		return findForwardOnLine(t, point.Line, line, point.Col, kind == charFindT, target, count)
	case charFindBigF, charFindBigT:
		return findBackwardOnLine(t, point.Line, line, point.Col, kind == charFindBigT, target, count)
	default:
		return at, false
	}
}

func findForwardOnLine(
	t *text.Text,
	bufLine int,
	line []byte,
	fromCol int,
	till bool,
	target rune,
	count int,
) (text.Offset, bool) {
	start := fromCol
	if count == 1 {
		start = fromCol + 1
	}
	if start < 0 {
		start = 0
	}

	matches := 0
	for col := start; col < len(line); {
		r, size := utf8.DecodeRune(line[col:])
		if r == target {
			matches++
			if matches == count {
				if till {
					if col == 0 {
						return 0, false
					}
					return t.OffsetOf(text.Point{Line: bufLine, Col: col - 1}), true
				}
				return t.OffsetOf(text.Point{Line: bufLine, Col: col}), true
			}
		}
		col += size
	}

	return 0, false
}

func findBackwardOnLine(
	t *text.Text,
	bufLine int,
	line []byte,
	fromCol int,
	till bool,
	target rune,
	count int,
) (text.Offset, bool) {
	searchEnd := fromCol - 1
	if searchEnd >= len(line) {
		searchEnd = len(line) - 1
	}
	if searchEnd < 0 {
		return 0, false
	}

	matches := 0
	for col := searchEnd; col >= 0; {
		for col > 0 && !utf8.RuneStart(line[col]) {
			col--
		}

		r, size := utf8.DecodeRune(line[col:])
		if r == target {
			matches++
			if matches == count {
				if till {
					after := col + size
					if after > len(line) {
						return 0, false
					}
					return t.OffsetOf(text.Point{Line: bufLine, Col: after}), true
				}
				return t.OffsetOf(text.Point{Line: bufLine, Col: col}), true
			}
		}

		if col == 0 {
			break
		}
		col--
	}

	return 0, false
}

func (e *Editor) clearPendingCharFind() {
	e.pendingCharFind = charFindNone
}

func (e *Editor) executeCharFind(
	frame *frame.Frame,
	win *window.Window,
	buf *buffer.Buffer,
	kind charFindKind,
	target rune,
	count int,
) error {
	if kind == charFindNone {
		return nil
	}

	t := buf.Text
	at := windowCursor(win)
	next, ok := findCharOnLine(t, at, kind, target, count)
	if !ok {
		frame.Dirty = true
		return nil
	}

	e.lastCharFindKind = kind
	e.lastCharFindRune = target

	setWindowCursor(buf, win, next)
	view := layoutViewForWindow(t, win, frame)
	rememberColumn(buf, win, view)
	frame.Dirty = true
	return nil
}

func (e *Editor) beginPendingCharFind(kind charFindKind) {
	e.clearPendingOp()
	e.pendingCharFind = kind
}

func (e *Editor) consumePendingCharFind(
	frame *frame.Frame,
	win *window.Window,
	buf *buffer.Buffer,
	keys string,
) (bool, error) {
	if e.pendingCharFind == charFindNone {
		return false, nil
	}

	kind := e.pendingCharFind
	e.clearPendingCharFind()

	target, ok := runeFromKey(keys)
	if !ok {
		e.clearPendingCount()
		return true, nil
	}

	count := e.takeCount()
	return true, e.executeCharFind(frame, win, buf, kind, target, count)
}
