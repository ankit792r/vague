package editor

import (
	"vague/bonefire/buffer"
	"vague/bonefire/text"
	"vague/bonefire/window"
)

func (e *Editor) enterInsert(frame *Frame, win *window.Window, buf *buffer.Buffer, at text.Offset) error {
	if buf.ReadOnly {
		return buffer.ErrReadOnly
	}

	e.pendingKey = ""
	setWindowCursor(buf, win, at)
	e.Mode = InsertMode
	e.beginInsertGroup(buf, at)
	frame.dirty = true
	return nil
}

func (e *Editor) beginInsertGroup(buf *buffer.Buffer, at text.Offset) {
	if e.insertGroup != 0 {
		return
	}

	buf.BeginEdit(at)
	e.insertGroup = buf.ID
}

func (e *Editor) leaveInsert(win *window.Window, buf *buffer.Buffer) {
	e.Mode = NormalMode

	if e.insertGroup == 0 {
		return
	}

	if e.insertGroup == buf.ID {
		buf.EndEdit(windowCursor(win))
	}

	e.insertGroup = 0
}

func (e *Editor) undoTo(frame *Frame, win *window.Window, buf *buffer.Buffer, redo bool) error {
	var (
		at text.Offset
		ok bool
	)

	if redo {
		at, ok = buf.Redo()
	} else {
		at, ok = buf.Undo()
	}

	if !ok {
		return nil
	}

	setWindowCursor(buf, win, clampToLine(buf.Text, at))
	view := layoutViewForWindow(buf.Text, win, frame)
	rememberColumn(buf, win, view)
	frame.dirty = true
	return nil
}

func clampToLine(t *text.Text, off text.Offset) text.Offset {
	off = t.Clamp(off)
	point := t.PointOf(off)
	line := t.Line(point.Line)
	col := point.Col
	if col > lastColumn(line) {
		col = lastColumn(line)
	}
	return t.OffsetOf(text.Point{Line: point.Line, Col: col})
}
