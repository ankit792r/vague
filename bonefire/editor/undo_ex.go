package editor

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"vague/bonefire/buffer"
	frame "vague/bonefire/frame"
	"vague/bonefire/text"
	"vague/bonefire/window"
)

type changeEntry struct {
	Seq  int
	Line int
	Col  int
	When time.Time
}

func (e *Editor) recordChangeEntry(buf *buffer.Buffer, off text.Offset) {
	if buf == nil {
		return
	}
	e.initChangeLog()
	pt := buf.Text.PointOf(off)
	entry := changeEntry{
		Seq:  buf.History.Seq(),
		Line: pt.Line + 1,
		Col:  pt.Col + 1,
		When: time.Now(),
	}
	log := e.changeLog[buf.ID]
	log = append(log, entry)
	const maxChanges = 100
	if len(log) > maxChanges {
		log = log[len(log)-maxChanges:]
	}
	e.changeLog[buf.ID] = log
}

func (e *Editor) initChangeLog() {
	if e.changeLog == nil {
		e.changeLog = make(map[uint64][]changeEntry)
	}
}

func (e *Editor) exUndolist(buf *buffer.Buffer) string {
	lines := buf.History.UndolistLines()
	return strings.Join(lines, "\n")
}

func (e *Editor) exChanges(buf *buffer.Buffer) string {
	e.initChangeLog()
	log := e.changeLog[buf.ID]
	if len(log) == 0 {
		return "No changes recorded"
	}
	var b strings.Builder
	start := len(log) - 20
	if start < 0 {
		start = 0
	}
	for i := start; i < len(log); i++ {
		c := log[i]
		fmt.Fprintf(&b, "%d  %d:%d  %s\n", c.Seq, c.Line, c.Col, c.When.Format("15:04:05"))
	}
	markDot := "."
	markQuote := `"`
	if off, ok := e.specialMarkOffset(buf, e.changeMarks); ok {
		pt := buf.Text.PointOf(off)
		fmt.Fprintf(&b, "change mark %s at %d:%d\n", markDot, pt.Line+1, pt.Col+1)
	}
	if off, ok := e.specialMarkOffset(buf, e.jumpMarks); ok {
		pt := buf.Text.PointOf(off)
		fmt.Fprintf(&b, "jump mark %s at %d:%d\n", markQuote, pt.Line+1, pt.Col+1)
	}
	return strings.TrimSpace(b.String())
}

func (e *Editor) exEarlier(
	frame *frame.Frame,
	win *window.Window,
	buf *buffer.Buffer,
	args string,
) error {
	steps, dur, err := parseEarlierLaterArgs(args)
	if err != nil {
		return err
	}
	if dur > 0 {
		seq, ok := buf.History.GoEarlierTime(dur)
		if !ok {
			frame.Dirty = true
			return nil
		}
		return e.gotoUndoSeq(frame, win, buf, seq)
	}
	if steps < 1 {
		steps = 1
	}
	for i := 0; i < steps; i++ {
		if err := e.undoTo(frame, win, buf, false); err != nil {
			return err
		}
	}
	return nil
}

func (e *Editor) exLater(
	frame *frame.Frame,
	win *window.Window,
	buf *buffer.Buffer,
	args string,
) error {
	steps, dur, err := parseEarlierLaterArgs(args)
	if err != nil {
		return err
	}
	if dur > 0 {
		seq, ok := buf.History.GoLaterTime(dur)
		if !ok {
			frame.Dirty = true
			return nil
		}
		return e.gotoUndoSeq(frame, win, buf, seq)
	}
	if steps < 1 {
		steps = 1
	}
	for i := 0; i < steps; i++ {
		if err := e.undoTo(frame, win, buf, true); err != nil {
			return err
		}
	}
	return nil
}

func (e *Editor) gotoUndoSeq(
	frame *frame.Frame,
	win *window.Window,
	buf *buffer.Buffer,
	seq int,
) error {
	cur, ok := buf.RestoreUndoSeq(seq)
	if !ok {
		frame.Dirty = true
		return nil
	}
	setWindowCursor(buf, win, clampToLine(buf.Text, cur))
	view := layoutViewForWindow(buf.Text, win, frame)
	rememberColumn(buf, win, view)
	frame.Dirty = true
	return nil
}

func parseEarlierLaterArgs(args string) (steps int, dur time.Duration, err error) {
	args = strings.TrimSpace(args)
	if args == "" {
		return 1, 0, nil
	}
	if d, ok := parseDurationArg(args); ok {
		return 0, d, nil
	}
	n, err := strconv.Atoi(args)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid count: %s", args)
	}
	return n, 0, nil
}

func parseDurationArg(s string) (time.Duration, bool) {
	if len(s) < 2 {
		return 0, false
	}
	unit := s[len(s)-1]
	num, err := strconv.Atoi(s[:len(s)-1])
	if err != nil || num < 1 {
		return 0, false
	}
	switch unit {
	case 's':
		return time.Duration(num) * time.Second, true
	case 'm':
		return time.Duration(num) * time.Minute, true
	case 'h':
		return time.Duration(num) * time.Hour, true
	case 'd':
		return time.Duration(num) * 24 * time.Hour, true
	default:
		return 0, false
	}
}
