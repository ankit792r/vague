package editor

import (
	"vague/bonefire/buffer"
	frame "vague/bonefire/frame"
	"vague/bonefire/text"
	"vague/bonefire/window"
)

// lineChangeRange returns the byte span deleted by dd on a single line.
func lineChangeRange(t *text.Text, line int) (from, to text.Offset) {
	from = t.LineStart(line)
	to = t.LineEnd(line)

	switch {
	case line+1 < t.LineCount():
		to = t.LineStart(line + 1)
	case line > 0:
		from = t.LineEnd(line - 1)
	}

	return from, to
}

func (e *Editor) yankLine(frame *frame.Frame, win *window.Window, buf *buffer.Buffer) error {
	t := buf.Text
	line := t.PointOf(windowCursor(win)).Line
	from, to := lineChangeRange(t, line)
	if to <= from {
		return nil
	}

	e.setRegisterLinewise(t.Slice(from, to))
	frame.Dirty = true
	return nil
}

func (e *Editor) pasteAfter(frame *frame.Frame, win *window.Window, buf *buffer.Buffer) error {
	text, kind, ok := e.registerText()
	if !ok {
		return nil
	}

	if kind == yankLinewise {
		return e.pasteLinewise(frame, win, buf, text, false)
	}

	return e.pasteCharwise(frame, win, buf, text, false)
}

func (e *Editor) pasteBefore(frame *frame.Frame, win *window.Window, buf *buffer.Buffer) error {
	text, kind, ok := e.registerText()
	if !ok {
		return nil
	}

	if kind == yankLinewise {
		return e.pasteLinewise(frame, win, buf, text, true)
	}

	return e.pasteCharwise(frame, win, buf, text, true)
}

func (e *Editor) pasteLinewise(
	frame *frame.Frame,
	win *window.Window,
	buf *buffer.Buffer,
	data []byte,
	before bool,
) error {
	if buf.ReadOnly {
		return buffer.ErrReadOnly
	}

	t := buf.Text
	at := windowCursor(win)
	line := t.PointOf(at).Line

	var insertAt text.Offset
	if before {
		insertAt = t.LineStart(line)
	} else if line+1 < t.LineCount() {
		insertAt = t.LineStart(line + 1)
	} else {
		insertAt = t.Len()
	}

	buf.BeginEdit(at)
	delta, err := buf.Insert(insertAt, data)
	if err != nil {
		return err
	}

	setWindowCursor(buf, win, delta.Start)
	buf.EndEdit(windowCursor(win))

	view := layoutViewForWindow(buf.Text, win, frame)
	rememberColumn(buf, win, view)
	frame.Dirty = true
	return nil
}

func (e *Editor) pasteCharwise(
	frame *frame.Frame,
	win *window.Window,
	buf *buffer.Buffer,
	data []byte,
	before bool,
) error {
	if buf.ReadOnly {
		return buffer.ErrReadOnly
	}

	t := buf.Text
	at := windowCursor(win)
	if before {
		if at > 0 {
			at = back(t, at)
		}
	} else {
		next := forwardChar(t, at)
		if next == at {
			next = t.LineEnd(t.PointOf(at).Line)
		}
		at = next
	}

	buf.BeginEdit(windowCursor(win))
	delta, err := buf.Insert(at, data)
	if err != nil {
		return err
	}

	setWindowCursor(buf, win, delta.NewEnd)
	buf.EndEdit(windowCursor(win))

	view := layoutViewForWindow(buf.Text, win, frame)
	rememberColumn(buf, win, view)
	frame.Dirty = true
	return nil
}
