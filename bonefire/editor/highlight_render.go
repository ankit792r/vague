package editor

import (
	"vague/backbone/process"
	"vague/bonefire/text"
)

func offsetRangeInViewport(
	from, to text.Offset,
	topLine int,
	meta []visualLine,
	lines []string,
	t *text.Text,
) *process.Selection {
	if to <= from {
		return nil
	}

	startPt := t.PointOf(from)
	endPt := t.PointOf(to - 1)

	sr, sc, ok1 := visualRowAt(meta, startPt)
	er, ec, ok2 := visualRowAt(meta, endPt)
	if !ok1 || !ok2 {
		return nil
	}

	sr -= topLine
	er -= topLine

	if sr >= len(lines) && er >= len(lines) {
		return nil
	}
	if sr < 0 {
		sr = 0
		sc = 0
	}
	if er >= len(lines) {
		er = len(lines) - 1
		if er >= 0 {
			lineLen := len(lines[er])
			if lineLen == 0 {
				ec = 0
			} else if ec >= lineLen {
				ec = lineLen - 1
			}
		}
	}

	if sr > er || (sr == er && sc > ec) {
		sr, sc, er, ec = er, ec, sr, sc
	}

	return &process.Selection{
		Visible: true,
		Start:   process.SelectionPoint{Row: sr, Column: sc},
		End:     process.SelectionPoint{Row: er, Column: ec},
	}
}
