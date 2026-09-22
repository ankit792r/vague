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
	repeatCaseChange
	repeatFormatChange
	repeatNumberChange
	repeatTextObject
	repeatIndentLines
	repeatChangeChars
)

type lastChange struct {
	kind     repeatKind
	op       opKind
	motion   motionKind
	linewise bool

	caseKind caseChangeKind
	target   string
	numberDelta int
	indentDelta int
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

func (e *Editor) recordCaseChange(kind caseChangeKind, target string) {
	e.lastChange = lastChange{kind: repeatCaseChange, caseKind: kind, target: target}
}

func (e *Editor) recordFormatChange(target string) {
	e.lastChange = lastChange{kind: repeatFormatChange, target: target}
}

func (e *Editor) recordNumberChange(delta int) {
	e.lastChange = lastChange{kind: repeatNumberChange, numberDelta: delta}
}

func (e *Editor) recordTextObjectChange(op opKind, target string) {
	e.lastChange = lastChange{kind: repeatTextObject, op: op, target: target}
}

func (e *Editor) recordIndentChange(delta int) {
	e.lastChange = lastChange{kind: repeatIndentLines, indentDelta: delta}
}

func (e *Editor) recordChangeChars() {
	e.lastChange = lastChange{kind: repeatChangeChars}
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
	case repeatCaseChange:
		return e.applyCaseTarget(frame, win, buf, e.lastChange.caseKind, e.lastChange.target)
	case repeatFormatChange:
		if from, to, ok := textObjectRange(buf.Text, win, e.lastChange.target); ok {
			return e.formatRange(frame, win, buf, from, to, e.lastChange.target)
		}
		if motion, ok := motionForKey(e.lastChange.target); ok {
			from, to, _ := textRangeForMotion(buf.Text, win, opChange, motion)
			if to > from {
				return e.formatRange(frame, win, buf, from, to, e.lastChange.target)
			}
		}
		return nil
	case repeatNumberChange:
		return e.changeNumberAtCursor(frame, win, buf, e.lastChange.numberDelta)
	case repeatTextObject:
		return e.applyOperatorTextObject(frame, win, buf, e.lastChange.op, e.lastChange.target, 1)
	case repeatIndentLines:
		return e.indentLines(frame, win, buf, e.lastChange.indentDelta, 1)
	case repeatChangeChars:
		return e.changeChars(frame, win, buf, 1)
	default:
		return nil
	}
}
