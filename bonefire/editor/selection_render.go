package editor

import (
	"vague/backbone/process"
	"vague/bonefire/text"
	"vague/bonefire/window"
)

func selectionInViewport(
	ed *Editor,
	win *window.Window,
	topLine int,
	meta []visualLine,
	lines []string,
	t *text.Text,
) *process.Selection {
	if ed.Mode != VisualMode && ed.Mode != VisualLineMode {
		return nil
	}

	anchor := ed.visualAnchor
	head := windowCursor(win)

	if ed.Mode == VisualLineMode {
		aLine := t.PointOf(anchor).Line
		hLine := t.PointOf(head).Line
		if aLine > hLine {
			aLine, hLine = hLine, aLine
		}

		startRow := -1
		endRow := -1
		endCol := 0
		for i, vl := range meta {
			if vl.BufferLine < aLine || vl.BufferLine > hLine {
				continue
			}

			vr := i - topLine
			if vr < 0 || vr >= len(lines) {
				continue
			}

			if startRow < 0 || vr < startRow {
				startRow = vr
			}
			if vr > endRow {
				endRow = vr
			}

			if col := len(lines[vr]); col > endCol {
				endCol = col
			}
			if endCol == 0 {
				endCol = 1
			}
		}

		if startRow < 0 {
			return nil
		}

		return &process.Selection{
			Visible:  true,
			Linewise: true,
			Start:    process.SelectionPoint{Row: startRow, Column: 0},
			End:      process.SelectionPoint{Row: endRow, Column: endCol - 1},
		}
	}

	startOff, endOff := anchor, head
	if anchor > head {
		startOff, endOff = head, anchor
	}

	startPt := t.PointOf(startOff)
	endPt := t.PointOf(endOff)

	sr, sc, ok1 := visualRowAt(meta, startPt)
	er, ec, ok2 := visualRowAt(meta, endPt)
	if !ok1 || !ok2 {
		return nil
	}

	sr -= topLine
	er -= topLine

	if sr >= len(lines) && er >= len(lines) {
		return nil
	}
	if sr < 0 {
		sr = 0
		sc = 0
	}
	if er >= len(lines) {
		er = len(lines) - 1
		if er >= 0 {
			lineLen := len(lines[er])
			if lineLen == 0 {
				ec = 0
			} else if ec >= lineLen {
				ec = lineLen - 1
			}
		}
	}

	if sr > er || (sr == er && sc > ec) {
		sr, sc, er, ec = er, ec, sr, sc
	}

	return &process.Selection{
		Visible: true,
		Start:   process.SelectionPoint{Row: sr, Column: sc},
		End:     process.SelectionPoint{Row: er, Column: ec},
	}
}
