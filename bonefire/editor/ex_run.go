package editor

import (
	"fmt"
	"strings"

	"vague/bonefire/buffer"
	frame "vague/bonefire/frame"
	"vague/bonefire/window"
)

func (e *Editor) RunExLine(
	frame *frame.Frame,
	win *window.Window,
	buf *buffer.Buffer,
	line string,
) (echo string, err error) {
	cmd, rng, err := parseExLine(line)
	if err != nil {
		return "", err
	}

	switch cmd.kind {
	case "substitute":
		if cmd.sub == nil {
			return "", fmt.Errorf("Invalid substitute")
		}
		if err := e.exSubstitute(frame, win, buf, cmd.sub, rng); err != nil {
			return "", err
		}
		return fmt.Sprintf("%d fewer lines", 0), nil
	case "global", "vglobal":
		return "", e.exGlobal(frame, win, buf, cmd, rng)
	case "normal":
		return "", e.exNormal(frame, win, buf, cmd.args, rng)
	case "read", "r":
		return "", e.exRead(frame, win, buf, cmd.args, rng)
	case "write", "w":
		msg, err := e.exWrite(frame, win, buf, cmd.args, rng, cmd.bang, strings.HasPrefix(strings.TrimSpace(cmd.args), ">>"))
		return msg, err
	case "find", "sf":
		return "", fmt.Errorf("find: handled by workspace")
	case "edit", "e":
		if strings.TrimSpace(cmd.args) == "" {
			return "", fmt.Errorf("use :e! to reload")
		}
		return "", fmt.Errorf("use :edit via execute")
	case "bdelete", "bd", "bwipeout", "bw":
		return e.exBufferDelete(frame, win, buf, cmd.bang)
	case "cd":
		return e.exCd(frame, cmd.args)
	case "pwd":
		return frame.WorkDir, nil
	case "set":
		msg, err := e.exSet(frame, win, cmd.args)
		if err != nil {
			return "", err
		}
		switch msg {
		case "wrap":
			return "", e.setWrap(frame, win, true)
		case "nowrap":
			return "", e.setWrap(frame, win, false)
		case "number", "nu":
			return "", e.setNumber(frame, win, true)
		case "nonumber", "nonu":
			return "", e.setNumber(frame, win, false)
		default:
			if strings.HasSuffix(msg, "?") {
				return msg, nil
			}
			frame.Dirty = true
			return msg, nil
		}
	case "map", "nmap", "imap", "vmap":
		return e.exMap(cmd.kind, cmd.args)
	case "command":
		return e.exUserCommand(cmd.args)
	case "source":
		return "", e.exSource(frame, win, buf, cmd.args)
	case "reg", "registers":
		return e.exReg(), nil
	case "marks":
		return e.exMarks(), nil
	case "delm":
		return e.exDelMark(cmd.args), nil
	case "digraphs", "dig":
		return e.exDigraphs(cmd.args), nil
	case "jumps":
		return e.exJumps(), nil
	case "clearjumps":
		e.clearJumps()
		return "jump list cleared", nil
	case "undo", "u":
		return "", e.undoTo(frame, win, buf, false)
	case "redo", "red":
		if err := e.undoTo(frame, win, buf, true); err != nil {
			return "", err
		}
		return "redo", nil
	case "undolist":
		return e.exUndolist(buf), nil
	case "changes":
		return e.exChanges(buf), nil
	case "later":
		return "", e.exLater(frame, win, buf, cmd.args)
	case "earlier":
		return "", e.exEarlier(frame, win, buf, cmd.args)
	case "checktime":
		return "", fmt.Errorf("checktime: use workspace CheckTime")
	case "diffsplit", "diffoff", "diffget", "diffput":
		return "", fmt.Errorf("%s not implemented", cmd.kind)
	case "terminal", "term":
		return "", fmt.Errorf("terminal not implemented")
	case "nohlsearch", "noh":
		e.exNohlSearch(frame)
		return "", nil
	case "vimgrep", "vim", "lvimgrep", "lvim":
		return e.exVimgrep(cmd.args)
	case "ta", "tag", "tags":
		return e.exTag(cmd.args)
	case "sort":
		return "", e.exSort(frame, win, buf, cmd.args, rng)
	case "uniq":
		return "", e.exUniq(frame, win, buf, rng)
	default:
		return "", fmt.Errorf("Unknown Ex command: %s", cmd.kind)
	}
}
