package window

import "vague/bonefire/text"

// visualLine is one screen row and where it came from in the buffer.
type VisualLine struct {
	BufferLine int
	StartCol   int
	EndCol     int
	Text       string
}

type ViewLayout struct {
	Lines []string
	Meta  []VisualLine
}

// layoutView turns buffer text into screen rows and mapping metadata.
// StartCol and EndCol are byte offsets within each buffer line.
func LayoutView(t *text.Text, width, maxRows int, wrap bool) ViewLayout {
	// if width < 1 {
	// 	width = frame.DefaultWidth
	// }

	var out ViewLayout

	for bufLine := 0; bufLine < t.LineCount(); bufLine++ {
		line := t.Line(bufLine)
		if len(line) == 0 {
			out.appendVisual(VisualLine{
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

			out.appendVisual(VisualLine{
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
			end := min(start+width, len(line))

			out.appendVisual(VisualLine{
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

func (v *ViewLayout) appendVisual(row VisualLine, maxRows int) {
	if maxRows > 0 && len(v.Lines) >= maxRows {
		return
	}

	v.Lines = append(v.Lines, row.Text)
	v.Meta = append(v.Meta, row)
}
