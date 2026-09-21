package editor

import (
	"unicode"
	"unicode/utf8"

	"vague/bonefire/buffer"
	frame "vague/bonefire/frame"
	"vague/bonefire/text"
	"vague/bonefire/window"
)

type opKind int

const (
	opNone opKind = iota
	opDelete
	opYank
	opChange
)

type motionKind int

const (
	motionLine motionKind = iota
	motionWord
	motionToEOL
)

func (e *Editor) clearPendingOp() {
	e.pendingOp = opNone
}

func motionForKey(keys string) (motionKind, bool) {
	switch keys {
	case "w":
		return motionWord, true
	case "$":
		return motionToEOL, true
	default:
		return 0, false
	}
}

func textRangeForMotion(t *text.Text, win *window.Window, op opKind, motion motionKind) (from, to text.Offset, linewise bool) {
	at := windowCursor(win)
	line := t.PointOf(at).Line

	switch motion {
	case motionLine:
		if op == opChange {
			from = t.LineStart(line)
			to = t.LineEnd(line)
			if to <= from && line+1 < t.LineCount() {
				to = t.LineStart(line + 1)
			}
			return from, to, true
		}
		from, to = lineChangeRange(t, line)
		return from, to, true
	case motionWord:
		if op == opChange {
			end := wordEndExclusive(t, at)
			if end <= at {
				return at, at, false
			}
			return at, end, false
		}
		end := moveWordForwardOnce(t, at)
		if end <= at {
			return at, at, false
		}
		return at, end, false
	case motionToEOL:
		end := t.LineEnd(line)
		if at >= end {
			return at, at, false
		}
		return at, end, false
	default:
		return at, at, false
	}
}

func landingAfterDelete(t *text.Text, from, to text.Offset, motion motionKind) text.Offset {
	switch motion {
	case motionLine:
		line := t.PointOf(from).Line
		if line >= t.LineCount() {
			line = t.LineCount() - 1
		}
		return t.LineStart(line)
	default:
		return clampToLine(t, from)
	}
}

func (e *Editor) applyOperatorMotion(
	frame *frame.Frame,
	win *window.Window,
	buf *buffer.Buffer,
	op opKind,
	motion motionKind,
) error {
	t := buf.Text
	from, to, linewise := textRangeForMotion(t, win, op, motion)
	if to <= from {
		frame.Dirty = true
		return nil
	}

	data := t.Slice(from, to)

	switch op {
	case opYank:
		e.stashRegister(data, linewise)
		frame.Dirty = true
		return nil
	case opDelete, opChange:
		if buf.ReadOnly {
			return buffer.ErrReadOnly
		}

		at := windowCursor(win)
		e.stashRegister(data, linewise)

		buf.BeginEdit(at)
		if _, err := buf.Delete(from, to); err != nil {
			return err
		}

		landing := landingAfterDelete(t, from, to, motion)
		setWindowCursor(buf, win, landing)
		buf.EndEdit(windowCursor(win))

		view := layoutViewForWindow(buf.Text, win, frame)
		rememberColumn(buf, win, view)
		frame.Dirty = true

		if op == opChange {
			return e.enterInsert(frame, win, buf, landing)
		}
		return nil
	default:
		return nil
	}
}

func (e *Editor) joinLines(frame *frame.Frame, win *window.Window, buf *buffer.Buffer) error {
	if buf.ReadOnly {
		return buffer.ErrReadOnly
	}

	t := buf.Text
	at := windowCursor(win)
	line := t.PointOf(at).Line
	if line+1 >= t.LineCount() {
		return nil
	}

	end := t.LineEnd(line)
	nextStart := t.LineStart(line + 1)
	if nextStart <= end {
		return nil
	}

	needSpace := false
	lineBytes := t.Line(line)
	nextBytes := t.Line(line + 1)
	if len(lineBytes) > 0 && len(nextBytes) > 0 {
		last, _ := utf8.DecodeLastRune(lineBytes)
		first, _ := utf8.DecodeRune(nextBytes)
		if !unicode.IsSpace(last) && !unicode.IsSpace(first) {
			needSpace = true
		}
	}

	buf.BeginEdit(at)
	if needSpace {
		if _, err := buf.Replace(end, nextStart, []byte(" ")); err != nil {
			return err
		}
	} else if _, err := buf.Delete(end, nextStart); err != nil {
		return err
	}

	setWindowCursor(buf, win, clampToLine(t, end))
	buf.EndEdit(windowCursor(win))

	view := layoutViewForWindow(buf.Text, win, frame)
	rememberColumn(buf, win, view)
	frame.Dirty = true
	return nil
}
