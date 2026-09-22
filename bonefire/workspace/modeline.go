package workspace

import (
	"strings"

	"vague/bonefire/buffer"
	"vague/bonefire/window"
)

func applyModeline(win *window.Window, buf *buffer.Buffer) {
	if buf == nil || win == nil || buf.Text.LineCount() == 0 {
		return
	}
	line := string(buf.Text.Line(buf.Text.LineCount() - 1))
	idx := strings.Index(line, "vim:")
	if idx < 0 {
		idx = strings.Index(line, "vi:")
	}
	if idx < 0 {
		return
	}
	after := strings.TrimSpace(line[idx:])
	if strings.HasPrefix(after, "vim:") {
		after = strings.TrimSpace(after[4:])
	} else if strings.HasPrefix(after, "vi:") {
		after = strings.TrimSpace(after[3:])
	} else {
		return
	}
	after = strings.TrimSuffix(after, ":")
	for _, tok := range strings.Fields(after) {
		tok = strings.TrimSpace(tok)
		if tok == "" {
			continue
		}
		if strings.HasPrefix(tok, "no") && len(tok) > 2 {
			applyModelineToken(win, tok[2:], false)
			continue
		}
		applyModelineToken(win, tok, true)
	}
}

func applyModelineToken(win *window.Window, tok string, enable bool) {
	switch strings.ToLower(tok) {
	case "nu", "number":
		win.WindowOptions.Number = enable
	case "wrap":
		win.WindowOptions.Wrap = enable
	case "nowrap":
		win.WindowOptions.Wrap = !enable
	}
}
