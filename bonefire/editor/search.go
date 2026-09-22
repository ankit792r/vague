package editor

import (
	"bytes"
	"errors"

	"vague/bonefire/buffer"
	frame "vague/bonefire/frame"
	"vague/bonefire/text"
	"vague/bonefire/window"
)

var ErrPatternNotFound = errors.New("Pattern not found")

func (e *Editor) Search(
	frame *frame.Frame,
	win *window.Window,
	buf *buffer.Buffer,
	pattern string,
	forward bool,
) error {
	if pattern == "" {
		return ErrPatternNotFound
	}

	e.endIncsearchCommit()
	e.clearSearchContext()
	e.searchPattern = pattern
	e.searchForward = forward

	return e.runSearch(frame, win, buf, forward)
}

func (e *Editor) RepeatSearch(
	frame *frame.Frame,
	win *window.Window,
	buf *buffer.Buffer,
	sameDirection bool,
) error {
	if e.searchPattern == "" {
		return nil
	}

	forward := e.searchForward
	if !sameDirection {
		forward = !forward
	}

	return e.runSearch(frame, win, buf, forward)
}

func (e *Editor) clearSearchContext() {
	if e.Mode == VisualMode || e.Mode == VisualLineMode || e.Mode == VisualBlockMode {
		e.leaveVisual(nil)
	}
	e.clearPendingOp()
	e.pendingKey = ""
}

func (e *Editor) runSearch(
	frame *frame.Frame,
	win *window.Window,
	buf *buffer.Buffer,
	forward bool,
) error {
	pat := []byte(e.searchPattern)
	if len(pat) == 0 {
		return ErrPatternNotFound
	}

	data := buf.Text.Bytes()
	at := int(windowCursor(win))

	var (
		match int
		ok    bool
	)
	if forward {
		match, ok = findForward(data, pat, at)
	} else {
		match, ok = findBackward(data, pat, at)
	}

	if !ok {
		return ErrPatternNotFound
	}

	off := text.Offset(match)
	e.setSearchMatch(buf, off, off+text.Offset(len(pat)))
	setWindowCursor(buf, win, off)
	view := layoutViewForWindow(buf.Text, win, frame)
	rememberColumn(buf, win, view)
	frame.Dirty = true
	return nil
}

func findForward(data, pat []byte, at int) (int, bool) {
	if len(pat) == 0 {
		return 0, false
	}
	if at < 0 {
		at = 0
	}
	if at > len(data) {
		at = len(data)
	}

	for i := 0; i+len(pat) <= len(data); i++ {
		if bytes.Equal(data[i:i+len(pat)], pat) && i > at {
			return i, true
		}
	}

	for i := 0; i+len(pat) <= len(data); i++ {
		if bytes.Equal(data[i:i+len(pat)], pat) {
			return i, true
		}
	}

	return 0, false
}

func findBackward(data, pat []byte, at int) (int, bool) {
	if len(pat) == 0 {
		return 0, false
	}
	if at < 0 {
		at = 0
	}
	if at > len(data) {
		at = len(data)
	}

	best := -1
	for i := 0; i+len(pat) <= len(data); i++ {
		if bytes.Equal(data[i:i+len(pat)], pat) && i < at {
			best = i
		}
	}
	if best >= 0 {
		return best, true
	}

	for i := len(data) - len(pat); i >= 0; i-- {
		if bytes.Equal(data[i:i+len(pat)], pat) {
			return i, true
		}
	}

	return 0, false
}
