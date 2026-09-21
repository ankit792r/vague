package editor

import (
	"unicode/utf8"

	"vague/bonefire/buffer"
	frame "vague/bonefire/frame"
	"vague/bonefire/text"
	"vague/bonefire/window"
)

func windowCursor(win *window.Window) text.Offset {
	if win.Cursor == nil {
		return 0
	}

	return win.Cursor.Off
}

func setWindowCursor(buf *buffer.Buffer, win *window.Window, at text.Offset) {
	if win.Cursor == nil {
		win.Cursor = buf.Text.AddMarker(at, text.GravityRight)
		return
	}

	win.Cursor.Off = buf.Text.Clamp(at)
}

func windowPoint(buf *buffer.Buffer, win *window.Window) text.Point {
	return buf.Text.PointOf(windowCursor(win))
}

func rememberColumn(buf *buffer.Buffer, win *window.Window, view *viewLayout) {
	point := windowPoint(buf, win)
	if view != nil {
		if _, col, ok := visualRowAt(view.Meta, point); ok {
			win.DesiredCol = col
			return
		}
	}

	win.DesiredCol = point.Col
}

func moveToLineStart(t *text.Text, off text.Offset) text.Offset {
	point := t.PointOf(off)
	return t.LineStart(point.Line)
}

func moveToLineEnd(t *text.Text, off text.Offset, past bool) text.Offset {
	point := t.PointOf(off)
	line := t.Line(point.Line)
	if past {
		return t.LineEnd(point.Line)
	}

	return t.OffsetOf(text.Point{Line: point.Line, Col: lastColumn(line)})
}

func firstNonBlank(t *text.Text, off text.Offset) text.Offset {
	point := t.PointOf(off)
	line := t.Line(point.Line)
	col := 0

	for col < len(line) {
		r, size := utf8.DecodeRune(line[col:])
		if r != ' ' && r != '\t' {
			return t.OffsetOf(text.Point{Line: point.Line, Col: col})
		}
		col += size
	}

	return t.LineStart(point.Line)
}

func moveToBufferLine(t *text.Text, win *window.Window, line int) text.Offset {
	line = clampInt(line, 0, t.LineCount()-1)
	target := t.Line(line)
	col := win.DesiredCol
	if col > lastColumn(target) {
		col = lastColumn(target)
	}

	return t.OffsetOf(text.Point{Line: line, Col: col})
}

func moveLeft(t *text.Text, off text.Offset, count int) text.Offset {
	point := t.PointOf(off)
	line := t.Line(point.Line)

	for ; count > 0 && point.Col > 0; count-- {
		_, size := utf8.DecodeLastRune(line[:point.Col])
		point.Col -= size
	}

	return t.OffsetOf(point)
}

func moveRight(t *text.Text, off text.Offset, count int, past bool) text.Offset {
	point := t.PointOf(off)
	line := t.Line(point.Line)

	limit := len(line)
	if !past {
		limit = lastColumn(line)
	}

	for ; count > 0 && point.Col < limit; count-- {
		_, size := utf8.DecodeRune(line[point.Col:])
		if point.Col+size > limit {
			break
		}
		point.Col += size
	}

	return t.OffsetOf(point)
}

// moveVertical moves by buffer lines when wrap is off.
func moveVertical(t *text.Text, off text.Offset, desiredCol, delta int) text.Offset {
	point := t.PointOf(off)

	target := clampInt(point.Line+delta, 0, t.LineCount()-1)
	if target == point.Line {
		return off
	}

	line := t.Line(target)
	col := desiredCol
	if col > lastColumn(line) {
		col = lastColumn(line)
	}

	return t.OffsetOf(text.Point{Line: target, Col: col})
}

// moveVerticalVisual moves by screen rows through wrapped lines.
func moveVerticalVisual(view viewLayout, t *text.Text, off text.Offset, desiredCol, delta int, past bool) text.Offset {
	point := t.PointOf(off)
	currentRow, _, ok := visualRowAt(view.Meta, point)
	if !ok {
		return off
	}

	targetRow := currentRow + delta
	if targetRow < 0 || targetRow >= len(view.Meta) {
		return off
	}

	vl := view.Meta[targetRow]
	col := clampVisualCol(vl, desiredCol, past)

	return t.OffsetOf(text.Point{Line: vl.BufferLine, Col: vl.StartCol + col})
}

func clampVisualCol(vl visualLine, desiredCol int, past bool) int {
	segLen := vl.EndCol - vl.StartCol
	if segLen <= 0 {
		return 0
	}

	if past {
		if desiredCol > segLen {
			return segLen
		}
		return desiredCol
	}

	if desiredCol >= segLen {
		return segLen - 1
	}

	if desiredCol < 0 {
		return 0
	}

	return desiredCol
}

func back(t *text.Text, off text.Offset) text.Offset {
	if off <= 0 {
		return 0
	}

	point := t.PointOf(off)
	if point.Col == 0 {
		if point.Line == 0 {
			return 0
		}

		return t.LineEnd(point.Line - 1)
	}

	line := t.Line(point.Line)
	_, size := utf8.DecodeLastRune(line[:point.Col])

	return off - text.Offset(size)
}

func lastColumn(line []byte) int {
	if len(line) == 0 {
		return 0
	}

	_, size := utf8.DecodeLastRune(line)
	return len(line) - size
}

func clampInt(v, lo, hi int) int {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

func moveVerticalForWindow(
	t *text.Text,
	win *window.Window,
	frame *frame.Frame,
	off text.Offset,
	delta int,
	past bool,
) text.Offset {
	if !win.WindowOptions.Wrap {
		return moveVertical(t, off, win.DesiredCol, delta)
	}

	view := layoutView(t, frame.Width, 0, true)
	return moveVerticalVisual(view, t, off, win.DesiredCol, delta, past)
}

func layoutViewForWindow(t *text.Text, win *window.Window, frame *frame.Frame) *viewLayout {
	if !win.WindowOptions.Wrap {
		return nil
	}

	view := layoutView(t, frame.Width, 0, true)
	return &view
}
