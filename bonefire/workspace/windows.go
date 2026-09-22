package workspace

import (
	"fmt"
	"strings"

	"vague/bonefire/buffer"
	frame "vague/bonefire/frame"
	"vague/bonefire/text"
	"vague/bonefire/window"
)

func (w *Workspace) syncTabFromFrame(f *frame.Frame) {
	if len(f.Tabs) == 0 {
		return
	}
	if f.ActiveTab < 0 || f.ActiveTab >= len(f.Tabs) {
		f.ActiveTab = 0
	}
	f.Tabs[f.ActiveTab].Root = f.Root
	f.Tabs[f.ActiveTab].ActiveWindowID = f.ActiveWindowID
}

func (w *Workspace) loadTab(f *frame.Frame, index int) {
	w.syncTabFromFrame(f)
	if index < 0 || index >= len(f.Tabs) {
		return
	}
	f.ActiveTab = index
	tab := f.Tabs[index]
	f.Root = tab.Root
	f.ActiveWindowID = tab.ActiveWindowID
	f.Dirty = true
}

func (w *Workspace) registerWindow(f *frame.Frame, win *window.Window, newBuf *buffer.Buffer) {
	f.ActiveWindowID = win.Id
	w.Editor.SetCurrentBuffer(newBuf.ID)
	f.Dirty = true
}

func (w *Workspace) newWindowForFrame(f *frame.Frame, buf *buffer.Buffer, copyOpts *window.Window) *window.Window {
	win := newWindow(w.nextWindowID, f.ID, buf.ID)
	w.nextWindowID++
	if copyOpts != nil {
		win.WindowOptions = copyOpts.WindowOptions
		win.TopLine = copyOpts.TopLine
		at := windowCursor(copyOpts)
		win.Cursor = buf.Text.AddMarker(at, text.GravityRight)
	} else {
		win.Cursor = buf.Text.AddMarker(0, text.GravityRight)
	}
	w.Windows[win.Id] = win
	return win
}

func windowCursor(win *window.Window) text.Offset {
	if win == nil || win.Cursor == nil {
		return 0
	}
	return win.Cursor.Off
}

// SplitWindow opens a horizontal split (:split).
func (w *Workspace) SplitWindow(frameID uint64, path string) error {
	return w.splitFrame(frameID, false, path, false)
}

// VSplitWindow opens a vertical split (:vsplit).
func (w *Workspace) VSplitWindow(frameID uint64, path string) error {
	return w.splitFrame(frameID, true, path, false)
}

// NewWindow opens a horizontal split with a new scratch buffer (:new).
func (w *Workspace) NewWindow(frameID uint64) error {
	return w.splitFrame(frameID, false, "", true)
}

// VNewWindow opens a vertical split with a new scratch buffer (:vnew).
func (w *Workspace) VNewWindow(frameID uint64) error {
	return w.splitFrame(frameID, true, "", true)
}

func (w *Workspace) splitFrame(frameID uint64, vertical bool, path string, scratch bool) error {
	f, win, buf, err := w.FrameContext(frameID)
	if err != nil {
		return err
	}
	w.syncTabFromFrame(f)

	var newBuf *buffer.Buffer
	if scratch {
		newBuf = w.Editor.Scratch("*scratch*")
	} else if path != "" {
		abs, err := w.resolvePath(frameID, path)
		if err != nil {
			return err
		}
		if existing := w.Editor.FindBuffer(abs); existing != nil {
			newBuf = existing
		} else {
			newBuf, err = w.Editor.LoadBufferPathAt("", abs)
			if err != nil {
				return err
			}
		}
	} else {
		newBuf = buf
	}

	newWin := w.newWindowForFrame(f, newBuf, win)
	if !window.SplitLeaf(f.Root, win.Id, newWin.Id, vertical) {
		return fmt.Errorf("split failed")
	}
	w.registerWindow(f, newWin, newBuf)
	w.Editor.SetCurrentBuffer(newBuf.ID)
	return nil
}

// OnlyWindow keeps the active window and closes others.
func (w *Workspace) OnlyWindow(frameID uint64) error {
	f, _, _, err := w.FrameContext(frameID)
	if err != nil {
		return err
	}
	w.syncTabFromFrame(f)
	f.Root = window.OnlyWindow(f.Root, f.ActiveWindowID)
	f.Dirty = true
	return nil
}

// CloseWindow closes the active window.
func (w *Workspace) CloseWindow(frameID uint64) error {
	f, win, _, err := w.FrameContext(frameID)
	if err != nil {
		return err
	}
	leaves := f.Root.Leaves()
	if len(leaves) <= 1 {
		return fmt.Errorf("cannot close last window")
	}
	w.syncTabFromFrame(f)
	leaf := f.Root.Find(win.Id)
	if leaf == nil {
		return fmt.Errorf("window not in layout")
	}
	idx := 0
	for i, l := range leaves {
		if l.WindowID == win.Id {
			idx = i
			break
		}
	}
	focus := leaves[0].WindowID
	if idx+1 < len(leaves) {
		focus = leaves[idx+1].WindowID
	} else if idx > 0 {
		focus = leaves[idx-1].WindowID
	}
	f.Root = window.RemoveLeafRoot(f.Root, leaf)
	if f.Root == nil {
		return fmt.Errorf("cannot close last window")
	}
	delete(w.Windows, win.Id)
	f.ActiveWindowID = focus
	if b, ok := w.Editor.Buffers[w.Windows[focus].BufferId]; ok {
		w.Editor.SetCurrentBuffer(b.ID)
	}
	f.Dirty = true
	return nil
}

// HideWindow closes the window like :hide (same as close for now).
func (w *Workspace) HideWindow(frameID uint64) error {
	return w.CloseWindow(frameID)
}

// FocusWindow switches the active window.
func (w *Workspace) FocusWindow(frameID uint64, windowID uint64) error {
	f, ok := w.Frames[frameID]
	if !ok {
		return errNotFound("frame", frameID)
	}
	if f.Root.Find(windowID) == nil {
		return fmt.Errorf("window not found")
	}
	f.ActiveWindowID = windowID
	if win, ok := w.Windows[windowID]; ok {
		if buf, ok := w.Editor.Buffers[win.BufferId]; ok {
			w.Editor.SetCurrentBuffer(buf.ID)
		}
	}
	f.Dirty = true
	return nil
}

// WinCmd runs a :wincmd letter.
func (w *Workspace) WinCmd(frameID uint64, cmd string) error {
	cmd = strings.TrimSpace(cmd)
	if cmd == "" {
		return fmt.Errorf("wincmd: command required")
	}
	f, _, _, err := w.FrameContext(frameID)
	if err != nil {
		return err
	}
	switch cmd[0] {
	case 'h', 'j', 'k', 'l':
		id := window.NeighborWindow(f.Root, f.ActiveWindowID, rune(cmd[0]))
		if id == 0 {
			return nil
		}
		return w.FocusWindow(frameID, id)
	case 'H', 'L', 'J', 'K':
		leaves := f.Root.Leaves()
		if len(leaves) == 0 {
			return nil
		}
		target := leaves[0].WindowID
		if cmd[0] == 'L' || cmd[0] == 'J' {
			target = leaves[len(leaves)-1].WindowID
		}
		return w.FocusWindow(frameID, target)
	case 'w', 'W':
		return w.cycleWindow(frameID, cmd[0] == 'W')
	case 'c', 'o', 'q':
		return w.CloseWindow(frameID)
	case '+', '-', '<', '>':
		w.syncTabFromFrame(f)
		if window.AdjustWeight(f.Root, f.ActiveWindowID, rune(cmd[0]), 0.25) {
			f.Dirty = true
		}
		return nil
	case '=':
		window.Balance(f.Root)
		f.Dirty = true
		return nil
	case '_':
		if window.MaximizePane(f.Root, f.ActiveWindowID, false) {
			f.Dirty = true
		}
		return nil
	case '|':
		if window.MaximizePane(f.Root, f.ActiveWindowID, true) {
			f.Dirty = true
		}
		return nil
	default:
		return fmt.Errorf("wincmd: unknown %q", cmd)
	}
}

func (w *Workspace) cycleWindow(frameID uint64, reverse bool) error {
	f, _, _, err := w.FrameContext(frameID)
	if err != nil {
		return err
	}
	leaves := f.Root.Leaves()
	if len(leaves) <= 1 {
		return nil
	}
	idx := 0
	for i, l := range leaves {
		if l.WindowID == f.ActiveWindowID {
			idx = i
			break
		}
	}
	if reverse {
		idx--
		if idx < 0 {
			idx = len(leaves) - 1
		}
	} else {
		idx = (idx + 1) % len(leaves)
	}
	return w.FocusWindow(frameID, leaves[idx].WindowID)
}

// HandleCtrlWKey handles one key after <C-w> in normal mode.
func (w *Workspace) HandleCtrlWKey(frameID uint64, keys string) error {
	if keys == "" {
		return nil
	}
	if keys == "<C-w>" {
		return w.WinCmd(frameID, "w")
	}
	key := keys
	if strings.HasPrefix(keys, "<C-") && strings.HasSuffix(keys, ">") {
		key = keys[3 : len(keys)-1]
	}
	return w.WinCmd(frameID, key)
}

func (w *Workspace) SetWindowScrollOff(frameID uint64, off int) error {
	f, win, _, err := w.FrameContext(frameID)
	if err != nil {
		return err
	}
	win.WindowOptions.ScrollOff = off
	f.Dirty = true
	return nil
}

func (w *Workspace) SetLocalOption(frameID uint64, name string, enable bool) (string, error) {
	f, win, _, err := w.FrameContext(frameID)
	if err != nil {
		return "", err
	}
	token := name
	if !enable {
		token = "no" + name
	}
	if msg, err := w.Editor.ApplySetDisplayArgs(f, win, token, enable); err == nil {
		return msg, nil
	}
	switch name {
	case "wrap", "linebreak":
		win.WindowOptions.Wrap = enable
	case "nowrap":
		win.WindowOptions.Wrap = false
	case "number", "nu":
		win.WindowOptions.Number = enable
	case "nonumber", "nonu":
		win.WindowOptions.Number = false
	default:
		return "", fmt.Errorf("setlocal: unknown option %s", name)
	}
	f.Dirty = true
	return name, nil
}
