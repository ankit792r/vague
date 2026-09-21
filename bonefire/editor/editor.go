package editor

import (
	"vague/bonefire/buffer"
	"vague/bonefire/text"
)

type Mode int

const (
	NormalMode Mode = iota
	InsertMode
	VisualMode
	VisualLineMode
)

type Editor struct {
	Mode Mode

	Buffers        map[uint64]*buffer.Buffer
	currentBuffer  uint64
	insertGroup   uint64
	pendingKey    string
	pendingOp     opKind
	visualAnchor  text.Offset
	reg           register

	nextBufferID uint64
}

func NewEditor() *Editor {
	return &Editor{
		Mode:         NormalMode,
		Buffers:      make(map[uint64]*buffer.Buffer),
		nextBufferID: 1,
	}
}

func (e *Editor) SetMode(mode Mode) {
	e.Mode = mode
}

func (e *Editor) ResetInputState() {
	e.insertGroup = 0
	e.pendingKey = ""
	e.pendingOp = opNone
}

func (e *Editor) CurrentBuffer() *buffer.Buffer {
	if e.currentBuffer == 0 {
		return nil
	}

	return e.Buffers[e.currentBuffer]
}

func (e *Editor) SetCurrentBuffer(id uint64) {
	e.currentBuffer = id
}

// Scratch creates a buffer with no backing file.
func (e *Editor) Scratch(name string) *buffer.Buffer {
	buffer := buffer.NewScratch(e.nextBufferID, name)

	e.nextBufferID++
	e.Buffers[buffer.ID] = buffer
	e.currentBuffer = buffer.ID

	return buffer
}

func (e *Editor) FindBuffer(path string) *buffer.Buffer {
	for _, buf := range e.Buffers {
		if buf.Path == path {
			return buf
		}
	}

	return nil
}
