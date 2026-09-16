package editor

import "vague/bonefire/buffer"

type Mode int

const (
	NormalMode Mode = iota
	InsertMode
)

type Editor struct {
	Mode Mode

	Buffers       map[uint64]*buffer.Buffer
	CurrentBuffer uint64
	nextBufferID  uint64
}

func NewEditor() *Editor {
	return &Editor{
		Mode:         NormalMode,
		nextBufferID: 1,
	}
}

// Scratch creates a buffer with no backing file.
func (e *Editor) Scratch(name string) *buffer.Buffer {
	buffer := buffer.NewScratch(e.nextBufferID, name)

	e.nextBufferID++
	e.Buffers[buffer.ID] = buffer
	e.CurrentBuffer = buffer.ID

	return buffer
}
