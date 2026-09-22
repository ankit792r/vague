package workspace

import "vague/bonefire/editor"

// QuitFrame closes a frame after checking whether its buffer may be discarded.
func (w *Workspace) QuitFrame(frameID uint64, force bool) error {
	return w.withBufferLeave(frameID, force, func() error {
		return w.CloseFrame(frameID)
	})
}

// CloseFrame removes a frame and its windows.
func (w *Workspace) CloseFrame(frameID uint64) error {
	if _, ok := w.Frames[frameID]; !ok {
		return errNotFound("frame", frameID)
	}

	for id, win := range w.Windows {
		if win.FrameId == frameID {
			delete(w.Windows, id)
		}
	}

	delete(w.Frames, frameID)

	if w.CurrentFrameID == frameID {
		w.CurrentFrameID = 0
		for id := range w.Frames {
			w.CurrentFrameID = id
			break
		}
	}

	w.Editor.ResetInputState()
	if w.CurrentFrameID == 0 {
		w.Editor.SetMode(editor.NormalMode)
	}

	return nil
}
