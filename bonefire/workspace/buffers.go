package workspace

import (
	"vague/bonefire/buffer"
)

func (w *Workspace) loadBufferPath(frameID uint64, path string) (*buffer.Buffer, error) {
	abs, err := w.resolvePath(frameID, path)
	if err != nil {
		return nil, err
	}
	if existing := w.Editor.FindBuffer(abs); existing != nil {
		return existing, nil
	}
	return w.Editor.LoadBufferPathAt("", abs)
}
