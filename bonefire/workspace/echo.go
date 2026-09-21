package workspace

import (
	"vague/bonefire/editor"
	frame "vague/bonefire/frame"
)

const (
	EchoInfo  = editor.EchoInfo
	EchoError = editor.EchoError
)

func (w *Workspace) SetEcho(frameID uint64, message, kind string) error {
	f, ok := w.Frames[frameID]
	if !ok {
		return errNotFound("frame", frameID)
	}

	if message == "" {
		f.Echo = frame.StatusEcho{}
	} else {
		if kind == "" {
			kind = EchoInfo
		}
		f.Echo = frame.StatusEcho{Message: message, Kind: kind}
	}

	f.Dirty = true
	return nil
}

func (w *Workspace) ClearEcho(frameID uint64) {
	f, ok := w.Frames[frameID]
	if !ok || f.Echo.Message == "" {
		return
	}

	f.Echo = frame.StatusEcho{}
	f.Dirty = true
}
