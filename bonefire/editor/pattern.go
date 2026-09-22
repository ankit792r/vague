package editor

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
	"unicode/utf8"
)

// MagicMode selects which characters are special in a search pattern.
// See PATTERN.md for the supported subset.
type MagicMode int

const (
	MagicMagic MagicMode = iota // \m — default Vim magic
	MagicNoMagic                // \M
	MagicVeryMagic              // \v
	MagicVeryNoMagic            // \V
)

type parsedPattern struct {
	re          *regexp.Regexp
	literal     []byte
	useLiteral  bool
	ignoreCase  bool
	displayPat  string
}

func (e *Editor) parseSearchPattern(raw string) (parsedPattern, error) {
	pat := raw
	mode := e.searchOpts.Magic
	forceIgnore := false
	forceNoIgnore := false

	for len(pat) >= 2 && pat[0] == '\\' {
		switch pat[1] {
		case 'm':
			mode = MagicMagic
			pat = pat[2:]
		case 'M':
			mode = MagicNoMagic
			pat = pat[2:]
		case 'v':
			mode = MagicVeryMagic
			pat = pat[2:]
		case 'V':
			mode = MagicVeryNoMagic
			pat = pat[2:]
		case 'c':
			forceIgnore = true
			pat = pat[2:]
		case 'C':
			forceNoIgnore = true
			pat = pat[2:]
		default:
			goto donePrefix
		}
	}
donePrefix:

	if pat == "" {
		return parsedPattern{}, fmt.Errorf("empty pattern")
	}

	ignore := e.effectiveIgnoreCase(pat, forceIgnore, forceNoIgnore)
	if isPlainLiteral(pat, mode) {
		lit := []byte(pat)
		if ignore {
			lit = bytesToLower(bytes.Clone(lit))
		}
		return parsedPattern{
			literal:    lit,
			useLiteral: true,
			ignoreCase: ignore,
			displayPat: raw,
		}, nil
	}

	reStr, err := vimPatternToRegexp(pat, mode)
	if err != nil {
		return parsedPattern{}, err
	}
	if ignore {
		reStr = "(?i)" + reStr
	}
	re, err := regexp.Compile(reStr)
	if err != nil {
		return parsedPattern{}, fmt.Errorf("invalid pattern: %w", err)
	}
	return parsedPattern{
		re:         re,
		useLiteral: false,
		ignoreCase: ignore,
		displayPat: raw,
	}, nil
}

func (e *Editor) effectiveIgnoreCase(pat string, forceIgnore, forceNoIgnore bool) bool {
	if forceNoIgnore {
		return false
	}
	if forceIgnore {
		return true
	}
	if e.searchOpts.IgnoreCase {
		if e.searchOpts.SmartCase && patternHasUpperASCII(pat) {
			return false
		}
		return true
	}
	return false
}

func patternHasUpperASCII(s string) bool {
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			return true
		}
	}
	return false
}

func isPlainLiteral(pat string, mode MagicMode) bool {
	if mode == MagicVeryNoMagic {
		return true
	}
	for i := 0; i < len(pat); i++ {
		if pat[i] != '\\' && isSpecialInMode(pat[i], mode) {
			return false
		}
	}
	return true
}

func isSpecialInMode(b byte, mode MagicMode) bool {
	switch mode {
	case MagicVeryNoMagic:
		return false
	case MagicNoMagic:
		return b == '^' || b == '$'
	case MagicVeryMagic:
		switch b {
		case '.', '*', '+', '?', '^', '$', '[', ']', '(', ')', '{', '}', '|', '\\':
			return true
		default:
			return false
		}
	default: // MagicMagic
		switch b {
		case '.', '^', '$', '[', ']', '\\', '~':
			return true
		default:
			return false
		}
	}
}

func vimPatternToRegexp(pat string, mode MagicMode) (string, error) {
	var b strings.Builder
	i := 0
	for i < len(pat) {
		c := pat[i]
		if c == '\\' && i+1 < len(pat) {
			next := pat[i+1]
			switch next {
			case 'n':
				b.WriteString(`\n`)
				i += 2
				continue
			case 't':
				b.WriteString(`\t`)
				i += 2
				continue
			case '.', '*', '+', '?', '^', '$', '[', ']', '(', ')', '{', '}', '|', '\\':
				b.WriteByte('\\')
				b.WriteByte(next)
				i += 2
				continue
			default:
				r, size := utf8.DecodeRuneInString(pat[i+1:])
				if size == 0 {
					return "", fmt.Errorf("bad escape")
				}
				b.WriteString(regexp.QuoteMeta(string(r)))
				i += 1 + size
				continue
			}
		}

		if isSpecialInMode(c, mode) {
			switch c {
			case '.':
				b.WriteString(`[^\n]`)
			case '^':
				if i == 0 {
					b.WriteString(`(?m)^`)
				} else {
					b.WriteString(`\^`)
				}
			case '$':
				if i == len(pat)-1 {
					b.WriteString(`(?m)$`)
				} else {
					b.WriteString(`\$`)
				}
			case '[':
				end := strings.IndexByte(pat[i:], ']')
				if end < 0 {
					return "", fmt.Errorf("unclosed [")
				}
				b.WriteByte('[')
				b.WriteString(pat[i+1 : i+end])
				b.WriteByte(']')
				i += end + 1
				continue
			case '~':
				b.WriteString(`~`)
			default:
				b.WriteByte('\\')
				b.WriteByte(c)
			}
			i++
			continue
		}

		b.WriteString(regexp.QuoteMeta(string(c)))
		i++
	}
	return b.String(), nil
}

func findPatternForward(data []byte, pp parsedPattern, after int, wrap bool) (start, end int, ok bool) {
	if pp.useLiteral {
		return findLiteralForward(data, pp.literal, after, wrap, pp.ignoreCase)
	}
	re := pp.re
	if re == nil {
		return 0, 0, false
	}
	search := data
	offset := 0
	if after >= 0 && after < len(data) {
		search = data[after+1:]
		offset = after + 1
	}
	loc := re.FindIndex(search)
	if loc != nil {
		return offset + loc[0], offset + loc[1], true
	}
	if !wrap || after < 0 {
		return 0, 0, false
	}
	loc = re.FindIndex(data)
	if loc != nil {
		return loc[0], loc[1], true
	}
	return 0, 0, false
}

func findPatternBackward(data []byte, pp parsedPattern, before int, wrap bool) (start, end int, ok bool) {
	if pp.useLiteral {
		return findLiteralBackward(data, pp.literal, before, wrap, pp.ignoreCase)
	}
	re := pp.re
	if re == nil {
		return 0, 0, false
	}
	bestStart := -1
	bestEnd := -1
	all := re.FindAllIndex(data, -1)
	for _, loc := range all {
		if loc[0] < before && loc[0] > bestStart {
			bestStart = loc[0]
			bestEnd = loc[1]
		}
	}
	if bestStart >= 0 {
		return bestStart, bestEnd, true
	}
	if !wrap {
		return 0, 0, false
	}
	if len(all) == 0 {
		return 0, 0, false
	}
	loc := all[len(all)-1]
	return loc[0], loc[1], true
}

func findLiteralForward(data, pat []byte, after int, wrap, ignoreCase bool) (int, int, bool) {
	if len(pat) == 0 {
		return 0, 0, false
	}
	hay := data
	if ignoreCase {
		hay = bytesToLower(bytes.Clone(data))
	}
	at := after
	if at < 0 {
		at = -1
	}
	for i := 0; i+len(pat) <= len(hay); i++ {
		if bytesEqualFoldSlice(hay[i:i+len(pat)], pat) && i > at {
			return i, i + len(pat), true
		}
	}
	if !wrap {
		return 0, 0, false
	}
	for i := 0; i+len(pat) <= len(hay); i++ {
		if bytesEqualFoldSlice(hay[i:i+len(pat)], pat) {
			return i, i + len(pat), true
		}
	}
	return 0, 0, false
}

func findLiteralBackward(data, pat []byte, before int, wrap, ignoreCase bool) (int, int, bool) {
	if len(pat) == 0 {
		return 0, 0, false
	}
	hay := data
	if ignoreCase {
		hay = bytesToLower(bytes.Clone(data))
	}
	best := -1
	for i := 0; i+len(pat) <= len(hay); i++ {
		if bytesEqualFoldSlice(hay[i:i+len(pat)], pat) && i < before {
			best = i
		}
	}
	if best >= 0 {
		return best, best + len(pat), true
	}
	if !wrap {
		return 0, 0, false
	}
	for i := len(hay) - len(pat); i >= 0; i-- {
		if bytesEqualFoldSlice(hay[i:i+len(pat)], pat) {
			return i, i + len(pat), true
		}
	}
	return 0, 0, false
}

func bytesEqualFoldSlice(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func findAllPattern(data []byte, pp parsedPattern) [][2]int {
	if pp.useLiteral {
		return findAllLiteral(data, pp.literal, pp.ignoreCase)
	}
	if pp.re == nil {
		return nil
	}
	locs := pp.re.FindAllIndex(data, -1)
	out := make([][2]int, 0, len(locs))
	for _, loc := range locs {
		out = append(out, [2]int{loc[0], loc[1]})
	}
	return out
}

func findAllLiteral(data, pat []byte, ignoreCase bool) [][2]int {
	if len(pat) == 0 {
		return nil
	}
	hay := data
	if ignoreCase {
		hay = bytesToLower(bytes.Clone(data))
	}
	var out [][2]int
	for i := 0; i+len(pat) <= len(hay); i++ {
		if bytesEqualFoldSlice(hay[i:i+len(pat)], pat) {
			out = append(out, [2]int{i, i + len(pat)})
		}
	}
	return out
}

// findPatternInSlice returns the first match of pp in hay.
func findPatternInSlice(hay []byte, pp parsedPattern) (start, end int, ok bool) {
	if len(hay) == 0 {
		return 0, 0, false
	}
	all := findAllPattern(hay, pp)
	if len(all) == 0 {
		return 0, 0, false
	}
	m := all[0]
	return m[0], m[1], true
}

func (e *Editor) parseExPattern(raw string, flags subFlags) (parsedPattern, error) {
	prev := e.searchOpts.IgnoreCase
	if flags.ignoreCase && !flags.noIgnoreCase {
		e.searchOpts.IgnoreCase = true
	} else if flags.noIgnoreCase {
		e.searchOpts.IgnoreCase = false
	}
	pp, err := e.parseSearchPattern(raw)
	e.searchOpts.IgnoreCase = prev
	return pp, err
}
