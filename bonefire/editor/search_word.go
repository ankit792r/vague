package editor

import (
	"vague/bonefire/buffer"
	frame "vague/bonefire/frame"
	"vague/bonefire/text"
	"vague/bonefire/window"
)

// wordSearchSlice returns the keyword under/at cursor for * / # search.
func wordSearchSlice(t *text.Text, off text.Offset, wholeWord bool) []byte {
	if off > t.Len() {
		off = t.Len()
	}

	start := off
	if wholeWord {
		start = wordStartAt(t, off)
	} else {
		cls, _, _ := charAt(t, off)
		if cls == wcWhitespace {
			start = moveWordForwardOnce(t, off)
			start = wordStartAt(t, start)
		} else {
			start = wordStartAt(t, off)
		}
	}

	end := wordEndExclusive(t, start)
	if end <= start {
		return nil
	}
	return t.Slice(start, end)
}

func wordStartAt(t *text.Text, off text.Offset) text.Offset {
	if off >= t.Len() {
		off = lastBufferChar(t)
	}

	cls, _, _ := charAt(t, off)
	if cls == wcWhitespace {
		off = moveWordForwardOnce(t, off)
	}

	cls, _, _ = charAt(t, off)
	if cls == wcWhitespace {
		return off
	}

	start := off
	for start > 0 {
		prev := back(t, start)
		if charClassAt(t, prev) != cls {
			break
		}
		start = prev
	}
	return start
}

func charClassAt(t *text.Text, off text.Offset) wordClass {
	cls, _, _ := charAt(t, off)
	return cls
}

func (e *Editor) searchWord(
	frame *frame.Frame,
	win *window.Window,
	buf *buffer.Buffer,
	forward bool,
	wholeWord bool,
) error {
	word := wordSearchSlice(buf.Text, windowCursor(win), wholeWord)
	if len(word) == 0 {
		return ErrPatternNotFound
	}

	e.clearSearchContext()
	e.searchPattern = string(word)
	e.searchForward = forward

	return e.runSearch(frame, win, buf, forward)
}
