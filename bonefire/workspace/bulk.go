package workspace

import (
	"vague/bonefire/editor"
)

func (w *Workspace) WriteAll(frameID uint64, force bool) error {
	for _, buf := range w.Editor.Buffers {
		if !buf.Modified() || buf.Path == "" {
			continue
		}
		var err error
		if force {
			err = buf.SaveForce()
		} else {
			err = buf.Save()
		}
		if err != nil {
			return editor.WriteFileError(buf.Path, err)
		}
	}
	return nil
}

func (w *Workspace) QuitAll(force bool) error {
	ids := make([]uint64, 0, len(w.Frames))
	for id := range w.Frames {
		ids = append(ids, id)
	}
	for _, id := range ids {
		if err := w.QuitFrame(id, force); err != nil {
			return err
		}
	}
	return nil
}

func (w *Workspace) WQAll(frameID uint64, force bool) error {
	if err := w.WriteAll(frameID, force); err != nil {
		return err
	}
	return w.QuitAll(force)
}

func (w *Workspace) XAll(frameID uint64, force bool) error {
	if err := w.WriteAll(frameID, force); err != nil {
		return err
	}
	return w.QuitAll(force)
}
