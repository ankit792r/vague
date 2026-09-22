package workspace

import (
	"fmt"

	frame "vague/bonefire/frame"
	"vague/bonefire/window"
)

func (w *Workspace) TabNew(frameID uint64, path string) error {
	f, _, _, err := w.FrameContext(frameID)
	if err != nil {
		return err
	}
	w.syncTabFromFrame(f)

	label := fmt.Sprintf("Tab %d", len(f.Tabs)+1)
	var root *window.Node
	var active uint64

	if path != "" {
		buf, err := w.loadBufferPath(frameID, path)
		if err != nil {
			return err
		}
		win := w.newWindowForFrame(f, buf, nil)
		root = win.LeafNode()
		active = win.Id
		label = buf.Name
	} else if buf := w.Editor.CurrentBuffer(); buf != nil {
		win := w.newWindowForFrame(f, buf, nil)
		root = win.LeafNode()
		active = win.Id
		label = buf.Name
	} else {
		buf := w.Editor.Scratch("*scratch*")
		win := w.newWindowForFrame(f, buf, nil)
		root = win.LeafNode()
		active = win.Id
	}

	f.Tabs = append(f.Tabs, frame.TabPage{
		Root:           root,
		ActiveWindowID: active,
		Label:          label,
	})
	w.loadTab(f, len(f.Tabs)-1)
	return nil
}

func (w *Workspace) TabClose(frameID uint64) error {
	f, _, _, err := w.FrameContext(frameID)
	if err != nil {
		return err
	}
	if len(f.Tabs) <= 1 {
		return fmt.Errorf("cannot close last tab")
	}
	w.syncTabFromFrame(f)
	idx := f.ActiveTab
	f.Tabs = append(f.Tabs[:idx], f.Tabs[idx+1:]...)
	if f.ActiveTab >= len(f.Tabs) {
		f.ActiveTab = len(f.Tabs) - 1
	}
	w.loadTab(f, f.ActiveTab)
	return nil
}

func (w *Workspace) TabNext(frameID uint64) error {
	f, _, _, err := w.FrameContext(frameID)
	if err != nil {
		return err
	}
	if len(f.Tabs) <= 1 {
		return nil
	}
	w.syncTabFromFrame(f)
	next := (f.ActiveTab + 1) % len(f.Tabs)
	w.loadTab(f, next)
	return nil
}

func (w *Workspace) TabPrev(frameID uint64) error {
	f, _, _, err := w.FrameContext(frameID)
	if err != nil {
		return err
	}
	if len(f.Tabs) <= 1 {
		return nil
	}
	w.syncTabFromFrame(f)
	prev := f.ActiveTab - 1
	if prev < 0 {
		prev = len(f.Tabs) - 1
	}
	w.loadTab(f, prev)
	return nil
}
