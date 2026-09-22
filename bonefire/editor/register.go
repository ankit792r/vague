package editor

import "bytes"

type yankKind int

const (
	yankCharwise yankKind = iota
	yankLinewise
)

type register struct {
	text []byte
	kind yankKind
}

func (e *Editor) setRegisterLinewise(text []byte) {
	e.reg.text = bytes.Clone(text)
	e.reg.kind = yankLinewise
}

func (e *Editor) setRegisterCharwise(text []byte) {
	e.reg.text = bytes.Clone(text)
	e.reg.kind = yankCharwise
}

func (e *Editor) stashRegister(data []byte, linewise bool) {
	if linewise {
		e.setRegisterLinewise(data)
	} else {
		e.setRegisterCharwise(data)
	}
}

func (e *Editor) registerText() ([]byte, yankKind, bool) {
	if len(e.reg.text) == 0 {
		return nil, 0, false
	}
	return e.reg.text, e.reg.kind, true
}

func (e *Editor) UnnamedRegisterString() string {
	data, _, ok := e.registerText()
	if !ok {
		return ""
	}
	return string(data)
}
