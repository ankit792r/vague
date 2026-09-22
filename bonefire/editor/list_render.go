package editor

import (
	"strings"

	"vague/bonefire/window"
)

func applyListDisplay(lines []string, lc window.ListChars) []string {
	if len(lines) == 0 {
		return lines
	}
	out := make([]string, len(lines))
	for i, line := range lines {
		out[i] = displayListLine(line, lc)
	}
	return out
}

func displayListLine(line string, lc window.ListChars) string {
	var b strings.Builder
	runes := []rune(line)
	for i, r := range runes {
		if r == '\t' {
			b.WriteString(lc.Tab)
			continue
		}
		if r == ' ' {
			trail := true
			for _, rest := range runes[i+1:] {
				if rest != ' ' {
					trail = false
					break
				}
			}
			if trail && lc.Trail != "" {
				b.WriteString(lc.Trail)
			} else {
				b.WriteRune(' ')
			}
			continue
		}
		b.WriteRune(r)
	}
	if lc.Eol != "" {
		b.WriteString(lc.Eol)
	}
	return b.String()
}

func relativeOrAbsoluteLineNumbers(meta []visualLine, cursorBufLine int, number, relative bool) []int {
	if len(meta) == 0 {
		return nil
	}
	out := make([]int, len(meta))
	for i, vl := range meta {
		if relative {
			diff := vl.BufferLine - cursorBufLine
			if diff == 0 && number {
				out[i] = vl.BufferLine + 1
			} else if diff == 0 {
				out[i] = 0
			} else if diff < 0 {
				out[i] = -diff
			} else {
				out[i] = diff
			}
		} else {
			out[i] = vl.BufferLine + 1
		}
	}
	return out
}

func gutterWidth(lineCount int, number, relative, sign bool) int {
	cols := 0
	if number || relative {
		cols += gutterColumns(lineCount)
	}
	if sign {
		cols += 2
	}
	return cols
}
