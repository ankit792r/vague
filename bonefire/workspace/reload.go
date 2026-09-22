package workspace

import (
	"vague/bonefire/buffer"
	"vague/bonefire/editor"
	"vague/bonefire/text"
)

// ReloadCurrentBuffer reloads the window's buffer from disk.
func (w *Workspace) ReloadCurrentBuffer(frameID uint64, force bool) (*buffer.Buffer, error) {
	frame, win, buf, err := w.FrameContext(frameID)
	if err != nil {
		return nil, err
	}

	if buf.Path == "" {
		return nil, buffer.ErrNoFileName
	}

	if buf.Modified() && !force {
		return nil, editor.ErrNotSaved
	}

	if err := buf.Reload(); err != nil {
		return nil, err
	}

	win.TopLine = 0
	win.DesiredCol = 0
	win.Cursor = buf.Text.AddMarker(0, text.GravityRight)

	w.Editor.ClearSearchMatch()
	w.Editor.ResetInputState()
	w.Editor.SetMode(editor.NormalMode)
	frame.Dirty = true

	return buf, nil
}
