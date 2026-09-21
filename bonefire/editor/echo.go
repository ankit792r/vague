package editor

import "vague/bonefire/buffer"

const (
	EchoInfo  = "info"
	EchoError = "error"
)

func WriteEchoMessage(buf *buffer.Buffer) string {
	name := buf.Name
	if buf.Path != "" {
		name = buf.Path
	}

	return `"` + name + `" written`
}
