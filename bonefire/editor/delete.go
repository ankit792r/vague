package editor

import (
	"unicode/utf8"

	"vague/bonefire/buffer"
	frame "vague/bonefire/frame"
	"vague/bonefire/text"
	"vague/bonefire/window"
)

func (e *Editor) deleteCharRepeat(
	frame *frame.Frame,
	win *window.Window,
	buf *buffer.Buffer,
	count int,
) error {
	if count < 1 {
		count = 1
	}

	for i := 0; i < count; i++ {
		t := buf.Text
		at := windowCursor(win)
		end := forwardChar(t, at)
		if end == at {
			break
		}

		if err := e.deleteChar(frame, win, buf); err != nil {
			return err
		}
	}

	return nil
}

func (e *Editor) deleteChar(frame *frame.Frame, win *window.Window, buf *buffer.Buffer) error {
	if buf.ReadOnly {
		return buffer.ErrReadOnly
	}

	t := buf.Text
	at := windowCursor(win)
	end := forwardChar(t, at)
	if end == at {
		return nil
	}

	data := t.Slice(at, end)
	e.recordDelete(data, false)

	e.noteChangeAt(buf, at)
	buf.BeginEdit(at)
	if _, err := buf.Delete(at, end); err != nil {
		return err
	}
	setWindowCursor(buf, win, clampToLine(t, at))
	buf.EndEdit(windowCursor(win))

	view := layoutViewForWindow(t, win, frame)
	rememberColumn(buf, win, view)
	frame.Dirty = true
	e.recordDeleteCharChange()
	return nil
}

func (e *Editor) deleteLine(frame *frame.Frame, win *window.Window, buf *buffer.Buffer) error {
	if buf.ReadOnly {
		return buffer.ErrReadOnly
	}

	t := buf.Text
	at := windowCursor(win)
	first := t.PointOf(at).Line

	from, to := lineChangeRange(t, first)

	if to <= from {
		return nil
	}

	landingLine := first
	if landingLine >= t.LineCount() {
		landingLine = t.LineCount() - 1
	}
	landing := t.LineStart(landingLine)

	buf.BeginEdit(at)
	if _, err := buf.Delete(from, to); err != nil {
		return err
	}
	setWindowCursor(buf, win, landing)
	buf.EndEdit(windowCursor(win))

	view := layoutViewForWindow(t, win, frame)
	rememberColumn(buf, win, view)
	frame.Dirty = true
	return nil
}

func forwardChar(t *text.Text, off text.Offset) text.Offset {
	point := t.PointOf(off)
	line := t.Line(point.Line)
	if point.Col >= len(line) {
		return off
	}

	_, size := utf8.DecodeRune(line[point.Col:])
	return off + text.Offset(size)
}
