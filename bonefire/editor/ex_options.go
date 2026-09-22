package editor

import (
	frame "vague/bonefire/frame"
	"vague/bonefire/window"
)

func (e *Editor) setWrap(f *frame.Frame, win *window.Window, wrap bool) error {
	win.WindowOptions.Wrap = wrap
	f.Dirty = true
	return nil
}

func (e *Editor) setNumber(f *frame.Frame, win *window.Window, number bool) error {
	win.WindowOptions.Number = number
	f.Dirty = true
	return nil
}
