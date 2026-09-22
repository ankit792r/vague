package editor

import (
	"vague/backbone/process"
	"vague/bonefire/buffer"
	frame "vague/bonefire/frame"
	"vague/bonefire/text"
	"vague/bonefire/window"
)

// RenderPane builds one window's redraw region.
func RenderPane(
	f *frame.Frame,
	win *window.Window,
	buf *buffer.Buffer,
	ed *Editor,
	rect window.PaneRect,
	active bool,
) process.RedrawPane {
	modeName := modeNameFor(ed)

	lineCount := lineCountForGutter(buf.Text)
	contentWidth := layoutContentWidth(rect.Width, lineCount, win.WindowOptions)
	fullView := layoutView(buf.Text, contentWidth, 0, win.WindowOptions.Wrap)
	point := windowPoint(buf, win)
	ensureCursorVisibleHeight(win, rect.Height, fullView, point)
	view := sliceView(fullView, win.TopLine, rect.Height)
	row, col, visible := cursorViewportPos(win.TopLine, fullView.Meta, point)
	sel := selectionInViewport(ed, win, win.TopLine, fullView.Meta, view.Lines, buf.Text)
	searchMatch := searchMatchInViewport(ed, buf, win.TopLine, fullView.Meta, view.Lines, buf.Text)
	searchHighlights := searchHighlightsInViewport(ed, buf, win.TopLine, fullView.Meta, view.Lines, buf.Text)

	var lineNumbers []int
	gutterCols := gutterWidth(lineCount, win.WindowOptions.Number, win.WindowOptions.Display.RelativeNumber, win.WindowOptions.Display.SignColumn)
	opts := win.WindowOptions
	if opts.Number || opts.Display.RelativeNumber {
		lineNumbers = relativeOrAbsoluteLineNumbers(view.Meta, point.Line, opts.Number, opts.Display.RelativeNumber)
	}
	lines := view.Lines
	if opts.Display.List {
		lines = applyListDisplay(lines, opts.Display.ListChars)
	}
	showCmd := ""
	if opts.Display.ShowCmd {
		showCmd = ed.ShowCmdString()
	}
	mod := " "
	if buf.Modified() {
		mod = "+"
	}
	statusText := FormatStatusLine(opts.Display.StatusLine, buf.Name, mod, point.Line+1, point.Col+1, modeName)
	showMode := opts.Display.ShowMode

	return process.RedrawPane{
		WindowID:      win.Id,
		X:             rect.X,
		Y:             rect.Y,
		Columns:       rect.Width,
		Rows:          rect.Height,
		Active:        active,
		Wrap:          win.WindowOptions.Wrap,
		Number:        win.WindowOptions.Number,
		GutterColumns: gutterCols,
		RelativeNumber: opts.Display.RelativeNumber,
		List:          opts.Display.List,
		CursorLineRow: cursorLineRow(row, opts.Display.CursorLine && active),
		CursorColumn:  col,
		CursorColumnOn: opts.Display.CursorColumn && visible && active,
		ColorColumns:  opts.Display.ColorColumn,
		SignColumn:    opts.Display.SignColumn,
		ShowCmd:       showCmd,
		ShowMode:      showMode,
		Ruler:         opts.Display.Ruler,
		StatusLine:    statusText,
		Theme:         ed.ThemeName(),
		Buffer: process.RedrawBuffer{
			ID:       buf.ID,
			Name:     buf.Name,
			Modified: buf.Modified(),
		},
		Lines:            lines,
		LineNumbers:      lineNumbers,
		Cursor:           process.CursorPos{Row: row, Column: col, Visible: visible},
		Selection:        sel,
		SearchMatch:      searchMatch,
		SearchHighlights: searchHighlights,
		Mode:             modeName,
		Position: process.BufferPosition{
			Line:   point.Line + 1,
			Column: point.Col + 1,
		},
		LineMarks: ed.marksOnLineLabel(buf, point.Line),
	}
}

func modeNameFor(ed *Editor) string {
	switch ed.Mode {
	case InsertMode:
		return "insert"
	case VisualMode:
		return "visual"
	case VisualLineMode:
		return "visual-line"
	case VisualBlockMode:
		return "visual-block"
	case ReplaceMode:
		return "replace"
	default:
		return "normal"
	}
}

func cursorLineRow(row int, on bool) int {
	if !on {
		return -1
	}
	return row
}

func ensureCursorVisibleHeight(win *window.Window, height int, view viewLayout, point text.Point) {
	if height < 1 {
		height = frame.DefaultHeight
	}
	cursorRow, _, ok := visualRowAt(view.Meta, point)
	if !ok {
		return
	}
	total := len(view.Meta)
	maxTop := total - height
	if maxTop < 0 {
		maxTop = 0
	}
	top := win.TopLine
	if top > maxTop {
		top = maxTop
	}
	scrollOff := win.WindowOptions.ScrollOff
	if cursorRow < top+scrollOff {
		top = cursorRow - scrollOff
	} else if cursorRow >= top+height-scrollOff {
		top = cursorRow - height + scrollOff + 1
	}
	if top < 0 {
		top = 0
	}
	if top > maxTop {
		top = maxTop
	}
	win.TopLine = top
}
