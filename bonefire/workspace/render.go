package workspace

import (
	"vague/backbone/process"
	"vague/bonefire/editor"
)

func (w *Workspace) InvalidateFrame(frameID uint64) {
	if frame, ok := w.Frames[frameID]; ok {
		frame.Dirty = true
	}
}

func (w *Workspace) RenderRedraw(frameID uint64) (process.Redraw, bool) {
	frame, ok := w.Frames[frameID]
	if !ok || !frame.Dirty {
		return process.Redraw{}, false
	}

	win, ok := w.Windows[frame.ActiveWindowID]
	if !ok {
		return process.Redraw{}, false
	}

	buf, ok := w.Editor.Buffers[win.BufferId]
	if !ok {
		return process.Redraw{}, false
	}

	return editor.RenderRedraw(frame, win, buf, w.Editor)
}
