package editor

import (
	"vague/bonefire/buffer"
	frame "vague/bonefire/frame"
	"vague/bonefire/text"
	"vague/bonefire/window"
)

// WindowCursorForTest exposes cursor position for external tests.
func WindowCursorForTest(win *window.Window) text.Offset {
	return windowCursor(win)
}

// SetWindowCursorForTest moves the window cursor in tests.
func SetWindowCursorForTest(buf *buffer.Buffer, win *window.Window, at text.Offset) {
	setWindowCursor(buf, win, at)
}

// LayoutViewForTest exposes layout for external tests.
func LayoutViewForTest(t *text.Text, width, maxRows int, wrap bool) viewLayout {
	return layoutView(t, width, maxRows, wrap)
}

// EnsureCursorVisibleForTest runs scroll logic in tests.
func EnsureCursorVisibleForTest(win *window.Window, frame *frame.Frame, view viewLayout, point text.Point) {
	ensureCursorVisible(win, frame, view, point)
}
