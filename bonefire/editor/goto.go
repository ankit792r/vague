package editor

import (
	"errors"

	"vague/bonefire/buffer"
	frame "vague/bonefire/frame"
	"vague/bonefire/window"
)

// ErrInvalidLine is returned when a goto line number is out of range.
var ErrInvalidLine = errors.New("Invalid line number")

// GoToLine moves the cursor to a 1-based buffer line.
func (e *Editor) GoToLine(frame *frame.Frame, win *window.Window, buf *buffer.Buffer, line1 int) error {
	if line1 < 1 {
		return ErrInvalidLine
	}

	lineCount := buf.Text.LineCount()
	if line1 > lineCount {
		return ErrInvalidLine
	}

	line := line1 - 1
	setWindowCursor(buf, win, moveToBufferLine(buf.Text, win, line))
	view := layoutViewForWindow(buf.Text, win, frame)
	rememberColumn(buf, win, view)
	frame.Dirty = true
	return nil
}
