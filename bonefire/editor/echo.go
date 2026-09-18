package editor

import "vague/bonefire/buffer"

const (
	EchoInfo  = "info"
	EchoError = "error"
)

// StatusEcho is user-visible feedback shown on the status line until the next key.
type StatusEcho struct {
	Message string
	Kind    string
}

func (e *Editor) SetEcho(frameID uint64, message, kind string) error {
	frame, ok := e.Frames[frameID]
	if !ok {
		return errNotFound("frame", frameID)
	}

	if message == "" {
		frame.echo = StatusEcho{}
	} else {
		if kind == "" {
			kind = EchoInfo
		}
		frame.echo = StatusEcho{Message: message, Kind: kind}
	}

	frame.dirty = true
	return nil
}

func (e *Editor) ClearEcho(frameID uint64) {
	frame, ok := e.Frames[frameID]
	if !ok || frame.echo.Message == "" {
		return
	}

	frame.echo = StatusEcho{}
	frame.dirty = true
}

func WriteEchoMessage(buf *buffer.Buffer) string {
	name := buf.Name
	if buf.Path != "" {
		name = buf.Path
	}

	return `"` + name + `" written`
}
