package workspace

import (
	"fmt"
	"path/filepath"
)

func (w *Workspace) RunExLine(frameID uint64, line string) (string, error) {
	if msg, handled, err := w.tryExWindowCommand(frameID, line); handled {
		if err != nil {
			return "", err
		}
		w.Editor.SetLastCommand(line)
		if f, ok := w.Frames[frameID]; ok {
			f.Dirty = true
		}
		return msg, nil
	}
	f, win, buf, err := w.FrameContext(frameID)
	if err != nil {
		return "", err
	}

	echo, err := w.Editor.RunExLine(f, win, buf, line)
	if err != nil {
		return "", err
	}
	w.Editor.SetLastCommand(line)
	f.Dirty = true
	return echo, nil
}

func (w *Workspace) ExFindFile(frameID uint64, pattern string) error {
	f, ok := w.Frames[frameID]
	if !ok {
		return errNotFound("frame", frameID)
	}
	matches, _ := filepath.Glob(filepath.Join(f.WorkDir, pattern))
	if len(matches) == 0 {
		return fmt.Errorf("find: no matches")
	}
	_, err := w.OpenFile(frameID, matches[0], false)
	return err
}
