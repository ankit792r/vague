package editor

import (
	"bytes"
	"fmt"
	"strings"
)

type yankKind int

const (
	yankCharwise yankKind = iota
	yankLinewise
)

type regSlot struct {
	text []byte
	kind yankKind
}

func (s regSlot) empty() bool {
	return len(s.text) == 0
}

type registerID struct {
	name   string
	append bool
}

type registers struct {
	unnamed regSlot
	named   map[string]regSlot
	yank0   regSlot
	ring    [9]regSlot
	small   regSlot
	colon   string
	slash   string
	star    regSlot
	plus    regSlot
}

func newRegisters() registers {
	return registers{named: make(map[string]regSlot)}
}

const smallDeleteMax = 4096

func parseRegisterKeyASCII(keys string) (registerID, bool) {
	if len(keys) != 1 {
		return registerID{}, false
	}
	r := rune(keys[0])
	switch {
	case r >= 'a' && r <= 'z':
		return registerID{name: string(r)}, true
	case r >= 'A' && r <= 'Z':
		return registerID{name: strings.ToLower(string(r)), append: true}, true
	case r >= '0' && r <= '9':
		return registerID{name: string(r)}, true
	case r == '-', r == '_', r == '=', r == ':', r == '/', r == '*', r == '+', r == '%', r == '#':
		return registerID{name: string(r)}, true
	default:
		return registerID{}, false
	}
}

func parseRegisterName(name string) (registerID, bool) {
	if name == "" || name == `"` {
		return registerID{}, true
	}
	if len(name) == 2 && name[0] >= 'A' && name[0] <= 'Z' {
		return registerID{name: strings.ToLower(string(name[0])), append: true}, true
	}
	return parseRegisterKeyASCII(name)
}

func (e *Editor) beginRegisterSelect() {
	e.pendingRegQuote = true
}

func (e *Editor) consumeRegisterSelect(keys string) bool {
	if !e.pendingRegQuote {
		return false
	}
	e.pendingRegQuote = false
	id, ok := parseRegisterKeyASCII(keys)
	if ok {
		e.activeReg = id
	}
	return true
}

func (e *Editor) slotFromCell(data []byte, linewise bool) regSlot {
	kind := yankCharwise
	if linewise {
		kind = yankLinewise
	}
	return regSlot{text: bytes.Clone(data), kind: kind}
}

func (e *Editor) recordYank(data []byte, linewise bool) {
	if len(data) == 0 || e.activeReg.name == "_" {
		return
	}
	cell := e.slotFromCell(data, linewise)
	e.writeRegister(e.activeReg, cell)
	e.regs.yank0 = cell
}

func (e *Editor) recordDelete(data []byte, linewise bool) {
	if len(data) == 0 {
		return
	}
	cell := e.slotFromCell(data, linewise)
	if e.activeReg.name != "_" {
		e.pushDeleteRing(cell)
		if len(data) <= smallDeleteMax {
			e.regs.small = cell
		}
		e.writeRegister(e.activeReg, cell)
	}
}

func (e *Editor) pushDeleteRing(cell regSlot) {
	for i := len(e.regs.ring) - 1; i > 0; i-- {
		e.regs.ring[i] = e.regs.ring[i-1]
	}
	e.regs.ring[0] = cell
}

func numberedRegister(name string) (int, bool) {
	if len(name) != 1 || name[0] < '1' || name[0] > '9' {
		return 0, false
	}
	return int(name[0] - '0'), true
}

func (e *Editor) writeRegister(id registerID, cell regSlot) {
	if id.name == "_" || id.name == "=" {
		return
	}

	switch id.name {
	case "":
		e.regs.unnamed = cell
	case "0":
		e.regs.yank0 = cell
		e.regs.unnamed = cell
	case "-":
		e.regs.small = cell
		e.regs.unnamed = cell
	case "*":
		e.regs.star = cell
		e.regs.unnamed = cell
		clipboardWritePrimary(string(cell.text))
		clipboardOSC52(cell.text)
	case "+":
		e.regs.plus = cell
		e.regs.unnamed = cell
		clipboardWriteClipboard(string(cell.text))
		clipboardOSC52(cell.text)
	default:
		if n, ok := numberedRegister(id.name); ok {
			e.regs.ring[n-1] = cell
			e.regs.unnamed = cell
			return
		}
		if id.name >= "a" && id.name <= "z" {
			if id.append {
				prev := e.regs.named[id.name]
				combined := append(bytes.Clone(prev.text), cell.text...)
				kind := prev.kind
				if cell.kind == yankLinewise {
					kind = yankLinewise
				}
				cell = regSlot{text: combined, kind: kind}
			}
			e.regs.named[id.name] = cell
			e.regs.unnamed = cell
		}
	}
}

func (e *Editor) readRegister(id registerID) ([]byte, yankKind, bool) {
	name := id.name
	switch name {
	case "":
		return e.slotBytes(e.regs.unnamed)
	case "0":
		return e.slotBytes(e.regs.yank0)
	case "-":
		return e.slotBytes(e.regs.small)
	case ":":
		if e.regs.colon == "" {
			return nil, 0, false
		}
		return []byte(e.regs.colon), yankCharwise, true
	case "/":
		if e.regs.slash == "" {
			return nil, 0, false
		}
		return []byte(e.regs.slash), yankCharwise, true
	case "*":
		if text, ok := clipboardReadPrimary(); ok {
			return []byte(text), yankCharwise, true
		}
		return e.slotBytes(e.regs.star)
	case "+":
		if text, ok := clipboardReadClipboard(); ok {
			return []byte(text), yankCharwise, true
		}
		return e.slotBytes(e.regs.plus)
	case "_", "=":
		return nil, 0, false
	default:
		if n, ok := numberedRegister(name); ok {
			return e.slotBytes(e.regs.ring[n-1])
		}
		if cell, ok := e.regs.named[name]; ok {
			return e.slotBytes(cell)
		}
	}
	return nil, 0, false
}

func (e *Editor) slotBytes(s regSlot) ([]byte, yankKind, bool) {
	if s.empty() {
		return nil, 0, false
	}
	return bytes.Clone(s.text), s.kind, true
}

func (e *Editor) registerText() ([]byte, yankKind, bool) {
	return e.readRegister(e.activeReg)
}

func (e *Editor) SetRegister(name string, text string, linewise bool) error {
	id, ok := parseRegisterName(name)
	if !ok {
		return fmt.Errorf("invalid register: %q", name)
	}
	cell := e.slotFromCell([]byte(text), linewise)
	e.writeRegister(id, cell)
	return nil
}

func (e *Editor) GetRegister(name string) (string, bool) {
	id, ok := parseRegisterName(name)
	if !ok {
		return "", false
	}
	data, _, ok := e.readRegister(id)
	if !ok {
		return "", false
	}
	return string(data), true
}

func (e *Editor) SetLastSearch(pattern string) {
	e.regs.slash = pattern
}

func (e *Editor) SetLastCommand(line string) {
	e.regs.colon = line
}

func (e *Editor) UnnamedRegisterString() string {
	data, _, ok := e.slotBytes(e.regs.unnamed)
	if !ok {
		return ""
	}
	return string(data)
}

func (e *Editor) FormatRegisters() string {
	var b strings.Builder
	b.WriteString("Type Name Content\n")
	appendLine := func(typ, name, content string) {
		if content == "" {
			return
		}
		if len(content) > 60 {
			content = content[:57] + "..."
		}
		fmt.Fprintf(&b, "%s %q %s\n", typ, name, content)
	}
	if s, _, ok := e.slotBytes(e.regs.unnamed); ok {
		appendLine("c", `"`, string(s))
	}
	for i := 9; i >= 1; i-- {
		if s, _, ok := e.slotBytes(e.regs.ring[i-1]); ok {
			appendLine("c", fmt.Sprintf("%d", i), string(s))
		}
	}
	if s, _, ok := e.slotBytes(e.regs.yank0); ok {
		appendLine("c", "0", string(s))
	}
	if s, _, ok := e.slotBytes(e.regs.small); ok {
		appendLine("c", "-", string(s))
	}
	appendLine("c", ":", e.regs.colon)
	appendLine("c", "/", e.regs.slash)
	for c := 'a'; c <= 'z'; c++ {
		name := string(c)
		if cell, ok := e.regs.named[name]; ok {
			if s, _, ok := e.slotBytes(cell); ok {
				appendLine("c", name, string(s))
			}
		}
	}
	if s, _, ok := e.slotBytes(e.regs.star); ok {
		appendLine("c", "*", string(s))
	}
	if s, _, ok := e.slotBytes(e.regs.plus); ok {
		appendLine("c", "+", string(s))
	}
	return strings.TrimSpace(b.String())
}

func (e *Editor) setRegisterLinewise(text []byte) {
	e.recordYank(text, true)
}

func (e *Editor) setRegisterCharwise(text []byte) {
	e.recordYank(text, false)
}

func (e *Editor) stashRegister(data []byte, linewise bool, isDelete bool) {
	if isDelete {
		e.recordDelete(data, linewise)
	} else {
		e.recordYank(data, linewise)
	}
}
