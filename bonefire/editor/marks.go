package editor

import (
	"fmt"
	"sort"
	"strings"
	"vague/bonefire/buffer"
	frame "vague/bonefire/frame"
	"vague/bonefire/shada"
	"vague/bonefire/text"
	"vague/bonefire/window"
)

type markJumpMode int

const (
	markJumpNone markJumpMode = iota
	markJumpExact
	markJumpLine
)

func (e *Editor) initMarks() {
	if e.localMarks == nil {
		e.localMarks = make(map[uint64]map[string]*text.Marker)
	}
	if e.changeMarks == nil {
		e.changeMarks = make(map[uint64]*text.Marker)
	}
	if e.jumpMarks == nil {
		e.jumpMarks = make(map[uint64]*text.Marker)
	}
	if e.fileMarks == nil {
		e.fileMarks = make(map[string]shada.FileMark)
	}
	if e.shada == nil {
		store, _ := shada.Load()
		e.shada = store
		for name, fm := range store.FileMarks {
			e.fileMarks[name] = fm
		}
	}
}

func (e *Editor) localMarkTable(bufID uint64) map[string]*text.Marker {
	e.initMarks()
	if e.localMarks[bufID] == nil {
		e.localMarks[bufID] = make(map[string]*text.Marker)
	}
	return e.localMarks[bufID]
}

func isLocalMarkName(r rune) bool {
	return r >= 'a' && r <= 'z'
}

func isFileMarkName(r rune) bool {
	return r >= 'A' && r <= 'Z'
}

func (e *Editor) setMarkFromNormal(buf *buffer.Buffer, win *window.Window, name rune) {
	e.initMarks()
	at := windowCursor(win)
	if isLocalMarkName(name) {
		key := string(name)
		table := e.localMarkTable(buf.ID)
		if old := table[key]; old != nil {
			buf.Text.RemoveMarker(old)
		}
		table[key] = buf.Text.AddMarker(at, text.GravityRight)
		return
	}
	if isFileMarkName(name) {
		pt := buf.Text.PointOf(at)
		key := string(name)
		fm := shada.FileMark{
			Path: buf.Path,
			Line: pt.Line + 1,
			Col:  pt.Col + 1,
		}
		if fm.Path == "" {
			fm.Path = buf.Name
		}
		e.fileMarks[key] = fm
		e.shada.FileMarks[key] = fm
		_ = e.shada.Save()
	}
}

func (e *Editor) deleteMarkName(name string) bool {
	e.initMarks()
	if len(name) != 1 {
		return false
	}
	r := rune(name[0])
	if isLocalMarkName(r) {
		for bufID, table := range e.localMarks {
			if m, ok := table[name]; ok {
				if buf := e.Buffers[bufID]; buf != nil && m != nil {
					buf.Text.RemoveMarker(m)
				}
				delete(table, name)
				return true
			}
		}
		return false
	}
	if isFileMarkName(r) {
		delete(e.fileMarks, name)
		delete(e.shada.FileMarks, name)
		_ = e.shada.Save()
		return true
	}
	return false
}

func (e *Editor) noteChangeAt(buf *buffer.Buffer, off text.Offset) {
	if buf == nil {
		return
	}
	e.initMarks()
	if old := e.changeMarks[buf.ID]; old != nil {
		buf.Text.RemoveMarker(old)
	}
	e.changeMarks[buf.ID] = buf.Text.AddMarker(off, text.GravityRight)
	e.recordChangeEntry(buf, off)
}

func (e *Editor) noteJumpFrom(buf *buffer.Buffer, off text.Offset) {
	if buf == nil {
		return
	}
	e.initMarks()
	if old := e.jumpMarks[buf.ID]; old != nil {
		buf.Text.RemoveMarker(old)
	}
	e.jumpMarks[buf.ID] = buf.Text.AddMarker(off, text.GravityRight)
}

func (e *Editor) jumpToMarkKey(
	frame *frame.Frame,
	win *window.Window,
	buf *buffer.Buffer,
	key string,
	mode markJumpMode,
	motionFrom text.Offset,
) error {
	if len(key) != 1 {
		return nil
	}
	ch := rune(key[0])
	var targetBuf *buffer.Buffer
	var off text.Offset
	var ok bool

	switch {
	case ch == '.':
		targetBuf = buf
		off, ok = e.specialMarkOffset(buf, e.changeMarks)
	case ch == '"':
		targetBuf = buf
		off, ok = e.specialMarkOffset(buf, e.jumpMarks)
	case isLocalMarkName(ch):
		targetBuf = buf
		off, ok = e.localMarkOffset(buf, key)
	case isFileMarkName(ch):
		targetBuf, off, ok = e.fileMarkOffset(ch)
	default:
		return nil
	}
	if !ok {
		return nil
	}

	if targetBuf != buf {
		e.SetCurrentBuffer(targetBuf.ID)
		buf = targetBuf
	}

	if mode == markJumpLine {
		pt := buf.Text.PointOf(off)
		off = firstNonBlankOnLine(buf.Text, pt.Line)
	}

	setWindowCursor(buf, win, off)
	e.finishMotionJump(buf, win, motionFrom)
	view := layoutViewForWindow(buf.Text, win, frame)
	rememberColumn(buf, win, view)
	frame.Dirty = true
	return nil
}

func (e *Editor) specialMarkOffset(buf *buffer.Buffer, table map[uint64]*text.Marker) (text.Offset, bool) {
	e.initMarks()
	m := table[buf.ID]
	if m == nil {
		return 0, false
	}
	return m.Off, true
}

func (e *Editor) localMarkOffset(buf *buffer.Buffer, name string) (text.Offset, bool) {
	e.initMarks()
	m := e.localMarkTable(buf.ID)[name]
	if m == nil {
		return 0, false
	}
	return m.Off, true
}

func (e *Editor) fileMarkOffset(name rune) (*buffer.Buffer, text.Offset, bool) {
	e.initMarks()
	fm, ok := e.fileMarks[string(name)]
	if !ok {
		return nil, 0, false
	}
	var buf *buffer.Buffer
	if fm.Path != "" {
		buf = e.FindBuffer(fm.Path)
	}
	if buf == nil {
		for _, b := range e.Buffers {
			if b.Name == fm.Path || b.Path == fm.Path {
				buf = b
				break
			}
		}
	}
	if buf == nil {
		return nil, 0, false
	}
	line := fm.Line - 1
	if line < 0 {
		line = 0
	}
	if line >= buf.Text.LineCount() {
		line = buf.Text.LineCount() - 1
	}
	col := fm.Col - 1
	if col < 0 {
		col = 0
	}
	off := buf.Text.OffsetOf(text.Point{Line: line, Col: col})
	return buf, buf.Text.Clamp(off), true
}

func (e *Editor) marksOnLineLabel(buf *buffer.Buffer, line int) string {
	if buf == nil {
		return ""
	}
	e.initMarks()
	var names []string
	for name, m := range e.localMarkTable(buf.ID) {
		if m == nil {
			continue
		}
		if buf.Text.PointOf(m.Off).Line == line {
			names = append(names, name)
		}
	}
	if len(names) == 0 {
		return ""
	}
	sort.Strings(names)
	return strings.Join(names, "")
}

func (e *Editor) formatMarksListing() string {
	e.initMarks()
	type row struct {
		sortKey string
		line    string
	}
	var rows []row

	for bufID, table := range e.localMarks {
		buf := e.Buffers[bufID]
		if buf == nil {
			continue
		}
		for name, m := range table {
			if m == nil {
				continue
			}
			pt := buf.Text.PointOf(m.Off)
			label := bufLabel(buf)
			rows = append(rows, row{
				sortKey: fmt.Sprintf("%s %s", label, name),
				line:    fmt.Sprintf("%s  %s  %d  %d", name, label, pt.Line+1, pt.Col+1),
			})
		}
	}
	for name, fm := range e.fileMarks {
		rows = append(rows, row{
			sortKey: fmt.Sprintf("%s %s", fm.Path, name),
			line:    fmt.Sprintf("%s  %s  %d  %d", name, fm.Path, fm.Line, fm.Col),
		})
	}
	if len(rows) == 0 {
		return "No marks set"
	}
	sort.Slice(rows, func(i, j int) bool { return rows[i].sortKey < rows[j].sortKey })
	var b strings.Builder
	for i, r := range rows {
		if i > 0 {
			b.WriteByte('\n')
		}
		b.WriteString(r.line)
	}
	return b.String()
}

func bufLabel(buf *buffer.Buffer) string {
	if buf.Path != "" {
		return buf.Path
	}
	return buf.Name
}

func parseMarkJumpKey(keys string) (string, bool) {
	if keys == "." || keys == `"` {
		return keys, true
	}
	if len(keys) != 1 {
		return "", false
	}
	r := rune(keys[0])
	if isLocalMarkName(r) || isFileMarkName(r) {
		return keys, true
	}
	return "", false
}

func isMarkSetKey(keys string) bool {
	if len(keys) != 1 {
		return false
	}
	r := rune(keys[0])
	return isLocalMarkName(r) || isFileMarkName(r)
}
