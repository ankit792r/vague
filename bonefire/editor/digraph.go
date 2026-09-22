package editor

import (
	"fmt"
	"sort"
	"strings"
	"unicode/utf8"
)

var digraphTable = map[string]rune{
	"Co":  '\u00a9',
	"e'":  '\u00e9',
	"e`":  '\u00e8',
	"u'":  '\u00fa',
	"a'":  '\u00e1',
	"o'":  '\u00f3',
	"i'":  '\u00ed',
	"c,":  '\u00e7',
	"u:":  '\u00fc',
	"o:":  '\u00f6',
	"a:":  '\u00e4',
	"~N":  '\u00f1',
	"~n":  '\u00f1',
	"->":  '\u2192',
	"<-":  '\u2190',
	"!=":  '\u2260',
	"=<":  '\u2264',
	"=>":  '\u2265',
	"Inf": '\u221e',
	"+-":  '\u00b1',
	"1/2": '\u00bd',
	"1/4": '\u00bc',
	"3/4": '\u00be',
}

func lookupDigraph(a, b rune) (rune, bool) {
	key := string([]rune{a, b})
	if r, ok := digraphTable[key]; ok {
		return r, true
	}
	return 0, false
}

func (e *Editor) consumeDigraphInsert(keys string) (string, bool) {
	if !e.pendingDigraph {
		return "", false
	}
	if len(keys) != 1 || keys[0] == '<' {
		e.pendingDigraph = false
		e.digraphFirst = ""
		return "", false
	}
	first := e.digraphFirst
	e.pendingDigraph = false
	e.digraphFirst = ""
	if first == "" {
		return "", false
	}
	r, ok := lookupDigraph(rune(first[0]), rune(keys[0]))
	if !ok {
		return "", false
	}
	return string(r), true
}

func (e *Editor) beginDigraph() {
	e.pendingDigraph = true
	e.digraphFirst = ""
}

func (e *Editor) setDigraphFirst(keys string) bool {
	if len(keys) != 1 || keys[0] == '<' {
		e.pendingDigraph = false
		return false
	}
	e.digraphFirst = keys
	return true
}

func (e *Editor) exDigraphs(args string) string {
	args = strings.TrimSpace(args)
	if args != "" {
		parts := strings.Fields(args)
		if len(parts) >= 3 {
			a, b := parts[0], parts[1]
			if len(a) == 1 && len(b) == 1 {
				ch := parts[2]
				r, _ := utf8.DecodeRuneInString(ch)
				if r != utf8.RuneError {
					digraphTable[a+b] = r
					return "digraph added"
				}
			}
		}
		return "usage: :digraphs xy char"
	}
	type row struct {
		key string
		val rune
	}
	var rows []row
	for k, v := range digraphTable {
		rows = append(rows, row{k, v})
	}
	sort.Slice(rows, func(i, j int) bool { return rows[i].key < rows[j].key })
	var b strings.Builder
	for i, r := range rows {
		if i > 0 {
			b.WriteByte('\n')
		}
		fmt.Fprintf(&b, "%s  %U  %c", r.key, r.val, r.val)
	}
	return b.String()
}
