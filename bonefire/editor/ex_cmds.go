package editor

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"vague/bonefire/buffer"
	frame "vague/bonefire/frame"
	"vague/bonefire/window"
)

func (e *Editor) exGlobal(
	frame *frame.Frame,
	win *window.Window,
	buf *buffer.Buffer,
	cmd exCommand,
	rng exRange,
) error {
	invert := cmd.kind == "vglobal"
	rest := cmd.raw
	if len(rest) < 2 {
		return fmt.Errorf("Invalid global command")
	}
	delim := rest[1]
	pat, exCmd, ok := splitGlobalPattern(rest[2:], delim)
	if !ok {
		return fmt.Errorf("Invalid global pattern")
	}

	t := buf.Text
	startLine, endLine := e.resolveExRange(rng, t, win)
	patBytes := []byte(pat)

	for line := startLine; line <= endLine; line++ {
		lineBytes := t.Line(line)
		if !bytesContainsMatch(lineBytes, patBytes, subFlags{ignoreCase: true}) {
			if !invert {
				continue
			}
		} else if invert {
			continue
		}

		setWindowCursor(buf, win, t.LineStart(line))
		if _, err := e.RunExLine(frame, win, buf, exCmd); err != nil {
			return err
		}
	}
	frame.Dirty = true
	return nil
}

func splitGlobalPattern(s string, delim byte) (pat, cmd string, ok bool) {
	for i := 0; i < len(s); i++ {
		if s[i] == delim && (i == 0 || s[i-1] != '\\') {
			return s[:i], strings.TrimSpace(s[i+1:]), true
		}
	}
	return "", "", false
}

func bytesContainsMatch(hay, needle []byte, flags subFlags) bool {
	return findSubMatch(hay, needle, flags) >= 0
}

func (e *Editor) exNormal(
	frame *frame.Frame,
	win *window.Window,
	buf *buffer.Buffer,
	keys string,
	rng exRange,
) error {
	if strings.TrimSpace(keys) == "" {
		return fmt.Errorf("normal: keys required")
	}
	t := buf.Text
	startLine, endLine := e.resolveExRange(rng, t, win)
	savedMode := e.Mode
	e.Mode = NormalMode

	for line := startLine; line <= endLine; line++ {
		setWindowCursor(buf, win, t.LineStart(line))
		if err := e.feedNormalKeys(frame, win, buf, keys); err != nil {
			e.Mode = savedMode
			return err
		}
	}
	e.Mode = savedMode
	frame.Dirty = true
	return nil
}

func (e *Editor) feedNormalKeys(
	frame *frame.Frame,
	win *window.Window,
	buf *buffer.Buffer,
	keys string,
) error {
	for _, r := range keys {
		if err := e.HandleInput(frame, win, buf, string(r)); err != nil {
			return err
		}
	}
	return nil
}

func (e *Editor) exRead(
	frame *frame.Frame,
	win *window.Window,
	buf *buffer.Buffer,
	path string,
	rng exRange,
) error {
	path = strings.TrimSpace(path)
	if path == "" {
		return fmt.Errorf("read: file name required")
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	at := windowCursor(win)
	if rng.startLine >= 0 {
		t := buf.Text
		at = t.LineStart(rng.startLine)
	}
	if len(data) > 0 && data[len(data)-1] != '\n' {
		data = append(data, '\n')
	}
	if _, err := buf.Insert(at, data); err != nil {
		return err
	}
	frame.Dirty = true
	return nil
}

func (e *Editor) exWrite(
	frame *frame.Frame,
	win *window.Window,
	buf *buffer.Buffer,
	args string,
	rng exRange,
	bang, appendMode bool,
) (string, error) {
	args = strings.TrimSpace(args)
	if appendMode || strings.HasPrefix(args, ">>") {
		path := strings.TrimSpace(strings.TrimPrefix(args, ">>"))
		if path == "" && buf.Path != "" {
			path = buf.Path
		}
		if path == "" {
			return "", buffer.ErrNoFileName
		}
		f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return "", err
		}
		defer f.Close()
		start, end := e.resolveExRange(rng, buf.Text, win)
		var data []byte
		for line := start; line <= end; line++ {
			data = append(data, buf.Text.Line(line)...)
			data = append(data, '\n')
		}
		if _, err := f.Write(data); err != nil {
			return "", err
		}
		return fmt.Sprintf(`"%s" appended`, filepath.Base(path)), nil
	}
	_ = bang
	_ = rng
	return "", fmt.Errorf("use :write for full buffer save")
}


func (e *Editor) exBufferDelete(
	frame *frame.Frame,
	win *window.Window,
	buf *buffer.Buffer,
	bang bool,
) (string, error) {
	if buf.Modified() && !bang {
		return "", fmt.Errorf("No write since last change (add ! to override)")
	}
	return fmt.Sprintf(`buffer "%s" deleted (stub)`, buf.Name), nil
}

func (e *Editor) exCd(f *frame.Frame, path string) (string, error) {
	path = strings.TrimSpace(path)
	if path == "" {
		return f.WorkDir, nil
	}
	abs, err := filepath.Abs(filepath.Join(f.WorkDir, path))
	if err != nil {
		return "", err
	}
	f.WorkDir = abs
	return abs, nil
}

func (e *Editor) exSet(
	frame *frame.Frame,
	win *window.Window,
	args string,
) (string, error) {
	args = strings.TrimSpace(args)
	if args == "" {
		return "", fmt.Errorf("set: option required")
	}
	lower := strings.ToLower(args)
	if strings.HasSuffix(lower, "?") {
		key := strings.TrimSpace(args[:len(args)-1])
		return fmt.Sprintf("%s?", key), nil
	}
	if strings.Contains(lower, "+=") {
		return "", fmt.Errorf("set += not implemented")
	}
	if strings.HasSuffix(lower, "&") {
		return "option reset", nil
	}
	switch lower {
	case "wrap", "nowrap", "number", "nu", "nonumber", "nonu":
		return lower, nil
	default:
		return "", fmt.Errorf("Unknown option: %s", args)
	}
}

func (e *Editor) exMap(kind, args string) (string, error) {
	e.ensureMaps()
	args = strings.TrimSpace(args)
	if args == "" {
		var b strings.Builder
		b.WriteString("Maps:\n")
		for k, v := range e.userMaps {
			fmt.Fprintf(&b, "  %s -> %s\n", k, v)
		}
		return strings.TrimSpace(b.String()), nil
	}
	parts := strings.Fields(args)
	if len(parts) < 2 {
		return "", fmt.Errorf("map: lhs and rhs required")
	}
	key := kind + ":" + parts[0]
	e.userMaps[key] = strings.Join(parts[1:], " ")
	return "map defined", nil
}

func (e *Editor) ensureMaps() {
	if e.userMaps == nil {
		e.userMaps = make(map[string]string)
	}
}

func (e *Editor) exUserCommand(args string) (string, error) {
	return "", fmt.Errorf("user :command definitions not implemented")
}

func (e *Editor) exSource(
	frame *frame.Frame,
	win *window.Window,
	buf *buffer.Buffer,
	path string,
) error {
	path = strings.TrimSpace(path)
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "\"") {
			continue
		}
		if _, err := e.RunExLine(frame, win, buf, line); err != nil {
			return err
		}
	}
	return nil
}

func (e *Editor) exReg() string {
	data, _, ok := e.registerText()
	if !ok {
		return "Type Name Content\n\" (empty)"
	}
	return fmt.Sprintf("Type Name Content\n\" %q", string(data))
}

func (e *Editor) exMarks() string {
	if len(e.marks) == 0 {
		return "No marks set"
	}
	var b strings.Builder
	for name, off := range e.marks {
		fmt.Fprintf(&b, "%s %d\n", name, off)
	}
	return strings.TrimSpace(b.String())
}

func (e *Editor) exDelMark(args string) string {
	args = strings.TrimSpace(args)
	if args == "" {
		return "delm: mark required"
	}
	delete(e.marks, args)
	return "mark deleted"
}

func (e *Editor) exJumps() string {
	if len(e.jumps) == 0 {
		return "jump list empty"
	}
	return fmt.Sprintf("%d jumps (use Ctrl-o / Ctrl-i)", len(e.jumps))
}

func (e *Editor) clearJumps() {
	e.jumps = nil
	e.jumpPos = 0
}
