package editor

import (
	"fmt"
	"vague/bonefire/buffer"
	"vague/bonefire/window"
)

// HandleInput applies one key in Vim notation for the given frame.
func (e *Editor) HandleInput(frameID uint64, keys string) error {
	switch e.Mode {
	case NormalMode:
		return e.normalKey(frameID, keys)
	default:
		return nil
	}
}

func (e *Editor) normalKey(frameID uint64, keys string) error {
	frame, win, buf, err := e.frameContext(frameID)
	if err != nil {
		return err
	}

	lines := splitLogicalLines(buf.Text)
	cursor := windowPoint(win)

	switch keys {
	case "h", "<Left>", "<BS>":
		cursor = moveLeft(lines, cursor, 1)
	case "l", "<Right>", "<Space>":
		cursor = moveRight(lines, cursor, 1)
	case "j", "<Down>", "<CR>":
		cursor = moveVertical(lines, cursor, win.DesiredCol, 1)
	case "k", "<Up>":
		cursor = moveVertical(lines, cursor, win.DesiredCol, -1)
	default:
		return nil
	}

	setWindowPoint(win, cursor)
	rememberColumn(win, cursor)
	frame.dirty = true
	return nil
}

func (e *Editor) frameContext(frameID uint64) (*Frame, *window.Window, *buffer.Buffer, error) {
	frame, ok := e.Frames[frameID]
	if !ok {
		return nil, nil, nil, fmt.Errorf("frame %d not found", frameID)
	}

	win, ok := e.Windows[frame.ActiveWindowID]
	if !ok {
		return nil, nil, nil, fmt.Errorf("window %d not found", frame.ActiveWindowID)
	}

	buf, ok := e.Buffers[win.BufferId]
	if !ok {
		return nil, nil, nil, fmt.Errorf("buffer %d not found", win.BufferId)
	}

	return frame, win, buf, nil
}
