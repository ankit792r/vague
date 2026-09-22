package workspace

import (
	"vague/bonefire/buffer"
	"vague/bonefire/editor"
)

type confirmPrompt struct {
	frameID uint64
	yes     func() error
}

func (w *Workspace) withBufferLeave(frameID uint64, force bool, action func() error) error {
	_, _, buf, err := w.FrameContext(frameID)
	if err != nil {
		return err
	}
	if !buf.Modified() || force {
		return action()
	}
	if w.Editor.FileOptionQuery("hidden") == "set" {
		return action()
	}
	if w.Editor.FileOptionQuery("autowrite") == "set" && buf.Path != "" {
		if err := buf.Save(); err != nil {
			return err
		}
		return action()
	}
	if w.Editor.FileOptionQuery("confirm") == "set" {
		w.pendingConfirm = &confirmPrompt{
			frameID: frameID,
			yes:     action,
		}
		_ = w.SetEcho(frameID, `"`+buf.Name+`" changed. confirm_answer y/n/a/q/l`, editor.EchoInfo)
		return editor.ErrConfirmPending
	}
	return editor.ErrNotSaved
}

func (w *Workspace) AnswerConfirm(frameID uint64, answer string, forceAll bool) error {
	p := w.pendingConfirm
	if p == nil || p.frameID != frameID {
		return nil
	}
	w.pendingConfirm = nil

	_, _, buf, err := w.FrameContext(frameID)
	if err != nil {
		return err
	}

	switch answer {
	case "y", "yes":
		if buf.Path != "" {
			if err := buf.Save(); err != nil {
				return err
			}
		}
		return p.yes()
	case "a", "all":
		if err := w.WriteAll(frameID, false); err != nil {
			return err
		}
		return p.yes()
	case "q", "quit":
		return editor.ErrNotSaved
	case "n", "no", "l", "last":
		return p.yes()
	default:
		return editor.ErrNotSaved
	}
}

func (w *Workspace) maybeAutoread(frameID uint64, buf *buffer.Buffer) {
	if w.Editor.FileOptionQuery("autoread") != "set" {
		return
	}
	if !buf.ChangedOnDisk() {
		return
	}
	if buf.Modified() {
		_ = w.SetEcho(frameID, buf.Name+" changed; use :e! to reload", editor.EchoInfo)
		return
	}
	_ = buf.Reload()
}

func (w *Workspace) CheckTime(frameID uint64) (string, error) {
	var changed int
	for _, buf := range w.Editor.Buffers {
		if buf.Path == "" {
			continue
		}
		if !buf.ChangedOnDisk() {
			continue
		}
		changed++
		if w.Editor.FileOptionQuery("autoread") == "set" && !buf.Modified() {
			_ = buf.Reload()
		}
	}
	if changed == 0 {
		return "checktime: OK", nil
	}
	return "checktime: file(s) changed on disk", nil
}
