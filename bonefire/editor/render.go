package editor

import "vague/backbone/process"

// InvalidateFrame forces the next redraw to include buffer content.
func (e *Editor) InvalidateFrame(frameID uint64) {
	if frame, ok := e.Frames[frameID]; ok {
		frame.dirty = true
	}
}

// RenderRedraw produces the update a client needs for a frame.
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

	lines := visualLines(
		splitLogicalLines(buf.Text),
		frame.Width,
		frame.Height,
		win.WindowOptions.Wrap,
	)

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
		Lines: lines,
		Mode:  mode,
	}, true
}
