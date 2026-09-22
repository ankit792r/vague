package editor

import (
	"vague/bonefire/buffer"
	frame "vague/bonefire/frame"
	"vague/bonefire/text"
	"vague/bonefire/window"
)

type caseChangeKind int

const (
	caseChangeNone caseChangeKind = iota
	caseChangeLower
	caseChangeUpper
	caseChangeToggle
)

func (e *Editor) applyCaseOnRange(
	frame *frame.Frame,
	win *window.Window,
	buf *buffer.Buffer,
	from, to text.Offset,
	kind caseChangeKind,
	target string,
) error {
	if buf.ReadOnly || to <= from {
		frame.Dirty = true
		return nil
	}

	data := buf.Text.Slice(from, to)
	switch kind {
	case caseChangeLower:
		data = bytesToLower(data)
	case caseChangeUpper:
		data = bytesToUpper(data)
	case caseChangeToggle:
		data = toggleCaseBytes(data)
	default:
		return nil
	}

	if _, err := buf.Replace(from, to, data); err != nil {
		return err
	}
	setWindowCursor(buf, win, from)
	frame.Dirty = true
	e.recordCaseChange(kind, target)
	return nil
}

func (e *Editor) applyCaseTarget(
	frame *frame.Frame,
	win *window.Window,
	buf *buffer.Buffer,
	kind caseChangeKind,
	target string,
) error {
	if from, to, ok := textObjectRange(buf.Text, win, target); ok {
		return e.applyCaseOnRange(frame, win, buf, from, to, kind, target)
	}
	if motion, ok := motionForKey(target); ok {
		from, to, _ := textRangeForMotion(buf.Text, win, opChange, motion)
		if to > from {
			return e.applyCaseOnRange(frame, win, buf, from, to, kind, target)
		}
	}
	frame.Dirty = true
	return nil
}

func bytesToLower(b []byte) []byte {
	out := make([]byte, len(b))
	copy(out, b)
	for i, c := range out {
		if c >= 'A' && c <= 'Z' {
			out[i] = c + ('a' - 'A')
		}
	}
	return out
}

func bytesToUpper(b []byte) []byte {
	out := make([]byte, len(b))
	copy(out, b)
	for i, c := range out {
		if c >= 'a' && c <= 'z' {
			out[i] = c - ('a' - 'A')
		}
	}
	return out
}
