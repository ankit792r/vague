package editor

import (
	"bytes"

	"vague/bonefire/buffer"
	frame "vague/bonefire/frame"
	"vague/bonefire/text"
	"vague/bonefire/window"
)

func (e *Editor) resolveExRange(
	rng exRange,
	t *text.Text,
	win *window.Window,
) (startLine, endLine int) {
	lineCount := t.LineCount()
	cur := t.PointOf(windowCursor(win)).Line

	switch {
	case rng.wholeBuf:
		return 0, lineCount - 1
	case rng.startLine == -2 && rng.endLine == -2:
		return e.visualExLineRange(t, win)
	case rng.startLine < 0 && rng.endLine < 0:
		line := cur
		if rng.startLine == -4 {
			line = lineCount - 1
		}
		return line, line
	case rng.startLine >= 0 && rng.endLine >= 0:
		a, b := rng.startLine, rng.endLine
		if a > b {
			a, b = b, a
		}
		if a < 0 {
			a = 0
		}
		if b >= lineCount {
			b = lineCount - 1
		}
		return a, b
	default:
		return cur, cur
	}
}

func (e *Editor) visualExLineRange(t *text.Text, win *window.Window) (int, int) {
	if e.lastVisualMode == NormalMode && e.lastVisualAnchor == 0 && e.lastVisualHead == 0 {
		cur := t.PointOf(windowCursor(win)).Line
		return cur, cur
	}
	a := t.PointOf(e.lastVisualAnchor).Line
	b := t.PointOf(e.lastVisualHead).Line
	if a > b {
		a, b = b, a
	}
	return a, b
}

func (e *Editor) exSubstitute(
	frame *frame.Frame,
	win *window.Window,
	buf *buffer.Buffer,
	sub *exSubstitute,
	rng exRange,
) error {
	if buf.ReadOnly {
		return buffer.ErrReadOnly
	}

	pp, err := e.parseExPattern(sub.pattern, sub.flags)
	if err != nil {
		return err
	}

	t := buf.Text
	startLine, endLine := e.resolveExRange(rng, t, win)
	repl := []byte(sub.replacement)

	for line := startLine; line <= endLine; line++ {
		lineStart := t.LineStart(line)
		lineEnd := t.LineEnd(line)
		lineBytes := t.Slice(lineStart, lineEnd)

		changed := false
		var out []byte
		searchFrom := 0

		for searchFrom <= len(lineBytes) {
			idx, matchEnd, ok := findPatternInSlice(lineBytes[searchFrom:], pp)
			if !ok {
				out = append(out, lineBytes[searchFrom:]...)
				break
			}
			idx += searchFrom
			matchEnd += searchFrom
			out = append(out, lineBytes[searchFrom:idx]...)
			out = append(out, repl...)
			changed = true
			searchFrom = matchEnd
			if !sub.flags.global {
				out = append(out, lineBytes[searchFrom:]...)
				break
			}
		}

		if !changed {
			continue
		}

		if _, err := buf.Replace(lineStart, lineEnd, out); err != nil {
			return err
		}
	}

	frame.Dirty = true
	return nil
}

func findSubMatch(hay, needle []byte, flags subFlags) int {
	if len(needle) == 0 {
		return -1
	}
	if flags.ignoreCase && !flags.noIgnoreCase {
		hayLower := bytesToLower(bytes.Clone(hay))
		needleLower := bytesToLower(bytes.Clone(needle))
		return bytes.Index(hayLower, needleLower)
	}
	return bytes.Index(hay, needle)
}

func (e *Editor) lineMatchesPattern(line []byte, pattern string) bool {
	pp, err := e.parseSearchPattern(pattern)
	if err != nil {
		return findSubMatch(line, []byte(pattern), subFlags{ignoreCase: true}) >= 0
	}
	_, _, ok := findPatternInSlice(line, pp)
	return ok
}
