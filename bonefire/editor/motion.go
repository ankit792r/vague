package editor

import (
	"unicode/utf8"

	"vague/bonefire/buffer"
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

func rememberColumn(buf *buffer.Buffer, win *window.Window) {
	win.DesiredCol = windowPoint(buf, win).Col
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
