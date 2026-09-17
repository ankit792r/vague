package editor

import "vague/bonefire/text"

// visualLine is one screen row and where it came from in the buffer.
type visualLine struct {
	BufferLine int
	StartCol   int
	EndCol     int
	Text       string
}

type viewLayout struct {
	Lines []string
	Meta  []visualLine
}

// layoutView turns buffer text into screen rows and mapping metadata.
// StartCol and EndCol are byte offsets within each buffer line.
func layoutView(t *text.Text, width, maxRows int, wrap bool) viewLayout {
	if width < 1 {
		width = defaultFrameWidth
	}

	var out viewLayout

	for bufLine := 0; bufLine < t.LineCount(); bufLine++ {
		line := t.Line(bufLine)
		if len(line) == 0 {
			out.appendVisual(visualLine{
				BufferLine: bufLine,
				Text:       "",
			}, maxRows)
			if maxRows > 0 && len(out.Lines) >= maxRows {
				return out
			}
			continue
		}

		if !wrap {
			chunk := line
			if len(chunk) > width {
				chunk = line[:width]
			}

			out.appendVisual(visualLine{
				BufferLine: bufLine,
				StartCol:   0,
				EndCol:     len(chunk),
				Text:       string(chunk),
			}, maxRows)

			if maxRows > 0 && len(out.Lines) >= maxRows {
				return out
			}
			continue
		}

		for start := 0; start < len(line); start += width {
			end := start + width
			if end > len(line) {
				end = len(line)
			}

			out.appendVisual(visualLine{
				BufferLine: bufLine,
				StartCol:   start,
				EndCol:     end,
				Text:       string(line[start:end]),
			}, maxRows)

			if maxRows > 0 && len(out.Lines) >= maxRows {
				return out
			}
		}
	}

	return out
}

func (v *viewLayout) appendVisual(row visualLine, maxRows int) {
	if maxRows > 0 && len(v.Lines) >= maxRows {
		return
	}

	v.Lines = append(v.Lines, row.Text)
	v.Meta = append(v.Meta, row)
}

func cursorScreenPos(meta []visualLine, point text.Point) (row, col int, visible bool) {
	lastIdx := -1
	for i, vl := range meta {
		if vl.BufferLine == point.Line {
			lastIdx = i
		}
	}

	for i, vl := range meta {
		if vl.BufferLine != point.Line {
			continue
		}

		if point.Col >= vl.StartCol && point.Col < vl.EndCol {
			return i, point.Col - vl.StartCol, true
		}

		if point.Col == vl.EndCol && i == lastIdx {
			return i, point.Col - vl.StartCol, true
		}

		if vl.StartCol == vl.EndCol && point.Col == 0 {
			return i, 0, true
		}
	}

	return 0, 0, false
}

func visualLines(t *text.Text, width, maxRows int, wrap bool) []string {
	return layoutView(t, width, maxRows, wrap).Lines
}
