package editor

import (
	"vague/bonefire/buffer"
	frame "vague/bonefire/frame"
	"vague/bonefire/window"
)

// HandleInput applies one key in Vim notation for the given frame window.
func (e *Editor) HandleInput(frame *frame.Frame, win *window.Window, buf *buffer.Buffer, keys string) error {
	switch e.Mode {
	case InsertMode:
		return e.insertKey(frame, win, buf, keys)
	case ReplaceMode:
		return e.replaceKey(frame, win, buf, keys)
	case VisualMode, VisualLineMode, VisualBlockMode:
		return e.visualKey(frame, win, buf, keys)
	default:
		err := e.normalKey(frame, win, buf, keys)
		if e.insertNormalOnce {
			e.insertNormalOnce = false
			e.Mode = InsertMode
			frame.Dirty = true
		}
		return err
	}
}

func (e *Editor) insertKey(frame *frame.Frame, win *window.Window, buf *buffer.Buffer, keys string) error {
	if buf.ReadOnly {
		return buffer.ErrReadOnly
	}

	t := buf.Text

	switch keys {
	case "<Esc>":
		e.insertCompleteActive = false
		e.leaveInsert(win, buf)
		setWindowCursor(buf, win, moveLeft(t, windowCursor(win), 1))
		view := layoutViewForWindow(t, win, frame)
		rememberColumn(buf, win, view)
		frame.Dirty = true
		return nil
	case "<C-w>":
		return e.insertDeleteWordBack(frame, win, buf)
	case "<C-u>":
		return e.insertDeleteToLineStart(frame, win, buf)
	case "<C-k>":
		e.beginDigraph()
		frame.Dirty = true
		return nil
	case "<C-a>":
		return e.insertMoveLineStart(frame, win, buf)
	case "<C-e>":
		return e.insertMoveLineEnd(frame, win, buf)
	case "<C-o>":
		e.insertNormalOnce = true
		e.Mode = NormalMode
		frame.Dirty = true
		return nil
	case "<C-r>":
		e.pendingInsertReg = true
		frame.Dirty = true
		return nil
	case "<C-t>":
		return e.insertAdjustIndent(frame, win, buf, 1)
	case "<C-d>":
		return e.insertAdjustIndent(frame, win, buf, -1)
	case "<C-x>":
		e.beginInsertCompletion(win, buf)
		frame.Dirty = true
		return nil
	case "<C-n>":
		if e.insertCompleteActive {
			return e.insertCompletionCycle(frame, win, buf, 1)
		}
		return nil
	case "<C-p>":
		if e.insertCompleteActive {
			return e.insertCompletionCycle(frame, win, buf, -1)
		}
		return nil
	case "<CR>":
		return e.insertBytes(frame, win, buf, []byte("\n"))
	case "<Space>":
		return e.insertBytes(frame, win, buf, []byte(" "))
	case "<BS>":
		return e.deleteBack(frame, win, buf)
	case "<Left>":
		setWindowCursor(buf, win, moveLeft(t, windowCursor(win), 1))
		frame.Dirty = true
		return nil
	case "<Right>":
		setWindowCursor(buf, win, moveRight(t, windowCursor(win), 1, true))
		frame.Dirty = true
		return nil
	case "<Up>":
		setWindowCursor(buf, win, moveVerticalForWindow(t, win, frame, windowCursor(win), -1, true))
		frame.Dirty = true
		view := layoutViewForWindow(t, win, frame)
		rememberColumn(buf, win, view)
		return nil
	case "<Down>":
		setWindowCursor(buf, win, moveVerticalForWindow(t, win, frame, windowCursor(win), 1, true))
		frame.Dirty = true
		view := layoutViewForWindow(t, win, frame)
		rememberColumn(buf, win, view)
		return nil
	}

	if keys == "" || keys[0] == '<' {
		return nil
	}

	if e.pendingInsertReg {
		e.pendingInsertReg = false
		id, ok := parseRegisterKeyASCII(keys)
		if !ok {
			id = registerID{}
		}
		if data, _, ok := e.readRegister(id); ok {
			return e.insertRegisterPaste(frame, win, buf, data)
		}
		frame.Dirty = true
		return nil
	}

	if e.pendingDigraph {
		if e.digraphFirst == "" {
			if e.setDigraphFirst(keys) {
				frame.Dirty = true
			}
			return nil
		}
		if text, ok := e.consumeDigraphInsert(keys); ok {
			return e.insertBytes(frame, win, buf, []byte(text))
		}
		return nil
	}

	e.insertCompleteActive = false
	return e.insertBytes(frame, win, buf, []byte(keys))
}

func (e *Editor) insertBytes(frame *frame.Frame, win *window.Window, buf *buffer.Buffer, data []byte) error {
	at := windowCursor(win)
	e.noteChangeAt(buf, at)
	delta, err := buf.Insert(at, data)
	if err != nil {
		return err
	}

	setWindowCursor(buf, win, delta.NewEnd)
	view := layoutViewForWindow(buf.Text, win, frame)
	rememberColumn(buf, win, view)
	frame.Dirty = true
	return nil
}

func (e *Editor) deleteBack(frame *frame.Frame, win *window.Window, buf *buffer.Buffer) error {
	at := windowCursor(win)
	if at <= 0 {
		return nil
	}

	from := back(buf.Text, at)
	e.noteChangeAt(buf, from)
	if _, err := buf.Delete(from, at); err != nil {
		return err
	}

	setWindowCursor(buf, win, from)
	view := layoutViewForWindow(buf.Text, win, frame)
	rememberColumn(buf, win, view)
	frame.Dirty = true
	return nil
}
