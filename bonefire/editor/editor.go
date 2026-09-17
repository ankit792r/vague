package editor

import (
	"vague/bonefire/buffer"
	"vague/bonefire/window"
)

type Mode int

const (
	NormalMode Mode = iota
	InsertMode
)

type Editor struct {
	Mode Mode

	Buffers       map[uint64]*buffer.Buffer
	Frames        map[uint64]*Frame
	Windows       map[uint64]*window.Window
	CurrentBuffer uint64
	CurrentFrameID uint64

	nextBufferID uint64
	nextFrameID  uint64
	nextWindowID uint64
}

func NewEditor() *Editor {
	return &Editor{
		Mode:         NormalMode,
		Buffers:      make(map[uint64]*buffer.Buffer),
		Frames:       make(map[uint64]*Frame),
		Windows:      make(map[uint64]*window.Window),
		nextBufferID: 1,
		nextFrameID:  1,
		nextWindowID: 1,
	}
}

// NewFrame creates a display surface showing the current buffer.
// If nothing is open yet, a scratch buffer is created first.
func (e *Editor) NewFrame(width, height int) (*Frame, error) {
	if width <= 0 {
		width = defaultFrameWidth
	}
	if height <= 0 {
		height = defaultFrameHeight
	}

	buf := e.currentBuffer()
	if buf == nil {
		buf = e.Scratch("*scratch*")
	}

	frame := &Frame{
		ID:     e.nextFrameID,
		Width:  width,
		Height: height,
		dirty:  true,
	}
	e.nextFrameID++

	win := newWindow(e.nextWindowID, frame.ID, buf.ID)
	e.nextWindowID++
	e.Windows[win.Id] = win

	frame.Root = win.LeafNode()
	frame.ActiveWindowID = win.Id

	e.Frames[frame.ID] = frame
	if e.CurrentFrameID == 0 {
		e.CurrentFrameID = frame.ID
	}

	return frame, nil
}

func (e *Editor) currentBuffer() *buffer.Buffer {
	if e.CurrentBuffer == 0 {
		return nil
	}

	return e.Buffers[e.CurrentBuffer]
}

// Scratch creates a buffer with no backing file.
func (e *Editor) Scratch(name string) *buffer.Buffer {
	buffer := buffer.NewScratch(e.nextBufferID, name)

	e.nextBufferID++
	e.Buffers[buffer.ID] = buffer
	e.CurrentBuffer = buffer.ID

	return buffer
}

func newWindow(id, frameID, bufferID uint64) *window.Window {
	win := window.NewWindow(id, frameID)
	win.BufferId = bufferID
	return win
}
