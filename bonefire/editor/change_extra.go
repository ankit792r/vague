package editor

import (
	"vague/bonefire/buffer"
	frame "vague/bonefire/frame"
	"vague/bonefire/text"
	"vague/bonefire/window"
)

const defaultShiftWidth = 2

func (e *Editor) changeChars(frame *frame.Frame, win *window.Window, buf *buffer.Buffer, count int) error {
	if buf.ReadOnly {
		return buffer.ErrReadOnly
	}
	if count < 1 {
		count = 1
	}

	t := buf.Text
	at := windowCursor(win)
	end := at
	for i := 0; i < count; i++ {
		next := moveRight(t, end, 1, true)
		if next == end {
			break
		}
		end = next
	}

	if end > at {
		if _, err := buf.Delete(at, end); err != nil {
			return err
		}
	}
	setWindowCursor(buf, win, at)
	frame.Dirty = true
	e.recordChangeChars()
	return e.enterInsert(frame, win, buf, at)
}

func (e *Editor) replaceChar(frame *frame.Frame, win *window.Window, buf *buffer.Buffer, ch string) error {
	if buf.ReadOnly {
		return buffer.ErrReadOnly
	}
	if ch == "" || ch[0] == '<' {
		return nil
	}

	t := buf.Text
	at := windowCursor(win)
	end := moveRight(t, at, 1, true)
	if end <= at {
		return nil
	}

	if _, err := buf.Replace(at, end, []byte(ch[:1])); err != nil {
		return err
	}
	setWindowCursor(buf, win, moveRight(t, at, 1, false))
	frame.Dirty = true
	return nil
}

func (e *Editor) toggleCase(frame *frame.Frame, win *window.Window, buf *buffer.Buffer, count int) error {
	if buf.ReadOnly {
		return buffer.ErrReadOnly
	}
	if count < 1 {
		count = 1
	}

	t := buf.Text
	at := windowCursor(win)
	for i := 0; i < count; i++ {
		if at >= t.Len() {
			break
		}
		end := moveRight(t, at, 1, false)
		if end <= at {
			end = moveRight(t, at, 1, true)
		}
		if end <= at {
			break
		}
		seg := t.Slice(at, end)
		toggled := toggleCaseBytes(seg)
		if _, err := buf.Replace(at, end, toggled); err != nil {
			return err
		}
		at = end
	}
	setWindowCursor(buf, win, at)
	frame.Dirty = true
	return nil
}

func toggleCaseBytes(b []byte) []byte {
	out := make([]byte, len(b))
	copy(out, b)
	for i, c := range out {
		if c >= 'a' && c <= 'z' {
			out[i] = c - ('a' - 'A')
		} else if c >= 'A' && c <= 'Z' {
			out[i] = c + ('a' - 'A')
		}
	}
	return out
}

func (e *Editor) indentLines(frame *frame.Frame, win *window.Window, buf *buffer.Buffer, delta, count int) error {
	if buf.ReadOnly {
		return buffer.ErrReadOnly
	}
	if count < 1 {
		count = 1
	}

	t := buf.Text
	line := t.PointOf(windowCursor(win)).Line
	for i := 0; i < count; i++ {
		target := line + i
		if target >= t.LineCount() {
		 break
		}
		start := t.LineStart(target)
		if delta > 0 {
			ins := make([]byte, defaultShiftWidth)
			for j := range ins {
				ins[j] = ' '
			}
			if _, err := buf.Insert(start, ins); err != nil {
				return err
			}
		} else {
			lineBytes := t.Line(target)
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
