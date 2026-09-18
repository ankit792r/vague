package editor

import (
	"fmt"
	"path/filepath"

	"vague/bonefire/buffer"
	"vague/bonefire/text"
	"vague/bonefire/window"
)

// OpenFile loads path into the frame's window, reusing an existing buffer when possible.
func (e *Editor) OpenFile(frameID uint64, path string, force bool) (*buffer.Buffer, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	if existing := e.findBuffer(abs); existing != nil {
		if force {
			if err := existing.Reload(); err != nil {
				return nil, err
			}
		}

		return e.switchBuffer(frameID, existing)
	}

	buf, err := e.loadBufferPath(path)
	if err != nil {
		return nil, err
	}

	return e.switchBuffer(frameID, buf)
}

func (e *Editor) loadBufferPath(path string) (*buffer.Buffer, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	if existing := e.findBuffer(abs); existing != nil {
		e.CurrentBuffer = existing.ID
		return existing, nil
	}

	loaded, err := buffer.Load(abs)
	if err != nil {
		return nil, err
	}

	loaded.ID = e.nextBufferID
	e.nextBufferID++
	e.Buffers[loaded.ID] = loaded
	e.CurrentBuffer = loaded.ID

	return loaded, nil
}

// WriteFile saves the frame's buffer. An empty path writes to the buffer path.
func (e *Editor) WriteFile(frameID uint64, path string, force bool) error {
	_, _, buf, err := e.frameContext(frameID)
	if err != nil {
		return err
	}

	if path != "" {
		if force {
			return buf.SaveAsForce(path)
		}
		return buf.SaveAs(path)
	}

	if force {
		return buf.SaveForce()
	}

	return buf.Save()
}

func (e *Editor) findBuffer(path string) *buffer.Buffer {
	for _, buf := range e.Buffers {
		if buf.Path == path {
			return buf
		}
	}

	return nil
}

func (e *Editor) switchBuffer(frameID uint64, buf *buffer.Buffer) (*buffer.Buffer, error) {
	frame, ok := e.Frames[frameID]
	if !ok {
		return nil, errNotFound("frame", frameID)
	}

	win, ok := e.Windows[frame.ActiveWindowID]
	if !ok {
		return nil, errNotFound("window", frame.ActiveWindowID)
	}

	win.BufferId = buf.ID
	win.TopLine = 0
	win.DesiredCol = 0
	win.Cursor = buf.Text.AddMarker(0, text.GravityRight)

	e.CurrentBuffer = buf.ID
	e.Mode = NormalMode
	frame.dirty = true

	return buf, nil
}

// FrameContext returns the frame, its active window, and buffer.
func (e *Editor) FrameContext(frameID uint64) (*Frame, *window.Window, *buffer.Buffer, error) {
	return e.frameContext(frameID)
}

// BufferInfo builds the execute result for buffer commands.
func BufferInfo(buf *buffer.Buffer, current bool) map[string]any {
	return map[string]any{
		"id":       buf.ID,
		"name":     buf.Name,
		"path":     buf.Path,
		"modified": buf.Modified(),
		"readonly": buf.ReadOnly,
		"current":  current,
		"lines":    buf.Text.LineCount(),
		"bytes":    int(buf.Text.Len()),
	}
}

// OpenFileError wraps an open failure for the wire.
func OpenFileError(path string, err error) error {
	return fmt.Errorf("edit %q: %w", path, err)
}

// WriteFileError wraps a save failure for the wire.
func WriteFileError(path string, err error) error {
	if path == "" {
		return fmt.Errorf("write: %w", err)
	}

	return fmt.Errorf("write %q: %w", path, err)
}
