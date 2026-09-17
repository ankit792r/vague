package editor

import "vague/backbone/process"

func (e *Editor) InvalidateFrame(frameID uint64) {
	if frame, ok := e.Frames[frameID]; ok {
		frame.dirty = true
	}
}

func (e *Editor) RenderRedraw(frameID uint64) (process.Redraw, bool) {
	frame, ok := e.Frames[frameID]
	if !ok || !frame.dirty {
		return process.Redraw{}, false
	}

	win, ok := e.Windows[frame.ActiveWindowID]
	if !ok {
		return process.Redraw{}, false
	}

	buf, ok := e.Buffers[win.BufferId]
	if !ok {
		return process.Redraw{}, false
	}

	mode := "normal"
	if e.Mode == InsertMode {
		mode = "insert"
	}

	fullView := layoutView(buf.Text, frame.Width, 0, win.WindowOptions.Wrap)
	point := windowPoint(buf, win)
	ensureCursorVisible(win, frame, fullView, point)
	view := sliceView(fullView, win.TopLine, frame.Height)
	row, col, visible := cursorViewportPos(win.TopLine, fullView.Meta, point)

	frame.dirty = false

	return process.Redraw{
		FrameID: frameID,
		Full:    true,
		Columns: frame.Width,
		Rows:    frame.Height,
		Wrap:    win.WindowOptions.Wrap,
		Buffer: process.RedrawBuffer{
			ID:   buf.ID,
			Name: buf.Name,
		},
		Lines: view.Lines,
		Cursor: process.CursorPos{
			Row:     row,
			Column:  col,
			Visible: visible,
		},
		Mode: mode,
	}, true
}
