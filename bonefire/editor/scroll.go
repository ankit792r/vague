package editor

import (
	"vague/bonefire/text"
	"vague/bonefire/window"
)

func ensureCursorVisible(win *window.Window, frame *Frame, view viewLayout, point text.Point) {
	cursorRow, _, ok := visualRowAt(view.Meta, point)
	if !ok {
		return
	}

	height := frame.Height
	if height < 1 {
		height = defaultFrameHeight
	}

	total := len(view.Meta)
	maxTop := total - height
	if maxTop < 0 {
		maxTop = 0
	}

	top := win.TopLine
	if top > maxTop {
		top = maxTop
	}

	scrollOff := win.WindowOptions.ScrollOff

	if cursorRow < top+scrollOff {
		top = cursorRow - scrollOff
	} else if cursorRow >= top+height-scrollOff {
		top = cursorRow - height + scrollOff + 1
	}

	if top < 0 {
		top = 0
	}
	if top > maxTop {
		top = maxTop
	}

	win.TopLine = top
}

func sliceView(view viewLayout, topLine, height int) viewLayout {
	if topLine < 0 {
		topLine = 0
	}
	if topLine >= len(view.Meta) {
		return viewLayout{}
	}

	end := topLine + height
	if end > len(view.Meta) {
		end = len(view.Meta)
	}

	return viewLayout{
		Lines: view.Lines[topLine:end],
		Meta:  view.Meta[topLine:end],
	}
}

func cursorViewportPos(topLine int, meta []visualLine, point text.Point) (row, col int, visible bool) {
	cursorRow, col, ok := visualRowAt(meta, point)
	if !ok {
		return 0, 0, false
	}

	row = cursorRow - topLine
	if row < 0 {
		return 0, 0, false
	}

	return row, col, true
}
