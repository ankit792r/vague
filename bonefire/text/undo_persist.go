package text

import (
	"encoding/json"
	"time"
)

type undoNodeDTO struct {
	Seq          int       `json:"seq"`
	ParentSeq    int       `json:"parent_seq"`
	Edits        []Edit    `json:"edits"`
	CursorBefore Offset    `json:"cursor_before"`
	CursorAfter  Offset    `json:"cursor_after"`
	Time         time.Time `json:"time"`
}

type undoTreeDTO struct {
	Initial    []byte        `json:"initial"`
	NextSeq    int           `json:"next_seq"`
	CurrentSeq int           `json:"current_seq"`
	Nodes      []undoNodeDTO `json:"nodes"`
}

// MarshalUndo encodes the undo tree for persistence.
func (u *UndoTree) MarshalUndo() ([]byte, error) {
	dto := undoTreeDTO{
		Initial:    bytesClone(u.initial),
		NextSeq:    u.nextSeq,
		CurrentSeq: u.current.Seq,
	}
	for _, n := range u.collectNodes() {
		parentSeq := 0
		if n.Parent != nil {
			parentSeq = n.Parent.Seq
		}
		dto.Nodes = append(dto.Nodes, undoNodeDTO{
			Seq:          n.Seq,
			ParentSeq:    parentSeq,
			Edits:        append([]Edit(nil), n.Edits...),
			CursorBefore: n.CursorBefore,
			CursorAfter:  n.CursorAfter,
			Time:         n.Time,
		})
	}
	return json.Marshal(dto)
}

// UnmarshalUndo restores an undo tree snapshot.
func UnmarshalUndo(data []byte) (*UndoTree, error) {
	var dto undoTreeDTO
	if err := json.Unmarshal(data, &dto); err != nil {
		return nil, err
	}
	u := NewUndoTree()
	u.initial = bytesClone(dto.Initial)
	u.nextSeq = dto.NextSeq
	if u.nextSeq < 1 {
		u.nextSeq = 1
	}
	bySeq := map[int]*UndoNode{0: u.root}
	for _, raw := range dto.Nodes {
		node := &UndoNode{
			Seq:          raw.Seq,
			Edits:        append([]Edit(nil), raw.Edits...),
			CursorBefore: raw.CursorBefore,
			CursorAfter:  raw.CursorAfter,
			Time:         raw.Time,
		}
		bySeq[node.Seq] = node
	}
	for _, raw := range dto.Nodes {
		node := bySeq[raw.Seq]
		parent := bySeq[raw.ParentSeq]
		if parent == nil {
			parent = u.root
		}
		node.Parent = parent
		parent.Children = append(parent.Children, node)
	}
	if cur := bySeq[dto.CurrentSeq]; cur != nil {
		u.current = cur
	} else {
		u.current = u.root
	}
	for _, n := range u.root.Children {
		linkRedoChain(n)
	}
	return u, nil
}

func linkRedoChain(n *UndoNode) {
	if len(n.Children) > 0 {
		n.Redo = n.Children[len(n.Children)-1]
		for _, c := range n.Children {
			linkRedoChain(c)
		}
	}
}

func bytesClone(b []byte) []byte {
	if len(b) == 0 {
		return nil
	}
	return append([]byte(nil), b...)
}
