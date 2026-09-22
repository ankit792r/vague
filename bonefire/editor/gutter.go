package editor

import (
	"strconv"

	"vague/bonefire/text"
	"vague/bonefire/window"
)

// gutterColumns is the width reserved for line numbers (digits + separator space).
func gutterColumns(lineCount int) int {
	if lineCount < 1 {
		lineCount = 1
	}
	digits := len(strconv.Itoa(lineCount))
	if digits < 2 {
		digits = 2
	}
	return digits + 1
}

func layoutContentWidth(frameWidth, lineCount int, opts window.WindowOptions) int {
	width := frameWidth
	width -= gutterWidth(lineCount, opts.Number, opts.Display.RelativeNumber, opts.Display.SignColumn)
	if width < 1 {
		width = 1
	}
	return width
}

func bufferLineNumbers(meta []visualLine) []int {
	if len(meta) == 0 {
		return nil
	}
	out := make([]int, len(meta))
	for i, vl := range meta {
		out[i] = vl.BufferLine + 1
	}
	return out
}

func lineCountForGutter(t *text.Text) int {
	if t == nil {
		return 1
	}
	return t.LineCount()
}
