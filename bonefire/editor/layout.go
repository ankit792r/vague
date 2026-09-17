package editor

import "strings"

// visualLines turns logical buffer lines into screen rows.
//
// Wrapped rows contain only the characters that belong on that row — no
// trailing spaces to pad out to the window width. Unwrapped rows are
// truncated to width.
func visualLines(logical []string, width, maxRows int, wrap bool) []string {
	if width < 1 {
		width = defaultFrameWidth
	}

	var out []string

	for _, line := range logical {
		chunks := splitVisualLine(line, width, wrap)
		for _, chunk := range chunks {
			out = append(out, chunk)
			if maxRows > 0 && len(out) >= maxRows {
				return out
			}
		}
	}

	return out
}

func splitVisualLine(line string, width int, wrap bool) []string {
	runes := []rune(line)
	if len(runes) == 0 {
		return []string{""}
	}

	if !wrap {
		if len(runes) > width {
			return []string{string(runes[:width])}
		}

		return []string{line}
	}

	chunks := make([]string, 0, (len(runes)+width-1)/width)
	for start := 0; start < len(runes); start += width {
		end := start + width
		if end > len(runes) {
			end = len(runes)
		}

		chunks = append(chunks, string(runes[start:end]))
	}

	return chunks
}

func splitLogicalLines(text string) []string {
	if text == "" {
		return []string{""}
	}

	return strings.Split(text, "\n")
}
