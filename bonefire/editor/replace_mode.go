package editor

import (
	"vague/bonefire/buffer"
	frame "vague/bonefire/frame"
	"vague/bonefire/text"
	"vague/bonefire/window"
)

// ReplaceMode overwrites characters on insert (Vim R).
const ReplaceMode Mode = 5

func (e *Editor) enterReplace(frame *frame.Frame, win *window.Window, buf *buffer.Buffer, at text.Offset) error {
	e.pendingKey = ""
	e.clearPendingOp()
	e.beginInsertGroup(buf, at)
	setWindowCursor(buf, win, at)
	e.Mode = ReplaceMode
	frame.Dirty = true
	return nil
}

func (e *Editor) replaceKey(frame *frame.Frame, win *window.Window, buf *buffer.Buffer, keys string) error {
	if buf.ReadOnly {
		return buffer.ErrReadOnly
	}

	switch keys {
	case "<Esc>":
		e.leaveInsert(win, buf)
		frame.Dirty = true
		return nil
	case "<CR>":
		return e.insertBytes(frame, win, buf, []byte("\n"))
	case "<BS>":
		return e.deleteBack(frame, win, buf)
	}

	if keys == "" || keys[0] == '<' {
		return nil
	}

	t := buf.Text
	at := windowCursor(win)
	end := moveRight(t, at, 1, true)
	if end > at {
		if _, err := buf.Replace(at, end, []byte(keys[:1])); err != nil {
			return err
		}
		setWindowCursor(buf, win, moveRight(t, at, 1, false))
	} else {
		return e.insertBytes(frame, win, buf, []byte(keys[:1]))
	}
	frame.Dirty = true
	return nil
}
