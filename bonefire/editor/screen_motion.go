package editor

import (
	"vague/bonefire/buffer"
	frame "vague/bonefire/frame"
	"vague/bonefire/text"
	"vague/bonefire/window"
)

func moveToScreenLine(
	t *text.Text,
	win *window.Window,
	fm *frame.Frame,
	which rune,
) text.Offset {
	height := fm.Height
	if height < 1 {
		height = frame.DefaultHeight
	}

	fullView := layoutView(t, layoutContentWidth(fm.Width, lineCountForGutter(t), win.WindowOptions.Number), 0, win.WindowOptions.Wrap)
	total := len(fullView.Meta)
	if total == 0 {
		return windowCursor(win)
	}

	top := win.TopLine
	targetRow := top
	switch which {
	case 'H':
		targetRow = top
	case 'M':
		targetRow = top + (height-1)/2
	case 'L':
		targetRow = top + height - 1
	}

	if targetRow >= total {
		targetRow = total - 1
	}
	if targetRow < 0 {
		targetRow = 0
	}

	vl := fullView.Meta[targetRow]
	col := win.DesiredCol
	segLen := vl.EndCol - vl.StartCol
	if segLen <= 0 {
		col = 0
	} else if col >= segLen {
		col = segLen - 1
	}

	return t.OffsetOf(text.Point{Line: vl.BufferLine, Col: vl.StartCol + col})
}

func scrollCursorLine(win *window.Window, fm *frame.Frame, buf *buffer.Buffer, second rune) {
	height := fm.Height
	if height < 1 {
		height = frame.DefaultHeight
	}

	t := buf.Text
	fullView := layoutView(t, layoutContentWidth(fm.Width, lineCountForGutter(t), win.WindowOptions.Number), 0, win.WindowOptions.Wrap)
	point := windowPoint(buf, win)
	cursorRow, _, ok := visualRowAt(fullView.Meta, point)
	if !ok {
		return
	}

	switch second {
	case 't':
		win.TopLine = cursorRow
	case 'z':
		win.TopLine = cursorRow - (height-1)/2
	case 'b':
		win.TopLine = cursorRow - height + 1
	default:
		return
	}

	maxTop := len(fullView.Meta) - height
	if maxTop < 0 {
		maxTop = 0
	}
	if win.TopLine < 0 {
		win.TopLine = 0
	}
	if win.TopLine > maxTop {
		win.TopLine = maxTop
	}
}
