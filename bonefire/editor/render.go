package editor

import (
	"vague/backbone/process"
	"vague/bonefire/buffer"
	"vague/bonefire/display"
	"vague/bonefire/window"
)

// RenderRedraw builds the wire payload for a dirty frame.
func RenderRedraw(frame *display.Frame, win *window.Window, buf *buffer.Buffer, mode Mode) (process.Redraw, bool) {
	if frame == nil || !frame.Dirty {
		return process.Redraw{}, false
	}

	modeName := "normal"
	if mode == InsertMode {
		modeName = "insert"
	}

	fullView := layoutView(buf.Text, frame.Width, 0, win.WindowOptions.Wrap)
	point := windowPoint(buf, win)
	ensureCursorVisible(win, frame, fullView, point)
	view := sliceView(fullView, win.TopLine, frame.Height)
	row, col, visible := cursorViewportPos(win.TopLine, fullView.Meta, point)

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
		Mode: modeName,
		Echo: echo,
	}, true
}
