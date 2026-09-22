package editor

import (
	"vague/bonefire/buffer"
	frame "vague/bonefire/frame"
)

const (
	EchoInfo  = "info"
	EchoError = "error"
)

func setFrameInfoEcho(f *frame.Frame, message string) {
	f.Echo = frame.StatusEcho{Message: message, Kind: EchoInfo}
}

func WriteEchoMessage(buf *buffer.Buffer) string {
	name := buf.Name
	if buf.Path != "" {
		name = buf.Path
	}

	return `"` + name + `" written`
}
