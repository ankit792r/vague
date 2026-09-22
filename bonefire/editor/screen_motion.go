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

	fullView := layoutView(t, layoutContentWidth(fm.Width, lineCountForGutter(t), win.WindowOptions), 0, win.WindowOptions.Wrap)
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
	fullView := layoutView(t, layoutContentWidth(fm.Width, lineCountForGutter(t), win.WindowOptions), 0, win.WindowOptions.Wrap)
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

func moveG0(t *text.Text, win *window.Window, fm *frame.Frame, off text.Offset) text.Offset {
	view := layoutView(t, layoutContentWidth(fm.Width, lineCountForGutter(t), win.WindowOptions), 0, win.WindowOptions.Wrap)
	point := t.PointOf(off)
	row, _, ok := visualRowAt(view.Meta, point)
	if !ok {
		return off
	}
	vl := view.Meta[row]
	return t.OffsetOf(text.Point{Line: vl.BufferLine, Col: vl.StartCol})
}

func moveGScreenEnd(t *text.Text, win *window.Window, fm *frame.Frame, off text.Offset) text.Offset {
	view := layoutView(t, layoutContentWidth(fm.Width, lineCountForGutter(t), win.WindowOptions), 0, win.WindowOptions.Wrap)
	point := t.PointOf(off)
	row, _, ok := visualRowAt(view.Meta, point)
	if !ok {
		return off
	}
	vl := view.Meta[row]
	col := vl.EndCol - 1
	if col < vl.StartCol {
		col = vl.StartCol
	}
	return t.OffsetOf(text.Point{Line: vl.BufferLine, Col: col})
}

func moveGm(t *text.Text, win *window.Window, fm *frame.Frame, off text.Offset) text.Offset {
	view := layoutView(t, layoutContentWidth(fm.Width, lineCountForGutter(t), win.WindowOptions), 0, win.WindowOptions.Wrap)
	point := t.PointOf(off)
	row, _, ok := visualRowAt(view.Meta, point)
	if !ok {
		return off
	}
	vl := view.Meta[row]
	segLen := vl.EndCol - vl.StartCol
	col := vl.StartCol
	if segLen > 0 {
		col = vl.StartCol + segLen/2
	}
	return t.OffsetOf(text.Point{Line: vl.BufferLine, Col: col})
}

func pageScrollLines(fm *frame.Frame, half bool) int {
	height := fm.Height
	if height < 1 {
		height = frame.DefaultHeight
	}
	if half {
		n := height / 2
		if n < 1 {
			return 1
		}
		return n
	}
	n := height - 2
	if n < 1 {
		return 1
	}
	return n
}

func pageScrollVertical(
	t *text.Text,
	win *window.Window,
	fm *frame.Frame,
	off text.Offset,
	forward bool,
	half bool,
) text.Offset {
	delta := pageScrollLines(fm, half)
	if !forward {
		delta = -delta
	}
	return moveVerticalForWindow(t, win, fm, off, delta, false)
}
