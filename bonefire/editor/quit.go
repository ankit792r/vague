package editor

import "errors"

// ErrNotSaved means quit was refused because the buffer has unsaved changes.
var ErrNotSaved = errors.New("No write since last change (add ! to override)")

// QuitFrame closes a frame after checking whether its buffer may be discarded.
func (e *Editor) QuitFrame(frameID uint64, force bool) error {
	_, _, buf, err := e.frameContext(frameID)
	if err != nil {
		return err
	}

	if buf.Modified() && !force {
		return ErrNotSaved
	}

	return e.CloseFrame(frameID)
}

// CloseFrame removes a frame and its windows from the editor.
func (e *Editor) CloseFrame(frameID uint64) error {
	if _, ok := e.Frames[frameID]; !ok {
		return errNotFound("frame", frameID)
	}

	for id, win := range e.Windows {
		if win.FrameId == frameID {
			delete(e.Windows, id)
		}
	}

	delete(e.Frames, frameID)

	if e.CurrentFrameID == frameID {
		e.CurrentFrameID = 0
		for id := range e.Frames {
			e.CurrentFrameID = id
			break
		}
	}

	e.insertGroup = 0
	e.pendingKey = ""
	if e.CurrentFrameID == 0 {
		e.Mode = NormalMode
	}

	return nil
}
