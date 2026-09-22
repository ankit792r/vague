package editor

import (
	"bytes"
	"strings"
	"unicode"

	"vague/bonefire/buffer"
	frame "vague/bonefire/frame"
	"vague/bonefire/text"
	"vague/bonefire/window"
)

const defaultTextWidth = 78

func (e *Editor) formatRange(
	frame *frame.Frame,
	win *window.Window,
	buf *buffer.Buffer,
	from, to text.Offset,
	target string,
) error {
	if buf.ReadOnly || to <= from {
		frame.Dirty = true
		return nil
	}

	chunk := buf.Text.Slice(from, to)
	formatted := formatParagraphBytes(chunk, defaultTextWidth)
	if bytes.Equal(chunk, formatted) {
		frame.Dirty = true
		return nil
	}

	if _, err := buf.Replace(from, to, formatted); err != nil {
		return err
	}
	setWindowCursor(buf, win, from)
	frame.Dirty = true
	e.recordFormatChange(target)
	return nil
}

func formatParagraphBytes(b []byte, width int) []byte {
	if width < 1 {
		width = defaultTextWidth
	}
	text := string(b)
	words := strings.Fields(text)
	if len(words) == 0 {
		return b
	}

	var out strings.Builder
	lineLen := 0
	for i, w := range words {
		wLen := len(w)
		if i == 0 {
			out.WriteString(w)
			lineLen = wLen
			continue
		}
		if lineLen+1+wLen <= width {
			out.WriteByte(' ')
			out.WriteString(w)
			lineLen += 1 + wLen
			continue
		}
		out.WriteByte('\n')
		out.WriteString(w)
		lineLen = wLen
	}
	return []byte(out.String())
}

func scanNumberAt(t *text.Text, at text.Offset) (start, end text.Offset, value int, ok bool) {
	if at >= t.Len() {
		return 0, 0, 0, false
	}
	point := t.PointOf(at)
	line := t.Line(point.Line)
	col := point.Col
	if col >= len(line) {
		col = len(line) - 1
	}
	if col < 0 {
		return 0, 0, 0, false
	}

	startCol := col
	for startCol > 0 && unicode.IsDigit(rune(line[startCol-1])) {
		startCol--
	}
	endCol := col
	for endCol < len(line) && unicode.IsDigit(rune(line[endCol])) {
		endCol++
	}
	if endCol <= startCol {
		return 0, 0, 0, false
	}

	num := 0
	for c := startCol; c < endCol; c++ {
		num = num*10 + int(line[c]-'0')
	}
	start = t.OffsetOf(text.Point{Line: point.Line, Col: startCol})
	end = t.OffsetOf(text.Point{Line: point.Line, Col: endCol})
	return start, end, num, true
}

func (e *Editor) changeNumberAtCursor(
	frame *frame.Frame,
	win *window.Window,
	buf *buffer.Buffer,
	delta int,
) error {
	if buf.ReadOnly {
		return buffer.ErrReadOnly
	}
	t := buf.Text
	at := windowCursor(win)
	start, end, value, ok := scanNumberAt(t, at)
	if !ok {
		frame.Dirty = true
		return nil
	}
	value += delta
	if value < 0 {
		value = 0
	}
	repl := []byte(itoa(value))
	if _, err := buf.Replace(start, end, repl); err != nil {
		return err
	}
	setWindowCursor(buf, win, start+text.Offset(len(repl)))
	frame.Dirty = true
	e.recordNumberChange(delta)
	return nil
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var digits [20]byte
	i := len(digits)
	for n > 0 {
		i--
		digits[i] = byte('0' + n%10)
		n /= 10
	}
	return string(digits[i:])
}
