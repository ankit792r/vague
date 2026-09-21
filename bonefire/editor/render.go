package editor

import (
	"vague/backbone/process"
	"vague/bonefire/buffer"
	frame "vague/bonefire/frame"
	"vague/bonefire/window"
)

// RenderRedraw builds the wire payload for a dirty frame.
func RenderRedraw(frame *frame.Frame, win *window.Window, buf *buffer.Buffer, ed *Editor) (process.Redraw, bool) {
	if frame == nil || !frame.Dirty {
		return process.Redraw{}, false
	}

	modeName := "normal"
	switch ed.Mode {
	case InsertMode:
		modeName = "insert"
	case VisualMode:
		modeName = "visual"
	case VisualLineMode:
		modeName = "visual-line"
	}

	fullView := layoutView(buf.Text, frame.Width, 0, win.WindowOptions.Wrap)
	point := windowPoint(buf, win)
	ensureCursorVisible(win, frame, fullView, point)
	view := sliceView(fullView, win.TopLine, frame.Height)
	row, col, visible := cursorViewportPos(win.TopLine, fullView.Meta, point)
	sel := selectionInViewport(ed, win, win.TopLine, fullView.Meta, view.Lines, buf.Text)
	bufLine := point.Line + 1
	bufCol := point.Col + 1

	frame.Dirty = false

	var echo *process.StatusEcho
	if frame.Echo.Message != "" {
		echo = &process.StatusEcho{
			Message: frame.Echo.Message,
			Kind:    frame.Echo.Kind,
		}
	}

	return process.Redraw{
		FrameID: frame.ID,
		Full:    true,
		Columns: frame.Width,
		Rows:    frame.Height,
		Wrap:    win.WindowOptions.Wrap,
		Buffer: process.RedrawBuffer{
			ID:       buf.ID,
			Name:     buf.Name,
			Modified: buf.Modified(),
		},
		Lines: view.Lines,
		Cursor: process.CursorPos{
			Row:     row,
			Column:  col,
			Visible: visible,
		},
		Selection: sel,
		Mode:      modeName,
		Position: process.BufferPosition{
			Line:   bufLine,
			Column: bufCol,
		},
		Echo: echo,
	}, true
}
