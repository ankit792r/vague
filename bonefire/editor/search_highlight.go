package editor

import (
	"bytes"

	"vague/backbone/process"
	"vague/bonefire/buffer"
	"vague/bonefire/text"
)

func searchMatchInViewport(
	ed *Editor,
	buf *buffer.Buffer,
	topLine int,
	meta []visualLine,
	lines []string,
	t *text.Text,
) *process.Selection {
	if ed.searchPattern == "" || ed.searchMatchBuf != buf.ID {
		return nil
	}

	from := ed.searchMatchFrom
	to := ed.searchMatchTo
	if to <= from {
		return nil
	}

	if to > t.Len() {
		return nil
	}

	pat := []byte(ed.searchPattern)
	if len(pat) == 0 || !bytes.Equal(t.Slice(from, to), pat) {
		return nil
	}

	return offsetRangeInViewport(from, to, topLine, meta, lines, t)
}

func (e *Editor) setSearchMatch(buf *buffer.Buffer, from, to text.Offset) {
	e.searchMatchFrom = from
	e.searchMatchTo = to
	e.searchMatchBuf = buf.ID
}

func (e *Editor) ClearSearchMatch() {
	e.searchMatchFrom = 0
	e.searchMatchTo = 0
	e.searchMatchBuf = 0
}
