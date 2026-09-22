package editor

import (
	"vague/bonefire/buffer"
	frame "vague/bonefire/frame"
	"vague/bonefire/window"
)

type repeatKind int

const (
	repeatNone repeatKind = iota
	repeatDeleteChar
	repeatJoin
	repeatOperator
	repeatVisualOperator
)

type lastChange struct {
	kind     repeatKind
	op       opKind
	motion   motionKind
	linewise bool
}

func (e *Editor) recordOperatorChange(op opKind, motion motionKind) {
	e.lastChange = lastChange{
		kind:   repeatOperator,
		op:     op,
		motion: motion,
	}
}

func (e *Editor) recordVisualOperatorChange(op opKind, linewise bool) {
	e.lastChange = lastChange{
		kind:     repeatVisualOperator,
		op:       op,
		linewise: linewise,
	}
}

func (e *Editor) recordDeleteCharChange() {
	e.lastChange = lastChange{kind: repeatDeleteChar}
}

func (e *Editor) recordJoinChange() {
	e.lastChange = lastChange{kind: repeatJoin}
}

func (e *Editor) repeatLastChange(
	frame *frame.Frame,
	win *window.Window,
	buf *buffer.Buffer,
) error {
	switch e.lastChange.kind {
	case repeatDeleteChar:
		return e.deleteChar(frame, win, buf)
	case repeatJoin:
		return e.joinLines(frame, win, buf)
	case repeatOperator:
		return e.applyOperatorMotion(frame, win, buf, e.lastChange.op, e.lastChange.motion, 1)
	case repeatVisualOperator:
		motion := motionWord
		if e.lastChange.linewise {
			motion = motionLine
		}
		return e.applyOperatorMotion(frame, win, buf, e.lastChange.op, motion, 1)
	default:
		return nil
	}
}
