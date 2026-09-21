package workspace

import (
	"vague/bonefire/display"
	"vague/bonefire/editor"
)

const (
	EchoInfo  = editor.EchoInfo
	EchoError = editor.EchoError
)

func (w *Workspace) SetEcho(frameID uint64, message, kind string) error {
	frame, ok := w.Frames[frameID]
	if !ok {
		return errNotFound("frame", frameID)
	}

	if message == "" {
		frame.Echo = display.StatusEcho{}
	} else {
		if kind == "" {
			kind = EchoInfo
		}
		frame.Echo = display.StatusEcho{Message: message, Kind: kind}
	}

	frame.Dirty = true
	return nil
}

func (w *Workspace) ClearEcho(frameID uint64) {
	frame, ok := w.Frames[frameID]
	if !ok || frame.Echo.Message == "" {
		return
	}

	frame.Echo = display.StatusEcho{}
	frame.Dirty = true
}
