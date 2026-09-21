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

func (e *Editor) registerText() ([]byte, yankKind, bool) {
	if len(e.reg.text) == 0 {
		return nil, 0, false
	}
	return e.reg.text, e.reg.kind, true
}
