package editor

import (
	"vague/bonefire/buffer"
	frame "vague/bonefire/frame"
	"vague/bonefire/text"
	"vague/bonefire/window"
)

func (e *Editor) BeginIncsearch(win *window.Window, buf *buffer.Buffer, forward bool) {
	e.incsearchActive = true
	e.incsearchForward = forward
	e.incsearchReturnCursor = windowCursor(win)
	e.incsearchSavedPattern = e.searchPattern
	e.incsearchSavedMatchFrom = e.searchMatchFrom
	e.incsearchSavedMatchTo = e.searchMatchTo
	e.incsearchSavedMatchBuf = e.searchMatchBuf
}

func (e *Editor) CancelIncsearch(
	frame *frame.Frame,
	win *window.Window,
	buf *buffer.Buffer,
) {
	if !e.incsearchActive {
		e.ClearSearchMatch()
		frame.Dirty = true
		return
	}

	setWindowCursor(buf, win, e.incsearchReturnCursor)
	e.searchPattern = e.incsearchSavedPattern
	e.searchMatchFrom = e.incsearchSavedMatchFrom
	e.searchMatchTo = e.incsearchSavedMatchTo
	e.searchMatchBuf = e.incsearchSavedMatchBuf
	e.incsearchActive = false
	frame.Dirty = true
}

func (e *Editor) endIncsearchCommit() {
	e.incsearchActive = false
}

func (e *Editor) PreviewIncsearch(
	frame *frame.Frame,
	win *window.Window,
	buf *buffer.Buffer,
	pattern string,
) error {
	if !e.incsearchActive {
		return nil
	}

	if pattern == "" {
		setWindowCursor(buf, win, e.incsearchReturnCursor)
		e.ClearSearchMatch()
		frame.Dirty = true
		return nil
	}

	data := buf.Text.Bytes()
	at := int(e.incsearchReturnCursor)
	var (
		match int
		ok    bool
	)
	if e.incsearchForward {
		match, ok = findForward(data, []byte(pattern), at)
	} else {
		match, ok = findBackward(data, []byte(pattern), at)
	}

	if !ok {
		e.ClearSearchMatch()
		setWindowCursor(buf, win, e.incsearchReturnCursor)
		frame.Dirty = true
		return nil
	}

	off := text.Offset(match)
	e.setSearchMatch(buf, off, off+text.Offset(len(pattern)))
	setWindowCursor(buf, win, off)
	view := layoutViewForWindow(buf.Text, win, frame)
	rememberColumn(buf, win, view)
	frame.Dirty = true
	return nil
}
