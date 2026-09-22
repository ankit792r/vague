package text

import "time"

// Edit is one Replace, recorded so it can be undone or replayed.
type Edit struct {
	At       Offset
	Removed  []byte
	Inserted []byte
}

// Invert returns the edit that undoes e.
func (e Edit) Invert() Edit {
	return Edit{At: e.At, Removed: e.Inserted, Inserted: e.Removed}
}

// OldEnd is the end of the range this edit consumes when applied.
func (e Edit) OldEnd() Offset { return e.At + Offset(len(e.Removed)) }

// NewEnd is the end of the range this edit produces.
func (e Edit) NewEnd() Offset { return e.At + Offset(len(e.Inserted)) }

// UndoNode is one undoable change.
type UndoNode struct {
	Seq      int
	Parent   *UndoNode
	Children []*UndoNode
	Redo     *UndoNode
	Edits    []Edit

	CursorBefore Offset
	CursorAfter  Offset
	Time         time.Time
}

// UndoTree is a branching edit history.
type UndoTree struct {
	root    *UndoNode
	current *UndoNode
	nextSeq int

	group []Edit
	depth int
	before Offset

	initial []byte
}

func NewUndoTree() *UndoTree {
	root := &UndoNode{Time: time.Now()}
	return &UndoTree{root: root, current: root, nextSeq: 1}
}

func (u *UndoTree) Root() *UndoNode { return u.root }

func (u *UndoTree) Current() *UndoNode { return u.current }

func (u *UndoTree) Seq() int { return u.current.Seq }

func (u *UndoTree) Begin(cursor Offset) {
	u.depth++
	if u.depth == 1 {
		u.group = u.group[:0]
		u.before = cursor
	}
}

func (u *UndoTree) Record(e Edit) {
	if u.depth == 0 {
		u.Begin(e.At)
		u.group = append(u.group, e)
		u.Commit(e.NewEnd())
		return
	}
	u.group = append(u.group, e)
}

func (u *UndoTree) Commit(cursor Offset) {
	if u.depth == 0 {
		return
	}

	u.depth--
	if u.depth > 0 || len(u.group) == 0 {
		return
	}

	node := &UndoNode{
		Seq:          u.nextSeq,
		Parent:       u.current,
		Edits:        append([]Edit(nil), u.group...),
		CursorBefore: u.before,
		CursorAfter:  cursor,
		Time:         time.Now(),
	}
	u.nextSeq++

	u.current.Children = append(u.current.Children, node)
	u.current.Redo = node
	u.current = node
	u.group = u.group[:0]
}

func (u *UndoTree) InGroup() bool { return u.depth > 0 }

func (u *UndoTree) Undo() (edits []Edit, cursor Offset, ok bool) {
	node := u.current
	if node.Parent == nil {
		return nil, 0, false
	}

	edits = make([]Edit, 0, len(node.Edits))
	for i := len(node.Edits) - 1; i >= 0; i-- {
		edits = append(edits, node.Edits[i].Invert())
	}

	u.current = node.Parent
	u.current.Redo = node

	return edits, node.CursorBefore, true
}

func (u *UndoTree) Redo() (edits []Edit, cursor Offset, ok bool) {
	node := u.current.Redo
	if node == nil {
		if len(u.current.Children) == 0 {
			return nil, 0, false
		}
		node = u.current.Children[len(u.current.Children)-1]
	}

	u.current = node

	return append([]Edit(nil), node.Edits...), node.CursorAfter, true
}
