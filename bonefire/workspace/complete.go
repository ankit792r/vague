package workspace

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
)

var exCommandNames = []string{
	"edit", "e", "write", "w", "quit", "q", "wq", "x",
	"buffer", "b", "buffers", "ls", "bnext", "bn", "bprev", "bp",
	"goto", "go", "search", "set", "wrap", "nowrap", "number", "nonumber",
	"only", "help",
}

var setOptionNames = []string{
	"number", "nu", "nonumber", "nonu", "wrap", "nowrap",
	"ignorecase", "ic", "noignorecase", "noic",
	"smartcase", "scs", "nosmartcase", "noscs",
	"hlsearch", "hls", "nohlsearch", "nohls",
	"incsearch", "noincsearch",
	"wrapscan", "ws", "nowrapscan", "nows",
	"swapfile", "swf",
	"paste", "nopaste",
}

func (w *Workspace) Complete(frameID uint64, kind, prefix string) ([]string, error) {
	switch kind {
	case "command":
		return filterPrefix(exCommandNames, strings.ToLower(prefix)), nil
	case "set":
		return filterPrefix(setOptionNames, strings.ToLower(prefix)), nil
	case "buffer":
		_, _, buf, err := w.FrameContext(frameID)
		if err != nil {
			return nil, err
		}
		return w.completeBufferNames(buf.ID, strings.ToLower(prefix)), nil
	case "path":
		frame, ok := w.Frames[frameID]
		if !ok {
			return nil, errNotFound("frame", frameID)
		}
		return completePaths(frame.WorkDir, prefix), nil
	default:
		return nil, nil
	}
}

func (w *Workspace) completeBufferNames(currentID uint64, prefix string) []string {
	names := make([]string, 0, len(w.Editor.Buffers))
	for _, buf := range w.Editor.Buffers {
		name := buf.Name
		if buf.Path != "" {
			name = buf.Path
		}
		if prefix == "" || strings.HasPrefix(strings.ToLower(name), prefix) {
			names = append(names, name)
		}
	}
	sort.Strings(names)
	return names
}

func completePaths(workDir, prefix string) []string {
	dir := workDir
	base := prefix
	if strings.Contains(prefix, "/") {
		lastSlash := strings.LastIndex(prefix, "/")
		dir = filepath.Join(workDir, prefix[:lastSlash+1])
		base = prefix[lastSlash+1:]
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}

	out := make([]string, 0, len(entries))
	for _, ent := range entries {
		name := ent.Name()
		if base != "" && !strings.HasPrefix(name, base) {
			continue
		}
		if ent.IsDir() {
			name += "/"
		}
		if strings.Contains(prefix, "/") {
			name = filepath.Join(prefix[:strings.LastIndex(prefix, "/")+1], name)
		}
		out = append(out, name)
	}
	sort.Strings(out)
	return out
}

func filterPrefix(candidates []string, prefix string) []string {
	if prefix == "" {
		out := make([]string, len(candidates))
		copy(out, candidates)
		return out
	}
	out := make([]string, 0)
	for _, c := range candidates {
		if strings.HasPrefix(c, prefix) {
			out = append(out, c)
		}
	}
	return out
}
