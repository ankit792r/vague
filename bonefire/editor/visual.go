package editor

import (
	"vague/bonefire/buffer"
	frame "vague/bonefire/frame"
	"vague/bonefire/text"
	"vague/bonefire/window"
)

func (e *Editor) leaveVisual(win *window.Window) {
	if e.Mode == VisualMode || e.Mode == VisualLineMode || e.Mode == VisualBlockMode {
		e.lastVisualMode = e.Mode
		e.lastVisualAnchor = e.visualAnchor
		if win != nil {
			e.lastVisualHead = windowCursor(win)
		}
	}
	if e.Mode != InsertMode {
		e.Mode = NormalMode
	}
	e.visualAnchor = 0
}

func (e *Editor) reselectLastVisual(
	frame *frame.Frame,
	win *window.Window,
	buf *buffer.Buffer,
) error {
	if e.lastVisualMode == NormalMode {
		frame.Dirty = true
		return nil
	}
	e.Mode = e.lastVisualMode
	e.visualAnchor = e.lastVisualAnchor
	setWindowCursor(buf, win, e.lastVisualHead)
	e.clearPendingOp()
	e.pendingKey = ""
	frame.Dirty = true
	return nil
}

func (e *Editor) enterVisualChar(win *window.Window) {
	e.Mode = VisualMode
	e.visualAnchor = windowCursor(win)
	e.clearPendingOp()
	e.pendingKey = ""
}

func (e *Editor) enterVisualLine(win *window.Window) {
	e.Mode = VisualLineMode
	e.visualAnchor = windowCursor(win)
	e.clearPendingOp()
	e.pendingKey = ""
}

func (e *Editor) enterVisualBlock(win *window.Window) {
	e.Mode = VisualBlockMode
	e.visualAnchor = windowCursor(win)
	e.clearPendingOp()
	e.pendingKey = ""
}

func (e *Editor) visualRange(t *text.Text, win *window.Window) (from, to text.Offset, linewise bool) {
	anchor := e.visualAnchor
	head := windowCursor(win)

	if e.Mode == VisualLineMode {
		aLine := t.PointOf(anchor).Line
		hLine := t.PointOf(head).Line
		if aLine > hLine {
			aLine, hLine = hLine, aLine
		}
		from = t.LineStart(aLine)
		_, to = lineChangeRange(t, hLine)
		return from, to, true
	}

	if e.Mode == VisualBlockMode {
		aPt := t.PointOf(anchor)
		hPt := t.PointOf(head)
		lineLo, lineHi := aPt.Line, hPt.Line
		if lineLo > lineHi {
			lineLo, lineHi = lineHi, lineLo
		}
		colLo, colHi := aPt.Col, hPt.Col
		if colLo > colHi {
			colLo, colHi = colHi, colLo
		}
		from = t.OffsetOf(text.Point{Line: lineLo, Col: colLo})
		endCol := colHi
		lineBytes := t.Line(lineHi)
		if endCol >= len(lineBytes) {
			to = t.LineEnd(lineHi)
		} else {
			to = t.OffsetOf(text.Point{Line: lineHi, Col: endCol + 1})
		}
		if to <= from {
			to = from + 1
		}
		return from, to, false
	}

	start := anchor
	end := head
	if anchor > head {
		start, end = head, anchor
	}

	from = start
	if end < t.Len() {
		to = forwardChar(t, end)
	} else {
		to = t.Len()
	}
	if to <= from {
		to = from + 1
		if to > t.Len() {
			to = t.Len()
		}
	}

	return from, to, false
}

func (e *Editor) applyVisualOperator(
	frame *frame.Frame,
	win *window.Window,
	buf *buffer.Buffer,
	op opKind,
) error {
	e.lastVisualMode = e.Mode
	e.lastVisualAnchor = e.visualAnchor
	e.lastVisualHead = windowCursor(win)

	t := buf.Text
	from, to, linewise := e.visualRange(t, win)
	if to <= from {
		e.leaveVisual(win)
		frame.Dirty = true
		return nil
	}

	landing := clampToLine(t, from)
	if linewise {
		landing = t.LineStart(t.PointOf(from).Line)
	}

	err := e.applyOperatorRange(frame, win, buf, op, from, to, linewise, landing)
	if err != nil {
		return err
	}

	if (op == opDelete || op == opChange) && to > from {
		e.recordVisualOperatorChange(op, linewise)
	}

	e.visualAnchor = 0
	if e.Mode != InsertMode {
		e.leaveVisual(win)
	}
	return nil
}

func (e *Editor) visualKey(frame *frame.Frame, win *window.Window, buf *buffer.Buffer, keys string) error {
	t := buf.Text
	at := windowCursor(win)

	switch keys {
	case "<Esc>", "v":
		e.leaveVisual(win)
		frame.Dirty = true
		return nil
	case ">":
		return e.indentVisualSelection(frame, win, buf, 1)
	case "<":
		return e.indentVisualSelection(frame, win, buf, -1)
	case "=":
		return e.indentVisualSelection(frame, win, buf, 1)
	case "d", "x":
		return e.applyVisualOperator(frame, win, buf, opDelete)
	case "y":
		return e.applyVisualOperator(frame, win, buf, opYank)
	case "c":
		return e.applyVisualOperator(frame, win, buf, opChange)
	case "h", "<Left>":
		setWindowCursor(buf, win, moveLeft(t, at, 1))
	case "l", "<Right>", "<Space>":
		setWindowCursor(buf, win, moveRight(t, at, 1, false))
	case "j", "<Down>":
		setWindowCursor(buf, win, moveVerticalForWindow(t, win, frame, at, 1, false))
	case "k", "<Up>":
		setWindowCursor(buf, win, moveVerticalForWindow(t, win, frame, at, -1, false))
	case "w":
		setWindowCursor(buf, win, moveWordForward(t, at, 1))
	case "b":
		setWindowCursor(buf, win, moveWordBack(t, at, 1))
	case "e":
		setWindowCursor(buf, win, moveWordEnd(t, at, 1))
	case "0":
		setWindowCursor(buf, win, moveToLineStart(t, at))
	case "$":
		setWindowCursor(buf, win, moveToLineEnd(t, at, false))
	case "G":
		setWindowCursor(buf, win, moveToBufferLine(t, win, t.LineCount()-1))
	default:
		return nil
	}

	view := layoutViewForWindow(t, win, frame)
	rememberColumn(buf, win, view)
	frame.Dirty = true
	return nil
}
