package editor

import (
	"vague/bonefire/buffer"
	frame "vague/bonefire/frame"
	"vague/bonefire/window"
)

func (e *Editor) normalKey(frame *frame.Frame, win *window.Window, buf *buffer.Buffer, keys string) error {
	if e.consumeCountDigit(keys) {
		return nil
	}

	if handled, err := e.consumePendingCharFind(frame, win, buf, keys); handled {
		return err
	}

	t := buf.Text
	at := windowCursor(win)

	if e.pendingKey == "g" && keys == "g" {
		e.pendingKey = ""
		line := 0
		if c, explicit := e.takeCountOptional(); explicit {
			line = clampInt(c-1, 0, t.LineCount()-1)
		}
		setWindowCursor(buf, win, moveToBufferLine(t, win, line))
		view := layoutViewForWindow(t, win, frame)
		rememberColumn(buf, win, view)
		frame.Dirty = true
		return nil
	}

	if e.pendingKey == "g" && (keys == "*" || keys == "#") {
		e.pendingKey = ""
		forward := keys == "*"
		return e.searchWord(frame, win, buf, forward, true)
	}

	if e.pendingKey == "z" && len(keys) == 1 {
		e.pendingKey = ""
		scrollCursorLine(win, frame, buf, rune(keys[0]))
		frame.Dirty = true
		return nil
	}

	if keys == "*" || keys == "#" {
		forward := keys == "*"
		return e.searchWord(frame, win, buf, forward, false)
	}

	if e.pendingOp == opDelete && keys == "d" {
		e.clearPendingOp()
		count := e.takeCount()
		return e.applyOperatorMotion(frame, win, buf, opDelete, motionLine, count)
	}
	if e.pendingOp == opYank && keys == "y" {
		e.clearPendingOp()
		count := e.takeCount()
		return e.applyOperatorMotion(frame, win, buf, opYank, motionLine, count)
	}
	if e.pendingOp == opChange && keys == "c" {
		e.clearPendingOp()
		count := e.takeCount()
		return e.applyOperatorMotion(frame, win, buf, opChange, motionLine, count)
	}

	if e.pendingOp != opNone {
		op := e.pendingOp
		if motion, ok := motionForKey(keys); ok {
			e.clearPendingOp()
			count := e.takeCount()
			return e.applyOperatorMotion(frame, win, buf, op, motion, count)
		}
		e.clearPendingOp()
		e.clearPendingCount()
	}

	if e.pendingKey != "" {
		e.pendingKey = ""
	}

	switch keys {
	case "d":
		e.pendingOp = opDelete
		return nil
	case "y":
		e.pendingOp = opYank
		return nil
	case "c":
		e.pendingOp = opChange
		return nil
	case "D":
		return e.applyOperatorMotion(frame, win, buf, opDelete, motionToEOL, e.takeCount())
	case "p":
		return e.pasteAfter(frame, win, buf)
	case "P":
		return e.pasteBefore(frame, win, buf)
	case "J":
		return e.joinLinesRepeat(frame, win, buf, e.takeCount())
	case "x", "<Del>":
		return e.deleteCharRepeat(frame, win, buf, e.takeCount())
	case "v":
		e.enterVisualChar(win)
		frame.Dirty = true
		return nil
	case "V":
		e.enterVisualLine(win)
		frame.Dirty = true
		return nil
	case "i":
		return e.enterInsert(frame, win, buf, at)
	case "I":
		return e.enterInsert(frame, win, buf, firstNonBlank(t, at))
	case "a":
		return e.enterInsert(frame, win, buf, moveRight(t, at, 1, true))
	case "A":
		return e.enterInsert(frame, win, buf, moveToLineEnd(t, at, true))
	case "g":
		e.pendingKey = "g"
		return nil
	case "G":
		line := t.LineCount() - 1
		if c, explicit := e.takeCountOptional(); explicit {
			line = clampInt(c-1, 0, t.LineCount()-1)
		}
		setWindowCursor(buf, win, moveToBufferLine(t, win, line))
	case "H":
		setWindowCursor(buf, win, moveToScreenLine(t, win, frame, 'H'))
	case "M":
		setWindowCursor(buf, win, moveToScreenLine(t, win, frame, 'M'))
	case "L":
		setWindowCursor(buf, win, moveToScreenLine(t, win, frame, 'L'))
	case "z":
		e.pendingKey = "z"
		return nil
	case "0":
		setWindowCursor(buf, win, moveToLineStart(t, at))
	case "^":
		setWindowCursor(buf, win, firstNonBlank(t, at))
	case "_":
		setWindowCursor(buf, win, moveUnderscore(t, at, e.takeCount()))
	case "|":
		setWindowCursor(buf, win, moveToColumn(t, at, e.takeCount()))
	case "$":
		setWindowCursor(buf, win, moveToLineEnd(t, at, false))
	case "o":
		return e.openLine(frame, win, buf, false)
	case "O":
		return e.openLine(frame, win, buf, true)
	case "n":
		return e.RepeatSearch(frame, win, buf, true)
	case "N":
		return e.RepeatSearch(frame, win, buf, false)
	case ".":
		return e.repeatLastChange(frame, win, buf)
	case "u":
		return e.undoTo(frame, win, buf, false)
	case "<C-r>":
		return e.undoTo(frame, win, buf, true)
	case "h", "<Left>":
		setWindowCursor(buf, win, moveLeft(t, at, e.takeCount()))
	case "w":
		setWindowCursor(buf, win, moveWordForward(t, at, e.takeCount()))
	case "b":
		setWindowCursor(buf, win, moveWordBack(t, at, e.takeCount()))
	case "e":
		setWindowCursor(buf, win, moveWordEnd(t, at, e.takeCount()))
	case "l", "<Right>", "<Space>":
		setWindowCursor(buf, win, moveRight(t, at, e.takeCount(), false))
	case "j", "<Down>":
		setWindowCursor(buf, win, moveVerticalForWindow(t, win, frame, at, e.takeCount(), false))
	case "k", "<Up>":
		setWindowCursor(buf, win, moveVerticalForWindow(t, win, frame, at, -e.takeCount(), false))
	case "{":
		setWindowCursor(buf, win, moveParagraphBackward(t, at, e.takeCount()))
	case "}":
		setWindowCursor(buf, win, moveParagraphForward(t, at, e.takeCount()))
	case "(":
		setWindowCursor(buf, win, moveSentenceBackward(t, at, e.takeCount()))
	case ")":
		setWindowCursor(buf, win, moveSentenceForward(t, at, e.takeCount()))
	case "%":
		setWindowCursor(buf, win, moveMatchingBracket(t, at, e.takeCount()))
	case "[":
		setWindowCursor(buf, win, moveParagraphBackward(t, at, e.takeCount()))
	case "]":
		setWindowCursor(buf, win, moveParagraphForward(t, at, e.takeCount()))
	case "f":
		e.beginPendingCharFind(charFindF)
		return nil
	case "F":
		e.beginPendingCharFind(charFindBigF)
		return nil
	case "t":
		e.beginPendingCharFind(charFindT)
		return nil
	case "T":
		e.beginPendingCharFind(charFindBigT)
		return nil
	case ";":
		if e.lastCharFindKind != charFindNone {
			return e.executeCharFind(frame, win, buf, e.lastCharFindKind, e.lastCharFindRune, e.takeCount())
		}
		return nil
	case ",":
		if e.lastCharFindKind != charFindNone {
			kind := oppositeCharFind(e.lastCharFindKind)
			return e.executeCharFind(frame, win, buf, kind, e.lastCharFindRune, e.takeCount())
		}
		return nil
	default:
		return nil
	}

	view := layoutViewForWindow(t, win, frame)
	rememberColumn(buf, win, view)
	frame.Dirty = true
	return nil
}
