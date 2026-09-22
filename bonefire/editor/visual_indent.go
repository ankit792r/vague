package editor

import (
	"vague/bonefire/buffer"
	frame "vague/bonefire/frame"
	"vague/bonefire/text"
	"vague/bonefire/window"
)

func (e *Editor) indentVisualSelection(
	frame *frame.Frame,
	win *window.Window,
	buf *buffer.Buffer,
	delta int,
) error {
	if buf.ReadOnly {
		return buffer.ErrReadOnly
	}

	t := buf.Text
	anchor := e.visualAnchor
	head := windowCursor(win)

	if e.Mode == VisualBlockMode {
		aPt := t.PointOf(anchor)
		hPt := t.PointOf(head)
		lineLo, lineHi := aPt.Line, hPt.Line
		if lineLo > lineHi {
			lineLo, lineHi = lineHi, lineLo
		}
		colLo := aPt.Col
		if hPt.Col < colLo {
			colLo = hPt.Col
		}
		for line := lineLo; line <= lineHi; line++ {
			start := t.OffsetOf(text.Point{Line: line, Col: colLo})
			if delta > 0 {
				ins := make([]byte, defaultShiftWidth)
				for i := range ins {
					ins[i] = ' '
				}
				if _, err := buf.Insert(start, ins); err != nil {
					return err
				}
			} else {
				lineBytes := t.Line(line)
				remove := 0
				for remove < defaultShiftWidth && colLo+remove < len(lineBytes) && lineBytes[colLo+remove] == ' ' {
					remove++
				}
				if remove > 0 {
					if _, err := buf.Delete(start, start+text.Offset(remove)); err != nil {
						return err
					}
				}
			}
		}
		frame.Dirty = true
		e.recordIndentChange(delta)
		return nil
	}

	from, to, linewise := e.visualRange(t, win)
	if to <= from && !linewise {
		frame.Dirty = true
		return nil
	}

	lineStart := t.PointOf(from).Line
	lineEnd := t.PointOf(to).Line
	if linewise {
		// already full lines
	} else if to > 0 {
		lineEnd = t.PointOf(to - 1).Line
	}
	for line := lineStart; line <= lineEnd; line++ {
		start := t.LineStart(line)
		if delta > 0 {
			ins := make([]byte, defaultShiftWidth)
			for i := range ins {
				ins[i] = ' '
			}
			if _, err := buf.Insert(start, ins); err != nil {
				return err
			}
		} else {
			lineBytes := t.Line(line)
			remove := 0
			for remove < defaultShiftWidth && remove < len(lineBytes) && lineBytes[remove] == ' ' {
				remove++
			}
			if remove > 0 {
				if _, err := buf.Delete(start, start+text.Offset(remove)); err != nil {
					return err
				}
			}
		}
	}
	frame.Dirty = true
	e.recordIndentChange(delta)
	return nil
}
