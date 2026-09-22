package editor

import (
	"vague/bonefire/buffer"
	frame "vague/bonefire/frame"
	"vague/bonefire/window"
)

func (e *Editor) joinLinesRepeat(
	frame *frame.Frame,
	win *window.Window,
	buf *buffer.Buffer,
	count int,
) error {
	if count < 1 {
		count = 1
	}

	for i := 0; i < count; i++ {
		t := buf.Text
		at := windowCursor(win)
		line := t.PointOf(at).Line
		if line+1 >= t.LineCount() {
			break
		}

		if err := e.joinLines(frame, win, buf); err != nil {
			return err
		}
	}

	return nil
}
