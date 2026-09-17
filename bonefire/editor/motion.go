package editor

import "vague/bonefire/window"

func windowPoint(win *window.Window) Point {
	return Point{Line: win.CursorLine, Col: win.CursorCol}
}

func setWindowPoint(win *window.Window, p Point) {
	win.CursorLine = p.Line
	win.CursorCol = p.Col
}

func rememberColumn(win *window.Window, p Point) {
	win.DesiredCol = p.Col
}

func moveLeft(lines []string, p Point, count int) Point {
	if len(lines) == 0 {
		return p
	}

	p.Line = clamp(p.Line, 0, len(lines)-1)
	runes := []rune(lines[p.Line])
	p.Col = clamp(p.Col, 0, lastRuneCol(runes))

	for ; count > 0 && p.Col > 0; count-- {
		p.Col--
	}

	return p
}

func moveRight(lines []string, p Point, count int) Point {
	if len(lines) == 0 {
		return p
	}

	p.Line = clamp(p.Line, 0, len(lines)-1)
	runes := []rune(lines[p.Line])
	limit := lastRuneCol(runes)
	p.Col = clamp(p.Col, 0, limit)

	for ; count > 0 && p.Col < limit; count-- {
		p.Col++
	}

	return p
}

func moveVertical(lines []string, p Point, desiredCol, delta int) Point {
	if len(lines) == 0 {
		return Point{}
	}

	targetLine := clamp(p.Line+delta, 0, len(lines)-1)
	if targetLine == p.Line {
		return p
	}

	runes := []rune(lines[targetLine])
	col := desiredCol
	if col > lastRuneCol(runes) {
		col = lastRuneCol(runes)
	}

	return Point{Line: targetLine, Col: col}
}

func lastRuneCol(runes []rune) int {
	if len(runes) == 0 {
		return 0
	}

	return len(runes) - 1
}

func clamp(v, lo, hi int) int {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}
