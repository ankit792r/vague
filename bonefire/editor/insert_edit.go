package editor

import (
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"

	"vague/bonefire/buffer"
	frame "vague/bonefire/frame"
	"vague/bonefire/text"
	"vague/bonefire/window"
)

func (e *Editor) deleteInsertRange(
	frame *frame.Frame,
	win *window.Window,
	buf *buffer.Buffer,
	from, to text.Offset,
) error {
	if to <= from {
		return nil
	}
	e.noteChangeAt(buf, from)
	if _, err := buf.Delete(from, to); err != nil {
		return err
	}
	setWindowCursor(buf, win, from)
	view := layoutViewForWindow(buf.Text, win, frame)
	rememberColumn(buf, win, view)
	frame.Dirty = true
	return nil
}

func (e *Editor) insertDeleteWordBack(frame *frame.Frame, win *window.Window, buf *buffer.Buffer) error {
	t := buf.Text
	at := windowCursor(win)
	from := moveWordBackOnce(t, at)
	return e.deleteInsertRange(frame, win, buf, from, at)
}

func (e *Editor) insertDeleteToLineStart(frame *frame.Frame, win *window.Window, buf *buffer.Buffer) error {
	t := buf.Text
	at := windowCursor(win)
	from := moveToLineStart(t, at)
	return e.deleteInsertRange(frame, win, buf, from, at)
}

func (e *Editor) insertMoveLineStart(frame *frame.Frame, win *window.Window, buf *buffer.Buffer) error {
	t := buf.Text
	setWindowCursor(buf, win, moveToLineStart(t, windowCursor(win)))
	frame.Dirty = true
	return nil
}

func (e *Editor) insertMoveLineEnd(frame *frame.Frame, win *window.Window, buf *buffer.Buffer) error {
	t := buf.Text
	setWindowCursor(buf, win, moveToLineEnd(t, windowCursor(win), true))
	frame.Dirty = true
	return nil
}

func (e *Editor) insertAdjustIndent(frame *frame.Frame, win *window.Window, buf *buffer.Buffer, delta int) error {
	if buf.ReadOnly {
		return buffer.ErrReadOnly
	}
	return e.indentLines(frame, win, buf, delta, 1)
}

func (e *Editor) insertRegisterPaste(
	frame *frame.Frame,
	win *window.Window,
	buf *buffer.Buffer,
	data []byte,
) error {
	if !e.fileOpts.Paste {
		return e.insertBytes(frame, win, buf, data)
	}
	at := windowCursor(win)
	e.noteChangeAt(buf, at)
	delta, err := buf.Insert(at, data)
	if err != nil {
		return err
	}
	setWindowCursor(buf, win, delta.NewEnd)
	view := layoutViewForWindow(buf.Text, win, frame)
	rememberColumn(buf, win, view)
	frame.Dirty = true
	return nil
}

func wordPrefixBefore(t *text.Text, at text.Offset) (text.Offset, string) {
	point := t.PointOf(at)
	line := t.Line(point.Line)
	col := point.Col
	if col > len(line) {
		col = len(line)
	}
	startCol := col
	for startCol > 0 {
		r, size := utf8.DecodeLastRuneInString(string(line[:startCol]))
		if !isWordRune(r) {
			break
		}
		startCol -= size
	}
	start := t.OffsetOf(text.Point{Line: point.Line, Col: startCol})
	return start, string(line[startCol:col])
}

func isWordRune(r rune) bool {
	return r == '_' || unicode.IsLetter(r) || unicode.IsDigit(r)
}

func findWordCompletions(t *text.Text, prefix string) []string {
	if prefix == "" {
		return nil
	}
	lowerPrefix := strings.ToLower(prefix)
	seen := make(map[string]struct{})
	var matches []string
	scan := t.Bytes()
	i := 0
	for i < len(scan) {
		r, size := utf8.DecodeRune(scan[i:])
		if !isWordRune(r) {
			i += size
			continue
		}
		j := i
		for j < len(scan) {
			r2, sz := utf8.DecodeRune(scan[j:])
			if !isWordRune(r2) {
				break
			}
			j += sz
		}
		word := string(scan[i:j])
		if strings.HasPrefix(strings.ToLower(word), lowerPrefix) && word != prefix {
			if _, ok := seen[word]; !ok {
				seen[word] = struct{}{}
				matches = append(matches, word)
			}
		}
		i = j
	}
	sort.Strings(matches)
	return matches
}

func (e *Editor) beginInsertCompletion(win *window.Window, buf *buffer.Buffer) {
	t := buf.Text
	at := windowCursor(win)
	start, prefix := wordPrefixBefore(t, at)
	e.insertCompleteActive = true
	e.insertCompleteStart = start
	e.insertCompletePrefix = prefix
	e.insertCompleteMatches = findWordCompletions(t, prefix)
	e.insertCompleteIndex = -1
}

func (e *Editor) insertCompletionCycle(
	frame *frame.Frame,
	win *window.Window,
	buf *buffer.Buffer,
	dir int,
) error {
	if !e.insertCompleteActive || len(e.insertCompleteMatches) == 0 {
		frame.Dirty = true
		return nil
	}
	e.insertCompleteIndex += dir
	if e.insertCompleteIndex < 0 {
		e.insertCompleteIndex = len(e.insertCompleteMatches) - 1
	}
	if e.insertCompleteIndex >= len(e.insertCompleteMatches) {
		e.insertCompleteIndex = 0
	}
	word := e.insertCompleteMatches[e.insertCompleteIndex]
	t := buf.Text
	at := windowCursor(win)
	from := e.insertCompleteStart
	if from > at {
		from, _ = wordPrefixBefore(t, at)
	}
	if at > from {
		e.noteChangeAt(buf, from)
		if _, err := buf.Delete(from, at); err != nil {
			return err
		}
	}
	e.noteChangeAt(buf, from)
	delta, err := buf.Insert(from, []byte(word))
	if err != nil {
		return err
	}
	setWindowCursor(buf, win, delta.NewEnd)
	e.insertCompleteStart = from
	setFrameInfoEcho(frame, word)
	frame.Dirty = true
	return nil
}

func (e *Editor) visualReplaceChar(
	frame *frame.Frame,
	win *window.Window,
	buf *buffer.Buffer,
	ch string,
) error {
	if buf.ReadOnly || ch == "" || ch[0] == '<' {
		return nil
	}
	t := buf.Text
	from, to, linewise := e.visualRange(t, win)
	if to <= from {
		e.leaveVisual(win)
		frame.Dirty = true
		return nil
	}

	repl := buildVisualReplace(t.Slice(from, to), ch[0])
	e.noteChangeAt(buf, from)
	if _, err := buf.Replace(from, to, repl); err != nil {
		return err
	}

	landing := from
	if linewise {
		landing = t.LineStart(t.PointOf(from).Line)
	}
	setWindowCursor(buf, win, landing)
	e.leaveVisual(win)
	frame.Dirty = true
	return nil
}

func buildVisualReplace(segment []byte, ch byte) []byte {
	out := make([]byte, len(segment))
	for i, b := range segment {
		if b == '\n' {
			out[i] = '\n'
		} else {
			out[i] = ch
		}
	}
	return out
}
