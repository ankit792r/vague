package editor

import (
	"fmt"
	"sort"
	"strings"

	"vague/backbone/process"
	"vague/bonefire/buffer"
	frame "vague/bonefire/frame"
	"vague/bonefire/text"
	"vague/bonefire/window"
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
	if to <= from || to > t.Len() {
		return nil
	}

	return offsetRangeInViewport(from, to, topLine, meta, lines, t)
}

func searchHighlightsInViewport(
	ed *Editor,
	buf *buffer.Buffer,
	topLine int,
	meta []visualLine,
	lines []string,
	t *text.Text,
) []process.Selection {
	if !ed.searchOpts.HlSearch || ed.nohlSearch || ed.searchPattern == "" || ed.searchMatchBuf != buf.ID {
		return nil
	}

	pp, err := ed.compiledPattern()
	if err != nil {
		return nil
	}

	all := findAllPattern(t.Bytes(), pp)
	if len(all) == 0 {
		return nil
	}

	curFrom := int(ed.searchMatchFrom)
	curTo := int(ed.searchMatchTo)
	out := make([]process.Selection, 0, len(all))

	for _, m := range all {
		if m[0] == curFrom && m[1] == curTo {
			continue
		}
		sel := offsetRangeInViewport(text.Offset(m[0]), text.Offset(m[1]), topLine, meta, lines, t)
		if sel == nil || !sel.Visible {
			continue
		}
		out = append(out, *sel)
	}

	return out
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

func (e *Editor) exNohlSearch(frame *frame.Frame) {
	e.ClearNohlSearch(frame)
}

func (e *Editor) exVimgrep(args string) (string, error) {
	args = strings.TrimSpace(args)
	if args == "" {
		return "", fmt.Errorf("vimgrep: pattern required")
	}
	return "", fmt.Errorf("vimgrep: quickfix list not implemented (Phase 12); use / and :g for now")
}

func (e *Editor) exTag(args string) (string, error) {
	args = strings.TrimSpace(args)
	if args == "" {
		return "", fmt.Errorf("tag: name required")
	}
	return "", fmt.Errorf("tag %q: tags file search not implemented", args)
}

func (e *Editor) exSort(
	frame *frame.Frame,
	win *window.Window,
	buf *buffer.Buffer,
	args string,
	rng exRange,
) error {
	if buf.ReadOnly {
		return buffer.ErrReadOnly
	}

	t := buf.Text
	startLine, endLine := e.resolveExRange(rng, t, win)
	if startLine > endLine {
		return nil
	}
	type lineChunk struct {
		start text.Offset
		end   text.Offset
		text  string
	}
	chunks := make([]lineChunk, 0, endLine-startLine+1)
	for line := startLine; line <= endLine; line++ {
		ls := t.LineStart(line)
		le := t.LineEnd(line)
		chunks = append(chunks, lineChunk{
			start: ls,
			end:   le,
			text:  string(t.Slice(ls, le)),
		})
	}

	sort.Slice(chunks, func(i, j int) bool {
		return chunks[i].text < chunks[j].text
	})

	if len(chunks) == 0 {
		return nil
	}

	rangeStart := t.LineStart(startLine)
	rangeEnd := t.LineEnd(endLine)
	if endLine+1 < t.LineCount() {
		rangeEnd = t.LineStart(endLine + 1)
	} else {
		rangeEnd = t.Len()
	}

	var assembled []byte
	for i, c := range chunks {
		if i > 0 {
			assembled = append(assembled, '\n')
		}
		assembled = append(assembled, c.text...)
	}

	if _, err := buf.Replace(rangeStart, rangeEnd, assembled); err != nil {
		return err
	}
	frame.Dirty = true
	return nil
}

func (e *Editor) exUniq(
	frame *frame.Frame,
	win *window.Window,
	buf *buffer.Buffer,
	rng exRange,
) error {
	if buf.ReadOnly {
		return buffer.ErrReadOnly
	}

	t := buf.Text
	startLine, endLine := e.resolveExRange(rng, t, win)
	var kept []string
	var prev string
	hadPrev := false
	for line := startLine; line <= endLine; line++ {
		s := string(t.Line(line))
		if hadPrev && s == prev {
			continue
		}
		kept = append(kept, s)
		prev = s
		hadPrev = true
	}

	rangeStart := t.LineStart(startLine)
	rangeEnd := t.LineEnd(endLine)
	if endLine+1 < t.LineCount() {
		rangeEnd = t.LineStart(endLine + 1)
	} else {
		rangeEnd = t.Len()
	}
	body := strings.Join(kept, "\n")
	if _, err := buf.Replace(rangeStart, rangeEnd, []byte(body)); err != nil {
		return err
	}
	frame.Dirty = true
	return nil
}
