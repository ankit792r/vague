package editor

import (
	"regexp"
	"strconv"
	"strings"
)

type searchOffset struct {
	anchor byte // 's', 'e', 'b', or 0 → start of match
	delta  int
}

var searchOffsetNumeric = regexp.MustCompile(`^[+-]\d+$`)
var searchOffsetAnchored = regexp.MustCompile(`^([esb])([+-]?\d+)$`)

// splitSearchOffset peels a Vim-style trailing offset from a / or ? pattern.
func splitSearchOffset(raw string) (pattern string, off searchOffset) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return raw, searchOffset{}
	}

	anchoredAt := -1
	var anchoredOff searchOffset
	numericAt := -1
	var numericOff searchOffset

	for i := 0; i < len(raw); i++ {
		suffix := raw[i:]
		if i == 0 {
			continue
		}
		if off, ok := parseAnchoredOffsetSuffix(suffix); ok {
			if anchoredAt == -1 || i < anchoredAt {
				anchoredAt = i
				anchoredOff = off
			}
		}
		if (raw[i] == '+' || raw[i] == '-') && searchOffsetNumeric.MatchString(suffix) {
			n, err := strconv.Atoi(suffix)
			if err == nil {
				numericAt = i
				numericOff = searchOffset{delta: n}
			}
		}
	}

	if anchoredAt > 0 {
		return raw[:anchoredAt], anchoredOff
	}
	if numericAt > 0 {
		return raw[:numericAt], numericOff
	}
	return raw, searchOffset{}
}

func parseAnchoredOffsetSuffix(s string) (searchOffset, bool) {
	m := searchOffsetAnchored.FindStringSubmatch(s)
	if m == nil {
		return searchOffset{}, false
	}
	off := searchOffset{anchor: m[1][0]}
	if m[2] != "" {
		n, err := strconv.Atoi(m[2])
		if err != nil {
			return searchOffset{}, false
		}
		off.delta = n
	}
	return off, true
}

func applySearchOffset(start, end, bufLen int, off searchOffset) int {
	if end <= start {
		return clampSearchPos(start, bufLen)
	}

	pos := start
	switch off.anchor {
	case 'e':
		pos = end - 1
	case 's', 'b':
		pos = start
	default:
		pos = start
	}

	pos += off.delta
	return clampSearchPos(pos, bufLen)
}

func clampSearchPos(pos, bufLen int) int {
	if pos < 0 {
		return 0
	}
	if bufLen <= 0 {
		return 0
	}
	if pos >= bufLen {
		return bufLen - 1
	}
	return pos
}
