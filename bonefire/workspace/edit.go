package workspace

import (
	"strings"

	"vague/bonefire/buffer"
)

// ParseEditArgs splits `:e ++binary path` style arguments.
func ParseEditArgs(args string) (path string, binary bool, force bool) {
	args = strings.TrimSpace(args)
	for args != "" {
		if strings.HasPrefix(args, "++") {
			end := strings.IndexByte(args, ' ')
			token := args
			if end >= 0 {
				token = args[:end]
				args = strings.TrimSpace(args[end+1:])
			} else {
				args = ""
			}
			if token == "++binary" {
				binary = true
			}
			continue
		}
		if args == "!" || strings.HasPrefix(args, "! ") {
			force = true
			args = strings.TrimSpace(strings.TrimPrefix(args, "!"))
			continue
		}
		path = args
		break
	}
	return path, binary, force
}

func (w *Workspace) EditFile(frameID uint64, args string) (*buffer.Buffer, error) {
	path, binary, force := ParseEditArgs(args)
	if path == "" {
		_, err := w.ReloadCurrentBuffer(frameID, force)
		return nil, err
	}
	buf, err := w.OpenFile(frameID, path, force)
	if err != nil {
		return nil, err
	}
	if binary && buf != nil {
		buf.Binary = true
	}
	return buf, err
}
