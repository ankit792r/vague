package workspace

import (
	"strings"

	"vague/bonefire/editor"
)

func (w *Workspace) tryExWindowCommand(frameID uint64, line string) (string, bool, error) {
	kind, args, err := editor.ExCommandKind(line)
	if err != nil {
		return "", false, err
	}
	kind = strings.ToLower(kind)
	switch kind {
	case "split", "sp":
		err := w.SplitWindow(frameID, strings.TrimSpace(args))
		return "", true, err
	case "vsplit", "vs":
		err := w.VSplitWindow(frameID, strings.TrimSpace(args))
		return "", true, err
	case "new":
		err := w.NewWindow(frameID)
		return "", true, err
	case "vnew":
		err := w.VNewWindow(frameID)
		return "", true, err
	case "only":
		err := w.OnlyWindow(frameID)
		return "only", true, err
	case "close", "clo":
		err := w.CloseWindow(frameID)
		return "", true, err
	case "hide":
		err := w.HideWindow(frameID)
		return "", true, err
	case "wincmd", "winc":
		err := w.WinCmd(frameID, strings.TrimSpace(args))
		return "", true, err
	case "tabnew":
		err := w.TabNew(frameID, strings.TrimSpace(args))
		return "", true, err
	case "tabclose", "tabd":
		err := w.TabClose(frameID)
		return "", true, err
	case "tabnext", "tabn":
		err := w.TabNext(frameID)
		return "", true, err
	case "tabprevious", "tabp":
		err := w.TabPrev(frameID)
		return "", true, err
	case "args":
		return w.ArgsList(), true, nil
	case "argadd":
		err := w.ArgAdd(strings.Fields(args)...)
		return "", true, err
	case "argdo":
		msg, err := w.ArgDo(frameID, args)
		return msg, true, err
	case "setlocal":
		msg, err := w.exSetLocal(frameID, args)
		return msg, true, err
	default:
		return "", false, nil
	}
}

func (w *Workspace) exSetLocal(frameID uint64, args string) (string, error) {
	args = strings.TrimSpace(args)
	if args == "" {
		return "", nil
	}
	f, win, _, err := w.FrameContext(frameID)
	if err != nil {
		return "", err
	}
	parts := strings.Fields(args)
	for _, p := range parts {
		if _, err := w.Editor.ApplySetDisplayArgs(f, win, p, true); err == nil {
			continue
		}
		enable := true
		name := p
		if strings.HasPrefix(p, "no") && len(p) > 2 {
			enable = false
			name = p[2:]
		}
		if _, err := w.SetLocalOption(frameID, name, enable); err != nil {
			return "", err
		}
	}
	return args, nil
}
