package editor

import (
	"fmt"
	"strconv"
	"strings"

	frame "vague/bonefire/frame"
	"vague/bonefire/window"
)

// GlobalDisplayOpts are editor-wide :set options (colorscheme, tabline, etc.).
type GlobalDisplayOpts struct {
	TermGuiColors bool
	ColorScheme   string
	ConcealLevel  int
	FoldEnable    bool
	TabLineFmt    string
}

func defaultGlobalDisplayOpts() GlobalDisplayOpts {
	return GlobalDisplayOpts{ColorScheme: "vague"}
}

func (e *Editor) DisplayOptionQuery(name string) string {
	switch name {
	case "termguicolors", "tgc":
		return boolOpt(e.globalDisplay.TermGuiColors)
	case "colorscheme", "cs", "background", "bg":
		return e.globalDisplay.ColorScheme
	case "conceallevel", "cole":
		return strconv.Itoa(e.globalDisplay.ConcealLevel)
	case "foldenable", "fen":
		return boolOpt(e.globalDisplay.FoldEnable)
	default:
		return ""
	}
}

func splitSetNameValue(lower string) (name, value string) {
	if i := strings.Index(lower, "="); i >= 0 {
		return lower[:i], lower[i+1:]
	}
	return lower, ""
}

func (e *Editor) ApplySetDisplayArgs(
	frame *frame.Frame,
	win *window.Window,
	args string,
	enable bool,
) (string, error) {
	lower := strings.ToLower(strings.TrimSpace(args))
	if strings.HasPrefix(lower, "no") && len(lower) > 2 && !strings.Contains(lower, "=") {
		lower = lower[2:]
		enable = false
	}
	name, val := splitSetNameValue(lower)
	if name == "" {
		return "", fmt.Errorf("Unknown option")
	}
	msg, err := e.ApplySetDisplayOption(frame, win, name, enable, val)
	if err != nil {
		return "", err
	}
	return msg, nil
}

func (e *Editor) displayOptionValue(win *window.Window, name string) string {
	if win == nil {
		return e.DisplayOptionQuery(name)
	}
	d := win.WindowOptions.Display
	switch name {
	case "scrolloff", "so":
		return strconv.Itoa(win.WindowOptions.ScrollOff)
	case "sidescrolloff", "sis":
		return strconv.Itoa(d.SideScrollOff)
	case "relativenumber", "rnu":
		return boolOpt(d.RelativeNumber)
	case "cursorline", "cul":
		return boolOpt(d.CursorLine)
	case "list":
		return boolOpt(d.List)
	case "textwidth", "tw":
		return strconv.Itoa(d.TextWidth)
	case "colorcolumn", "cc":
		if len(d.ColorColumn) == 0 {
			return ""
		}
		parts := make([]string, len(d.ColorColumn))
		for i, c := range d.ColorColumn {
			parts[i] = strconv.Itoa(c)
		}
		return strings.Join(parts, ",")
	case "statusline", "stl":
		return d.StatusLine
	default:
		q := e.DisplayOptionQuery(name)
		if q != "" && q != "setlocal" {
			return q
		}
	}
	return ""
}

func (e *Editor) ApplySetDisplayOption(
	frame *frame.Frame,
	win *window.Window,
	name string,
	enable bool,
	rawValue string,
) (string, error) {
	if win == nil {
		return "", fmt.Errorf("no window")
	}
	d := &win.WindowOptions.Display
	switch name {
	case "relativenumber", "rnu":
		d.RelativeNumber = enable
	case "cursorline", "cul":
		d.CursorLine = enable
	case "cursorcolumn", "cuc":
		d.CursorColumn = enable
	case "list":
		d.List = enable
	case "signcolumn":
		d.SignColumn = enable || rawValue == "yes" || rawValue == "auto"
	case "showmode", "smd":
		d.ShowMode = enable
	case "showcmd", "sc":
		d.ShowCmd = enable
	case "ruler":
		d.Ruler = enable
	case "smoothscroll", "sm":
		d.SmoothScroll = enable
	case "scrolloff", "so":
		if rawValue != "" {
			n, err := strconv.Atoi(rawValue)
			if err != nil {
				return "", err
			}
			win.WindowOptions.ScrollOff = n
		} else if enable {
			if win.WindowOptions.ScrollOff == 0 {
				win.WindowOptions.ScrollOff = 2
			}
		} else {
			win.WindowOptions.ScrollOff = 0
		}
	case "sidescrolloff", "sis":
		if rawValue != "" {
			n, err := strconv.Atoi(rawValue)
			if err != nil {
				return "", err
			}
			d.SideScrollOff = n
		} else if enable {
			d.SideScrollOff = 8
		} else {
			d.SideScrollOff = 0
		}
	case "textwidth", "tw":
		if rawValue == "" {
			return "", fmt.Errorf("textwidth requires value")
		}
		n, err := strconv.Atoi(rawValue)
		if err != nil {
			return "", err
		}
		d.TextWidth = n
	case "wrapmargin", "wm":
		if rawValue == "" {
			return "", fmt.Errorf("wrapmargin requires value")
		}
		n, err := strconv.Atoi(rawValue)
		if err != nil {
			return "", err
		}
		d.WrapMargin = n
	case "colorcolumn", "cc":
		if rawValue != "" {
			d.ColorColumn = parseColorColumns(rawValue)
		} else if !enable {
			d.ColorColumn = nil
		}
	case "statusline", "stl":
		if rawValue != "" {
			d.StatusLine = rawValue
		}
	case "listchars", "lcs":
		if rawValue != "" {
			d.ListChars = parseListChars(rawValue)
		}
	case "termguicolors", "tgc":
		e.globalDisplay.TermGuiColors = enable
	case "background", "bg", "colorscheme", "cs":
		if rawValue == "" {
			return "", fmt.Errorf("%s requires value", name)
		}
		e.globalDisplay.ColorScheme = rawValue
	case "conceallevel", "cole":
		if rawValue != "" {
			n, err := strconv.Atoi(rawValue)
			if err != nil {
				return "", err
			}
			e.globalDisplay.ConcealLevel = n
		} else if !enable {
			e.globalDisplay.ConcealLevel = 0
		} else {
			return "", fmt.Errorf("conceallevel needs syntax support (stub)")
		}
	case "foldenable", "fen":
		if enable {
			return "", fmt.Errorf("foldenable: see Phase 11")
		}
		e.globalDisplay.FoldEnable = false
	case "tabline", "tal":
		if rawValue != "" {
			e.globalDisplay.TabLineFmt = rawValue
		}
	default:
		return "", fmt.Errorf("Unknown option: %s", name)
	}
	frame.Dirty = true
	return name, nil
}

func parseColorColumns(s string) []int {
	parts := strings.Split(s, ",")
	var out []int
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		n, err := strconv.Atoi(p)
		if err == nil && n > 0 {
			out = append(out, n)
		}
	}
	return out
}

func parseListChars(s string) window.ListChars {
	lc := window.DefaultListChars()
	for _, part := range strings.Split(s, ",") {
		kv := strings.SplitN(part, ":", 2)
		if len(kv) != 2 {
			continue
		}
		switch strings.TrimSpace(kv[0]) {
		case "tab":
			lc.Tab = kv[1]
		case "trail":
			lc.Trail = kv[1]
		case "eol":
			lc.Eol = kv[1]
		case "extends":
			lc.Extends = kv[1]
		}
	}
	return lc
}

func (e *Editor) ShowCmdString() string {
	if !e.pendingRegQuote && e.pendingOp == opNone && e.pendingCount == 0 && e.pendingKey == "" {
		return ""
	}
	var b strings.Builder
	if e.pendingCount > 0 {
		b.WriteString(strconv.Itoa(e.pendingCount))
	}
	switch e.pendingOp {
	case opDelete:
		b.WriteByte('d')
	case opChange:
		b.WriteByte('c')
	case opYank:
		b.WriteByte('y')
	}
	if e.pendingKey != "" {
		b.WriteString(e.pendingKey)
	}
	if e.pendingRegQuote {
		b.WriteByte('"')
	}
	return b.String()
}

func FormatStatusLine(tmpl, fileName, mod string, line, col int, mode string) string {
	if tmpl == "" {
		tmpl = "%f %m %l,%c"
	}
	r := strings.NewReplacer(
		"%f", fileName,
		"%m", mod,
		"%l", strconv.Itoa(line),
		"%c", strconv.Itoa(col),
		"%M", mode,
	)
	return r.Replace(tmpl)
}

func FormatTabLine(tmpl, name string, index int, active bool) string {
	if tmpl == "" {
		tmpl = "%N"
	}
	activeMark := " "
	if active {
		activeMark = ">"
	}
	r := strings.NewReplacer(
		"%N", name,
		"%n", strconv.Itoa(index),
		"%A", activeMark,
	)
	return r.Replace(tmpl)
}

func (e *Editor) ThemeName() string {
	if e.globalDisplay.ColorScheme == "" {
		return "vague"
	}
	return e.globalDisplay.ColorScheme
}

func (e *Editor) TabLineFormat() string {
	if e.globalDisplay.TabLineFmt == "" {
		return "%N"
	}
	return e.globalDisplay.TabLineFmt
}
