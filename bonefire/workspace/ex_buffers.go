package workspace

import (
	"fmt"
	"strings"

	"vague/bonefire/editor"
)

func (w *Workspace) tryExBufferCommand(frameID uint64, line string) (string, bool, error) {
	kind, args, err := editor.ExCommandKind(line)
	if err != nil {
		return "", false, err
	}
	kind = strings.ToLower(kind)
	switch kind {
	case "confirm":
		msg, err := w.RunConfirmEx(frameID, args)
		return msg, true, err
	case "wa", "wall":
		err := w.WriteAll(frameID, strings.Contains(line, "!"))
		return "", true, err
	case "qa", "qall":
		err := w.QuitAll(strings.HasSuffix(strings.TrimSpace(line), "!"))
		return "", true, err
	case "wqa", "wqall":
		force := strings.HasSuffix(strings.TrimSpace(line), "!")
		err := w.WQAll(frameID, force)
		return "", true, err
	case "xa", "xall":
		force := strings.HasSuffix(strings.TrimSpace(line), "!")
		err := w.XAll(frameID, force)
		return "", true, err
	case "lcd":
		return w.exLcd(frameID, args)
	case "xxd":
		return "xxd not implemented (use external tools)", true, nil
	case "checktime":
		msg, err := w.CheckTime(frameID)
		return msg, true, err
	case "telescope", "pick":
		return "", true, fmt.Errorf("fuzzy finder not implemented")
	default:
		return "", false, nil
	}
}

func (w *Workspace) exLcd(frameID uint64, path string) (string, bool, error) {
	frame, win, _, err := w.FrameContext(frameID)
	if err != nil {
		return "", true, err
	}
	path = strings.TrimSpace(path)
	if path == "" {
		if win.LocalWorkDir != "" {
			return win.LocalWorkDir, true, nil
		}
		return frame.WorkDir, true, nil
	}
	abs, err := w.resolvePath(frameID, path)
	if err != nil {
		return "", true, err
	}
	win.LocalWorkDir = abs
	frame.Dirty = true
	return abs, true, nil
}

func (w *Workspace) RunConfirmEx(frameID uint64, inner string) (string, error) {
	inner = strings.TrimSpace(inner)
	if inner == "" {
		return "", fmt.Errorf("confirm: command required")
	}
	var msg string
	err := w.withBufferLeave(frameID, false, func() error {
		var err error
		msg, err = w.RunExLine(frameID, inner)
		return err
	})
	if err != nil {
		return msg, err
	}
	return msg, nil
}
