package workspace

import (
	"vague/backbone/process"
	"vague/bonefire/editor"
	"vague/bonefire/window"
)

func (w *Workspace) InvalidateFrame(frameID uint64) {
	if frame, ok := w.Frames[frameID]; ok {
		frame.Dirty = true
	}
}

func (w *Workspace) RenderRedraw(frameID uint64) (process.Redraw, bool) {
	frame, ok := w.Frames[frameID]
	if !ok || !frame.Dirty {
		return process.Redraw{}, false
	}

	paneRects := window.LayoutPanes(frame.Root, frame.Width, frame.Height)
	if len(paneRects) == 0 {
		return process.Redraw{}, false
	}

	var panes []process.RedrawPane
	var activePane process.RedrawPane
	for _, rect := range paneRects {
		win, ok := w.Windows[rect.WindowID]
		if !ok {
			continue
		}
		buf, ok := w.Editor.Buffers[win.BufferId]
		if !ok {
			continue
		}
		active := rect.WindowID == frame.ActiveWindowID
		pane := editor.RenderPane(frame, win, buf, w.Editor, rect, active)
		panes = append(panes, pane)
		if active {
			activePane = pane
		}
	}

	var tabs []process.RedrawTab
	for i, tab := range frame.Tabs {
		tabs = append(tabs, process.RedrawTab{
			Label:  tab.Label,
			Active: i == frame.ActiveTab,
		})
	}

	var echo *process.StatusEcho
	if frame.Echo.Message != "" {
		echo = &process.StatusEcho{
			Message: frame.Echo.Message,
			Kind:    frame.Echo.Kind,
		}
	}

	frame.Dirty = false

	return process.Redraw{
		FrameID:          frame.ID,
		Full:             true,
		Columns:          frame.Width,
		Rows:             frame.Height,
		Wrap:             activePane.Wrap,
		Number:           activePane.Number,
		GutterColumns:    activePane.GutterColumns,
		Buffer:           activePane.Buffer,
		Lines:            activePane.Lines,
		LineNumbers:      activePane.LineNumbers,
		Cursor:           activePane.Cursor,
		Selection:        activePane.Selection,
		SearchMatch:      activePane.SearchMatch,
		SearchHighlights: activePane.SearchHighlights,
		Mode:             activePane.Mode,
		Position:         activePane.Position,
		LineMarks:        activePane.LineMarks,
		Echo:             echo,
		Panes:            panes,
		Tabs:             tabs,
		ActiveTab:        frame.ActiveTab,
	}, true
}
