package editor

import (
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
	pat, off := splitSearchOffset(pattern)
	e.searchPattern = pat
	e.searchOffset = off
	e.searchForward = forward
	e.nohlSearch = false

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
	pp, err := e.parseSearchPattern(e.searchPattern)
	if err != nil {
		return err
	}

	data := buf.Text.Bytes()
	at := int(windowCursor(win))
	wrap := e.searchOpts.WrapScan

	var (
		start, end int
		ok         bool
	)
	if forward {
		start, end, ok = findPatternForward(data, pp, at, wrap)
	} else {
		start, end, ok = findPatternBackward(data, pp, at, wrap)
	}

	if !ok {
		return ErrPatternNotFound
	}

	off := text.Offset(start)
	endOff := text.Offset(end)
	e.setSearchMatch(buf, off, endOff)
	cursor := applySearchOffset(start, end, len(data), e.searchOffset)
	setWindowCursor(buf, win, text.Offset(cursor))
	view := layoutViewForWindow(buf.Text, win, frame)
	rememberColumn(buf, win, view)
	frame.Dirty = true
	return nil
}

func (e *Editor) ClearNohlSearch(frame *frame.Frame) {
	e.nohlSearch = true
	frame.Dirty = true
}

func (e *Editor) compiledPattern() (parsedPattern, error) {
	if e.searchPattern == "" {
		return parsedPattern{}, ErrPatternNotFound
	}
	return e.parseSearchPattern(e.searchPattern)
}
