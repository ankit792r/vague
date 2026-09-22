package workspace

import (
	"fmt"
	"strings"
)

func (w *Workspace) ArgsList() string {
	if len(w.ArgList) == 0 {
		return "argument list empty"
	}
	var b strings.Builder
	for i, p := range w.ArgList {
		if i > 0 {
			b.WriteByte(' ')
		}
		if i == w.ArgIndex {
			b.WriteByte('[')
			b.WriteString(p)
			b.WriteByte(']')
		} else {
			b.WriteString(p)
		}
	}
	return b.String()
}

func (w *Workspace) ArgAdd(paths ...string) error {
	for _, p := range paths {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		w.ArgList = append(w.ArgList, p)
	}
	if w.ArgIndex < 0 && len(w.ArgList) > 0 {
		w.ArgIndex = 0
	}
	return nil
}

func (w *Workspace) ArgDo(frameID uint64, cmdLine string) (string, error) {
	cmdLine = strings.TrimSpace(cmdLine)
	if cmdLine == "" {
		return "", fmt.Errorf("argdo: command required")
	}
	if len(w.ArgList) == 0 {
		return "", fmt.Errorf("argument list empty")
	}
	var ran int
	for i, path := range w.ArgList {
		w.ArgIndex = i
		if _, err := w.OpenFile(frameID, path, false); err != nil {
			continue
		}
		if _, err := w.RunExLine(frameID, cmdLine); err != nil {
			return "", err
		}
		ran++
	}
	return fmt.Sprintf("argdo: %d files", ran), nil
}
