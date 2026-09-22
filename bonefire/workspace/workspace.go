package workspace

import (
	"fmt"
	"os"
	"path/filepath"
	"vague/bonefire/buffer"
	"vague/bonefire/editor"
	frame "vague/bonefire/frame"
	"vague/bonefire/text"
	"vague/bonefire/window"
)

// Workspace holds the editor core plus all frames and windows.
type Workspace struct {
	Editor *editor.Editor

	Frames  map[uint64]*frame.Frame
	Windows map[uint64]*window.Window

	CurrentFrameID uint64
	nextFrameID    uint64
	nextWindowID   uint64
}

func New() *Workspace {
	return &Workspace{
		Editor:       editor.NewEditor(),
		Frames:       make(map[uint64]*frame.Frame),
		Windows:      make(map[uint64]*window.Window),
		nextFrameID:  1,
		nextWindowID: 1,
	}
}

func errNotFound(kind string, id uint64) error {
	return fmt.Errorf("%s %d not found", kind, id)
}

func (w *Workspace) FrameContext(frameID uint64) (*frame.Frame, *window.Window, *buffer.Buffer, error) {
	frame, ok := w.Frames[frameID]
	if !ok {
		return nil, nil, nil, fmt.Errorf("frame %d not found", frameID)
	}

	win, ok := w.Windows[frame.ActiveWindowID]
	if !ok {
		return nil, nil, nil, fmt.Errorf("window %d not found", frame.ActiveWindowID)
	}

	buf, ok := w.Editor.Buffers[win.BufferId]
	if !ok {
		return nil, nil, nil, fmt.Errorf("buffer %d not found", win.BufferId)
	}

	return frame, win, buf, nil
}

func initialWorkDir(clientDir string, initialPaths ...string) string {
	if len(initialPaths) > 0 && initialPaths[0] != "" {
		if abs, err := resolvePathAgainst(defaultWorkDir(clientDir), initialPaths[0]); err == nil {
			return filepath.Dir(abs)
		}
	}

	return defaultWorkDir(clientDir)
}

func defaultWorkDir(clientDir string) string {
	if clientDir != "" {
		if abs, err := filepath.Abs(clientDir); err == nil && abs != "/" {
			return abs
		}
	}

	if home, err := os.UserHomeDir(); err == nil {
		return home
	}

	return "/"
}

func resolvePathAgainst(baseDir, path string) (string, error) {
	if path == "" {
		return "", nil
	}

	if filepath.IsAbs(path) {
		return filepath.Clean(path), nil
	}

	if baseDir == "" {
		baseDir = defaultWorkDir("")
	}

	return filepath.Abs(filepath.Join(baseDir, path))
}

func (w *Workspace) resolvePath(frameID uint64, path string) (string, error) {
	frame, ok := w.Frames[frameID]
	if !ok {
		return "", errNotFound("frame", frameID)
	}

	return resolvePathAgainst(frame.WorkDir, path)
}

func newWindow(id, frameID, bufferID uint64) *window.Window {
	win := window.NewWindow(id, frameID)
	win.BufferId = bufferID
	win.DesiredCol = 0
	return win
}

// NewFrame creates a display surface showing the current buffer.
func (w *Workspace) NewFrame(width, height int, workDir string, initialPath ...string) (*frame.Frame, error) {
	if width <= 0 {
		width = frame.DefaultWidth
	}
	if height <= 0 {
		height = frame.DefaultHeight
	}

	frameWorkDir := initialWorkDir(workDir, initialPath...)

	var (
		buf *buffer.Buffer
		err error
	)

	if len(initialPath) > 0 && initialPath[0] != "" {
		buf, err = w.Editor.LoadBufferPathAt(frameWorkDir, initialPath[0])
		if err != nil {
			return nil, err
		}
	} else if buf = w.Editor.CurrentBuffer(); buf == nil {
		buf = w.Editor.Scratch("*scratch*")
	}

	frame := &frame.Frame{
		ID:      w.nextFrameID,
		Width:   width,
		Height:  height,
		WorkDir: frameWorkDir,
		Dirty:   true,
	}
	w.nextFrameID++

	win := newWindow(w.nextWindowID, frame.ID, buf.ID)
	win.Cursor = buf.Text.AddMarker(0, text.GravityRight)
	w.nextWindowID++
	w.Windows[win.Id] = win

	frame.Root = win.LeafNode()
	frame.ActiveWindowID = win.Id

	w.Frames[frame.ID] = frame
	if w.CurrentFrameID == 0 {
		w.CurrentFrameID = frame.ID
	}

	return frame, nil
}

func (w *Workspace) ResizeFrame(frameID uint64, width, height int) error {
	frame, ok := w.Frames[frameID]
	if !ok {
		return errNotFound("frame", frameID)
	}

	if width > 0 {
		frame.Width = width
	}
	if height > 0 {
		frame.Height = height
	}

	frame.Dirty = true
	return nil
}

func (w *Workspace) SetWindowWrap(frameID uint64, wrap bool) error {
	frame, ok := w.Frames[frameID]
	if !ok {
		return errNotFound("frame", frameID)
	}

	win, ok := w.Windows[frame.ActiveWindowID]
	if !ok {
		return errNotFound("window", frame.ActiveWindowID)
	}

	win.WindowOptions.Wrap = wrap
	frame.Dirty = true
	return nil
}

func (w *Workspace) HandleInput(frameID uint64, keys string) error {
	w.ClearEcho(frameID)
	frame, win, buf, err := w.FrameContext(frameID)
	if err != nil {
		return err
	}
	return w.Editor.HandleInput(frame, win, buf, keys)
}

func (w *Workspace) Search(frameID uint64, pattern string, forward bool) error {
	frame, win, buf, err := w.FrameContext(frameID)
	if err != nil {
		return err
	}

	if err := w.Editor.Search(frame, win, buf, pattern, forward); err != nil {
		return err
	}

	w.ClearEcho(frameID)
	return nil
}

func (w *Workspace) SwitchToNextBuffer(frameID uint64) (*buffer.Buffer, error) {
	_, _, buf, err := w.FrameContext(frameID)
	if err != nil {
		return nil, err
	}

	next := w.Editor.NextBuffer(buf.ID)
	if next == nil {
		return nil, editor.ErrBufferNotFound
	}

	return w.switchBuffer(frameID, next)
}

func (w *Workspace) SwitchToPrevBuffer(frameID uint64) (*buffer.Buffer, error) {
	_, _, buf, err := w.FrameContext(frameID)
	if err != nil {
		return nil, err
	}

	prev := w.Editor.PrevBuffer(buf.ID)
	if prev == nil {
		return nil, editor.ErrBufferNotFound
	}

	return w.switchBuffer(frameID, prev)
}

func (w *Workspace) SwitchToBuffer(frameID uint64, spec string) (*buffer.Buffer, error) {
	target, err := w.Editor.ResolveBuffer(spec)
	if err != nil {
		return nil, err
	}

	return w.switchBuffer(frameID, target)
}

func (w *Workspace) BufferListMessage(frameID uint64) (string, error) {
	_, _, buf, err := w.FrameContext(frameID)
	if err != nil {
		return "", err
	}

	return w.Editor.FormatBufferList(buf.ID), nil
}

func (w *Workspace) OpenFile(frameID uint64, path string, force bool) (*buffer.Buffer, error) {
	abs, err := w.resolvePath(frameID, path)
	if err != nil {
		return nil, err
	}

	if existing := w.Editor.FindBuffer(abs); existing != nil {
		if force {
			if err := existing.Reload(); err != nil {
				return nil, err
			}
		}

		return w.switchBuffer(frameID, existing)
	}

	buf, err := w.Editor.LoadBufferPathAt("", abs)
	if err != nil {
		return nil, err
	}

	return w.switchBuffer(frameID, buf)
}

func (w *Workspace) WriteFile(frameID uint64, path string, force bool) error {
	_, _, buf, err := w.FrameContext(frameID)
	if err != nil {
		return err
	}

	if path != "" {
		abs, err := w.resolvePath(frameID, path)
		if err != nil {
			return err
		}

		if force {
			if err := buf.SaveAsForce(abs); err != nil {
				return editor.WriteFileError(path, err)
			}
			return nil
		}
		if err := buf.SaveAs(abs); err != nil {
			return editor.WriteFileError(path, err)
		}
		return nil
	}

	if force {
		if err := buf.SaveForce(); err != nil {
			return editor.WriteFileError(path, err)
		}
		return nil
	}

	if err := buf.Save(); err != nil {
		return editor.WriteFileError(path, err)
	}
	return nil
}

func (w *Workspace) switchBuffer(frameID uint64, buf *buffer.Buffer) (*buffer.Buffer, error) {
	frame, ok := w.Frames[frameID]
	if !ok {
		return nil, errNotFound("frame", frameID)
	}

	win, ok := w.Windows[frame.ActiveWindowID]
	if !ok {
		return nil, errNotFound("window", frame.ActiveWindowID)
	}

	win.BufferId = buf.ID
	win.TopLine = 0
	win.DesiredCol = 0
	win.Cursor = buf.Text.AddMarker(0, text.GravityRight)

	w.Editor.SetCurrentBuffer(buf.ID)
	w.Editor.ClearSearchMatch()
	w.Editor.SetMode(editor.NormalMode)
	frame.Dirty = true

	return buf, nil
}
