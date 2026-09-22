package editor

import (
	"vague/bonefire/buffer"
	frame "vague/bonefire/frame"
	"vague/bonefire/text"
	"vague/bonefire/window"
)

const maxJumpList = 100

func (e *Editor) pushJump(from text.Offset) {
	if e.jumpNav {
		return
	}
	if len(e.jumps) > 0 && e.jumps[e.jumpPos] == from {
		return
	}

	if len(e.jumps) == 0 {
		e.jumps = append(e.jumps, from)
		e.jumpPos = 0
		return
	}

	e.jumps = e.jumps[:e.jumpPos+1]
	e.jumps = append(e.jumps, from)
	e.jumpPos = len(e.jumps) - 1

	if len(e.jumps) > maxJumpList {
		e.jumps = e.jumps[len(e.jumps)-maxJumpList:]
		e.jumpPos = len(e.jumps) - 1
	}
}

func (e *Editor) jumpOlder(frame *frame.Frame, win *window.Window, buf *buffer.Buffer) error {
	if e.jumpPos == 0 {
		return nil
	}

	e.jumpNav = true
	e.jumpPos--
	setWindowCursor(buf, win, e.jumps[e.jumpPos])
	view := layoutViewForWindow(buf.Text, win, frame)
	rememberColumn(buf, win, view)
	frame.Dirty = true
	return nil
}

func (e *Editor) jumpNewer(frame *frame.Frame, win *window.Window, buf *buffer.Buffer) error {
	if e.jumpPos+1 >= len(e.jumps) {
		return nil
	}

	e.jumpNav = true
	e.jumpPos++
	setWindowCursor(buf, win, e.jumps[e.jumpPos])
	view := layoutViewForWindow(buf.Text, win, frame)
	rememberColumn(buf, win, view)
	frame.Dirty = true
	return nil
}

func (e *Editor) finishMotionJump(buf *buffer.Buffer, win *window.Window, from text.Offset) {
	if e.jumpNav {
		e.jumpNav = false
		return
	}
	if windowCursor(win) != from {
		e.noteJumpFrom(buf, from)
		e.pushJump(from)
	}
}
