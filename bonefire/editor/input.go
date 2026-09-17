package editor

import (
	"fmt"

	"vague/bonefire/buffer"
	"vague/bonefire/window"
)

// HandleInput applies one key in Vim notation for the given frame.
func (e *Editor) HandleInput(frameID uint64, keys string) error {
	switch e.Mode {
	case InsertMode:
		return e.insertKey(frameID, keys)
	default:
		return e.normalKey(frameID, keys)
	}
}

func (e *Editor) normalKey(frameID uint64, keys string) error {
	frame, win, buf, err := e.frameContext(frameID)
	if err != nil {
		return err
	}

	t := buf.Text
	at := windowCursor(win)

	switch keys {
	case "i":
		e.Mode = InsertMode
		frame.dirty = true
		return nil
	case "a":
		setWindowCursor(buf, win, moveRight(t, at, 1, true))
		e.Mode = InsertMode
		frame.dirty = true
		return nil
	case "h", "<Left>":
		setWindowCursor(buf, win, moveLeft(t, at, 1))
	case "l", "<Right>", "<Space>":
		setWindowCursor(buf, win, moveRight(t, at, 1, false))
	case "j", "<Down>":
		setWindowCursor(buf, win, moveVerticalForWindow(t, win, frame, at, 1, false))
	case "k", "<Up>":
		setWindowCursor(buf, win, moveVerticalForWindow(t, win, frame, at, -1, false))
	default:
		return nil
	}

	view := layoutViewForWindow(t, win, frame)
	rememberColumn(buf, win, view)
	frame.dirty = true
	return nil
}

func (e *Editor) insertKey(frameID uint64, keys string) error {
	frame, win, buf, err := e.frameContext(frameID)
	if err != nil {
		return err
	}

	if buf.ReadOnly {
		return buffer.ErrReadOnly
	}

	t := buf.Text

	switch keys {
	case "<Esc>":
		e.Mode = NormalMode
		setWindowCursor(buf, win, moveLeft(t, windowCursor(win), 1))
		view := layoutViewForWindow(t, win, frame)
		rememberColumn(buf, win, view)
		frame.dirty = true
		return nil
	case "<CR>":
		return e.insertBytes(frame, win, buf, []byte("\n"))
	case "<Space>":
		return e.insertBytes(frame, win, buf, []byte(" "))
	case "<BS>":
		return e.deleteBack(frame, win, buf)
	case "<Left>":
		setWindowCursor(buf, win, moveLeft(t, windowCursor(win), 1))
		frame.dirty = true
		return nil
	case "<Right>":
		setWindowCursor(buf, win, moveRight(t, windowCursor(win), 1, true))
		frame.dirty = true
		return nil
	case "<Up>":
		setWindowCursor(buf, win, moveVerticalForWindow(t, win, frame, windowCursor(win), -1, true))
		frame.dirty = true
		view := layoutViewForWindow(t, win, frame)
		rememberColumn(buf, win, view)
		return nil
	case "<Down>":
		setWindowCursor(buf, win, moveVerticalForWindow(t, win, frame, windowCursor(win), 1, true))
		frame.dirty = true
		view := layoutViewForWindow(t, win, frame)
		rememberColumn(buf, win, view)
		return nil
	}

	if keys == "" || keys[0] == '<' {
		return nil
	}

	return e.insertBytes(frame, win, buf, []byte(keys))
}

func (e *Editor) insertBytes(frame *Frame, win *window.Window, buf *buffer.Buffer, data []byte) error {
	delta := buf.Text.Insert(windowCursor(win), data)
	setWindowCursor(buf, win, delta.NewEnd)
	view := layoutViewForWindow(buf.Text, win, frame)
	rememberColumn(buf, win, view)
	frame.dirty = true
	return nil
}

func (e *Editor) deleteBack(frame *Frame, win *window.Window, buf *buffer.Buffer) error {
	at := windowCursor(win)
	if at <= 0 {
		return nil
	}

	from := back(buf.Text, at)
	buf.Text.Delete(from, at)
	setWindowCursor(buf, win, from)
	view := layoutViewForWindow(buf.Text, win, frame)
	rememberColumn(buf, win, view)
	frame.dirty = true
	return nil
}

func (e *Editor) frameContext(frameID uint64) (*Frame, *window.Window, *buffer.Buffer, error) {
	frame, ok := e.Frames[frameID]
	if !ok {
		return nil, nil, nil, fmt.Errorf("frame %d not found", frameID)
	}

	win, ok := e.Windows[frame.ActiveWindowID]
	if !ok {
		return nil, nil, nil, fmt.Errorf("window %d not found", frame.ActiveWindowID)
	}

	buf, ok := e.Buffers[win.BufferId]
	if !ok {
		return nil, nil, nil, fmt.Errorf("buffer %d not found", win.BufferId)
	}

	return frame, win, buf, nil
}
