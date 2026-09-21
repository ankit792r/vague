package editor

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"vague/bonefire/buffer"
)

var ErrBufferNotFound = errors.New("Buffer not found")

func (e *Editor) sortedBuffers() []*buffer.Buffer {
	list := make([]*buffer.Buffer, 0, len(e.Buffers))
	for _, buf := range e.Buffers {
		list = append(list, buf)
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i].ID < list[j].ID
	})

	return list
}

func (e *Editor) NextBuffer(fromID uint64) *buffer.Buffer {
	list := e.sortedBuffers()
	if len(list) == 0 {
		return nil
	}

	idx := 0
	for i, buf := range list {
		if buf.ID == fromID {
			idx = (i + 1) % len(list)
			return list[idx]
		}
	}

	return list[0]
}

func (e *Editor) PrevBuffer(fromID uint64) *buffer.Buffer {
	list := e.sortedBuffers()
	if len(list) == 0 {
		return nil
	}

	idx := len(list) - 1
	for i, buf := range list {
		if buf.ID == fromID {
			idx = i - 1
			if idx < 0 {
				idx = len(list) - 1
			}
			return list[idx]
		}
	}

	return list[len(list)-1]
}

func (e *Editor) ResolveBuffer(spec string) (*buffer.Buffer, error) {
	spec = strings.TrimSpace(spec)
	if spec == "" {
		return nil, ErrBufferNotFound
	}

	if n, err := strconv.Atoi(spec); err == nil {
		list := e.sortedBuffers()
		if n < 1 || n > len(list) {
			return nil, ErrBufferNotFound
		}
		return list[n-1], nil
	}

	if buf := e.FindBuffer(spec); buf != nil {
		return buf, nil
	}

	for _, buf := range e.Buffers {
		if buf.Name == spec {
			return buf, nil
		}
	}

	lower := strings.ToLower(spec)
	for _, buf := range e.Buffers {
		if strings.EqualFold(buf.Name, spec) {
			return buf, nil
		}
		if strings.HasSuffix(strings.ToLower(buf.Path), lower) {
			return buf, nil
		}
	}

	return nil, ErrBufferNotFound
}

func (e *Editor) FormatBufferList(currentID uint64) string {
	list := e.sortedBuffers()
	if len(list) == 0 {
		return "No buffers"
	}

	var b strings.Builder
	for i, buf := range list {
		if i > 0 {
			b.WriteByte('\n')
		}

		mark := " "
		if buf.ID == currentID {
			mark = "%"
		}
		mod := " "
		if buf.Modified() {
			mod = "+"
		}

		b.WriteString(fmt.Sprintf("  %d %s%s %q", i+1, mark, mod, buf.Name))
	}

	return b.String()
}
