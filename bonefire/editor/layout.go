package editor

import "strings"

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

// layoutView turns logical buffer lines into screen rows and mapping metadata.
func layoutView(logical []string, width, maxRows int, wrap bool) viewLayout {
	if width < 1 {
		width = defaultFrameWidth
	}

	var out viewLayout

	for bufLine, line := range logical {
		runes := []rune(line)
		if len(runes) == 0 {
			out.appendVisual(visualLine{
				BufferLine: bufLine,
				StartCol:   0,
				EndCol:     0,
				Text:       "",
			}, maxRows)
			if maxRows > 0 && len(out.Lines) >= maxRows {
				return out
			}
			continue
		}

		if !wrap && len(runes) > width {
			runes = runes[:width]
		}

		if !wrap {
			out.appendVisual(visualLine{
				BufferLine: bufLine,
				StartCol:   0,
				EndCol:     len(runes),
				Text:       string(runes),
			}, maxRows)
			if maxRows > 0 && len(out.Lines) >= maxRows {
				return out
			}
			continue
		}

		for start := 0; start < len(runes); start += width {
			end := start + width
			if end > len(runes) {
				end = len(runes)
			}

			out.appendVisual(visualLine{
				BufferLine: bufLine,
				StartCol:   start,
				EndCol:     end,
				Text:       string(runes[start:end]),
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

func cursorScreenPos(meta []visualLine, cursor Point) (row, col int, visible bool) {
	for i, vl := range meta {
		if vl.BufferLine != cursor.Line {
			continue
		}

		if cursor.Col >= vl.StartCol && cursor.Col < vl.EndCol {
			return i, cursor.Col - vl.StartCol, true
		}

		if vl.StartCol == vl.EndCol && cursor.Col == 0 {
			return i, 0, true
		}
	}

	return 0, 0, false
}

func splitLogicalLines(text string) []string {
	if text == "" {
		return []string{""}
	}

	return strings.Split(text, "\n")
}

// visualLines turns logical buffer lines into screen rows.
func visualLines(logical []string, width, maxRows int, wrap bool) []string {
	return layoutView(logical, width, maxRows, wrap).Lines
}
