package editor

import (
	"fmt"
	"strconv"
	"strings"
)

type exRange struct {
	startLine int
	endLine   int
	wholeBuf  bool
}

type exSubstitute struct {
	pattern     string
	replacement string
	flags       subFlags
}

type subFlags struct {
	global      bool // g
	ignoreCase  bool // i
	noIgnoreCase bool // I
}

type exCommand struct {
	kind string
	sub  *exSubstitute
	raw  string
	args string
	bang bool
}

func parseExLine(line string) (exCommand, exRange, error) {
	line = strings.TrimSpace(line)
	if line == "" {
		return exCommand{}, exRange{}, fmt.Errorf("No command")
	}

	bang := false
	if strings.HasSuffix(line, "!") {
		bang = true
		line = strings.TrimSpace(line[:len(line)-1])
	}

	rng, rest, err := parseExRangePrefix(line)
	if err != nil {
		return exCommand{}, exRange{}, err
	}
	if rest == "" {
		return exCommand{}, rng, fmt.Errorf("No command")
	}

	if isExSubstitute(rest) {
		sub, err := parseSubstitute(rest)
		if err != nil {
			return exCommand{}, exRange{}, err
		}
		return exCommand{kind: "substitute", sub: sub, bang: bang}, rng, nil
	}

	if strings.HasPrefix(rest, "g") || strings.HasPrefix(rest, "v") {
		kind := "global"
		if rest[0] == 'v' {
			kind = "vglobal"
		}
		return exCommand{kind: kind, raw: rest, bang: bang}, rng, nil
	}

	if strings.HasPrefix(rest, "norm") {
		args := strings.TrimSpace(rest[len("norm"):])
		if strings.HasPrefix(args, "al") {
			args = strings.TrimSpace(args[2:])
		}
		return exCommand{kind: "normal", args: args, bang: bang}, rng, nil
	}

	cmdName, args := splitExNameArgs(rest)
	cmdName = strings.ToLower(cmdName)

	return exCommand{kind: cmdName, args: args, bang: bang}, rng, nil
}

// isExSubstitute reports whether rest is :substitute (s/pat/repl/), not :set or :sort.
func isExSubstitute(rest string) bool {
	if len(rest) < 2 || rest[0] != 's' {
		return false
	}
	lower := strings.ToLower(rest)
	if strings.HasPrefix(lower, "set") {
		if len(lower) == 3 {
			return false
		}
		if lower[3] == ' ' || lower[3] == '\t' {
			return false
		}
	}
	if strings.HasPrefix(lower, "sort") || strings.HasPrefix(lower, "source") {
		return false
	}
	delim := rest[1]
	return delim != ' ' && delim != '\t'
}

func splitExNameArgs(rest string) (string, string) {
	i := 0
	for i < len(rest) && rest[i] != ' ' && rest[i] != '\t' {
		i++
	}
	if i >= len(rest) {
		return rest, ""
	}
	return rest[:i], strings.TrimSpace(rest[i:])
}

func parseExRangePrefix(line string) (exRange, string, error) {
	rng := exRange{startLine: -1, endLine: -1}

	if strings.HasPrefix(line, "%") {
		rng.wholeBuf = true
		return rng, strings.TrimSpace(line[1:]), nil
	}

	if strings.HasPrefix(line, "'<,'>") {
		rng.startLine = -2
		rng.endLine = -2
		return rng, strings.TrimSpace(line[len("'<,'>"):]), nil
	}

	// line number range: 5,10s or 5s or .,.$ etc — find where command starts (letter after range)
	cmdStart := findExCommandStart(line)
	if cmdStart <= 0 {
		return rng, line, nil
	}

	prefix := strings.TrimSpace(line[:cmdStart])
	rest := strings.TrimSpace(line[cmdStart:])
	if prefix == "" {
		return rng, rest, nil
	}

	if prefix == "." {
		rng.startLine = -3
		rng.endLine = -3
		return rng, rest, nil
	}
	if prefix == "$" {
		rng.startLine = -4
		rng.endLine = -4
		return rng, rest, nil
	}

	if strings.Contains(prefix, ",") {
		parts := strings.SplitN(prefix, ",", 2)
		a, err := parseLineSpec(parts[0])
		if err != nil {
			return exRange{}, "", err
		}
		b, err := parseLineSpec(parts[1])
		if err != nil {
			return exRange{}, "", err
		}
		rng.startLine = a
		rng.endLine = b
		return rng, rest, nil
	}

	n, err := parseLineSpec(prefix)
	if err != nil {
		return exRange{}, "", err
	}
	rng.startLine = n
	rng.endLine = n
	return rng, rest, nil
}

func findExCommandStart(line string) int {
	for i := 0; i < len(line); i++ {
		c := line[i]
		if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') {
			if i > 0 {
				prev := line[i-1]
				if prev == ',' || prev == ';' || prev == '%' || prev == '>' || prev == '\'' {
					continue
				}
			}
			return i
		}
	}
	return -1
}

func parseLineSpec(s string) (int, error) {
	s = strings.TrimSpace(s)
	if s == "." {
		return -3, nil
	}
	if s == "$" {
		return -4, nil
	}
	n, err := strconv.Atoi(s)
	if err != nil || n < 1 {
		return 0, fmt.Errorf("Invalid line number")
	}
	return n - 1, nil
}

func parseSubstitute(rest string) (*exSubstitute, error) {
	if len(rest) < 2 {
		return nil, fmt.Errorf("Invalid substitute command")
	}
	delim := rest[1]
	if delim == ' ' || delim == '\t' {
		return nil, fmt.Errorf("Invalid substitute delimiter")
	}

	body := rest[2:]
	pat, tail, ok := splitDelimField(body, delim)
	if !ok || pat == "" {
		return nil, fmt.Errorf("Invalid pattern")
	}
	repl, tail, ok := splitDelimField(tail, delim)
	if !ok {
		return nil, fmt.Errorf("Invalid replacement")
	}

	flags := subFlags{}
	for _, f := range tail {
		switch f {
		case 'g':
			flags.global = true
		case 'i':
			flags.ignoreCase = true
		case 'I':
			flags.noIgnoreCase = true
		case 'c':
			// confirm — no interactive UI yet
		default:
			return nil, fmt.Errorf("Invalid substitute flag: %c", f)
		}
	}

	return &exSubstitute{pattern: pat, replacement: repl, flags: flags}, nil
}

func splitDelimField(s string, delim byte) (string, string, bool) {
	for i := 0; i < len(s); i++ {
		if s[i] == delim && (i == 0 || s[i-1] != '\\') {
			return s[:i], s[i+1:], true
		}
	}
	return s, "", false
}
