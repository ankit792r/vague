package editor

import (
	"vague/bonefire/buffer"
	"vague/bonefire/display"
	"vague/bonefire/window"
)

func (e *Editor) openLine(frame *display.Frame, win *window.Window, buf *buffer.Buffer, above bool) error {
	if buf.ReadOnly {
		return buffer.ErrReadOnly
	}

	t := buf.Text
	point := t.PointOf(windowCursor(win))

	at := t.LineEnd(point.Line)
	if above {
		at = t.LineStart(point.Line)
	}

	e.pendingKey = ""
	e.beginInsertGroup(buf, windowCursor(win))

	delta, err := buf.Insert(at, []byte("\n"))
	if err != nil {
		return err
	}

	landing := delta.NewEnd
	if above {
		landing = delta.Start
	}

	setWindowCursor(buf, win, landing)
	e.Mode = InsertMode
	frame.Dirty = true
	return nil
}
