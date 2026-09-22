package text

import (
	"bytes"
	"sort"
	"time"
)

// SetInitial records the buffer bytes undo replay starts from.
func (u *UndoTree) SetInitial(b []byte) {
	u.initial = bytes.Clone(b)
}

func (u *UndoTree) Initial() []byte {
	return u.initial
}

func (u *UndoTree) findBySeq(seq int) *UndoNode {
	if seq <= 0 {
		return u.root
	}
	var found *UndoNode
	var walk func(n *UndoNode)
	walk = func(n *UndoNode) {
		if n.Seq == seq {
			found = n
			return
		}
		for _, c := range n.Children {
			walk(c)
		}
	}
	walk(u.root)
	return found
}

func (u *UndoTree) pathFromRoot(node *UndoNode) []*UndoNode {
	if node == nil || node == u.root {
		return nil
	}
	var path []*UndoNode
	for n := node; n != nil && n.Parent != nil; n = n.Parent {
		path = append([]*UndoNode{n}, path...)
	}
	return path
}

// ReplayToSeq rebuilds text at the given undo sequence and moves current there.
func (u *UndoTree) ReplayToSeq(seq int) ([]byte, Offset, bool) {
	target := u.findBySeq(seq)
	if target == nil {
		return nil, 0, false
	}
	tx := New(bytes.Clone(u.initial))
	for _, node := range u.pathFromRoot(target) {
		for _, e := range node.Edits {
			tx.Replace(e.At, e.OldEnd(), e.Inserted)
		}
	}
	u.current = target
	if len(target.Children) > 0 {
		u.current.Redo = target.Children[len(target.Children)-1]
	} else {
		u.current.Redo = nil
	}
	return tx.Bytes(), target.CursorAfter, true
}

func (u *UndoTree) collectNodes() []*UndoNode {
	var out []*UndoNode
	var walk func(n *UndoNode)
	walk = func(n *UndoNode) {
		if n != u.root && n.Seq > 0 {
			out = append(out, n)
		}
		for _, c := range n.Children {
			walk(c)
		}
	}
	walk(u.root)
	return out
}

// UndolistLines formats branching undo history for :undolist.
func (u *UndoTree) UndolistLines() []string {
	nodes := u.collectNodes()
	sort.Slice(nodes, func(i, j int) bool {
		return nodes[i].Seq < nodes[j].Seq
	})
	var lines []string
	for _, n := range nodes {
		mark := " "
		if n == u.current {
			mark = ">"
		}
		branch := ""
		if n.Parent != nil && len(n.Parent.Children) > 1 {
			branch = " *"
		}
		lines = append(lines, formatUndoLine(mark, n)+branch)
	}
	if len(lines) == 0 {
		lines = append(lines, ">0  (no changes yet)")
	}
	return lines
}

func formatUndoLine(mark string, n *UndoNode) string {
	edits := len(n.Edits)
	return mark + itoaUndo(n.Seq) + "  " + n.Time.Format("15:04:05") + "  " + itoaUndo(edits) + " edit(s)"
}

func itoaUndo(n int) string {
	if n == 0 {
		return "0"
	}
	var d [12]byte
	i := len(d)
	for n > 0 {
		i--
		d[i] = byte('0' + n%10)
		n /= 10
	}
	return string(d[i:])
}

// GoEarlierTime picks a sequence at or before current time minus d.
func (u *UndoTree) GoEarlierTime(d time.Duration) (seq int, ok bool) {
	if d <= 0 {
		return u.current.Seq, true
	}
	cutoff := u.current.Time.Add(-d)
	nodes := u.collectNodes()
	sort.Slice(nodes, func(i, j int) bool {
		return nodes[i].Time.Before(nodes[j].Time)
	})
	targetSeq := 0
	for _, n := range nodes {
		if !n.Time.After(cutoff) {
			targetSeq = n.Seq
		}
	}
	if targetSeq == u.current.Seq {
		return targetSeq, false
	}
	return targetSeq, true
}

// GoLaterTime picks a sequence after current but not newer than current+d.
func (u *UndoTree) GoLaterTime(d time.Duration) (seq int, ok bool) {
	if d <= 0 {
		return u.current.Seq, true
	}
	cutoff := u.current.Time.Add(d)
	nodes := u.collectNodes()
	sort.Slice(nodes, func(i, j int) bool {
		return nodes[i].Time.Before(nodes[j].Time)
	})
	targetSeq := u.current.Seq
	for _, n := range nodes {
		if n.Time.After(u.current.Time) && !n.Time.After(cutoff) {
			targetSeq = n.Seq
		}
	}
	if targetSeq == u.current.Seq {
		return targetSeq, false
	}
	return targetSeq, true
}

// UndoSteps moves current toward the root up to n times.
func (u *UndoTree) UndoSteps(n int) int {
	applied := 0
	for i := 0; i < n; i++ {
		if _, _, ok := u.Undo(); !ok {
			break
		}
		applied++
	}
	return applied
}

// RedoSteps moves current forward up to n times.
func (u *UndoTree) RedoSteps(n int) int {
	applied := 0
	for i := 0; i < n; i++ {
		if _, _, ok := u.Redo(); !ok {
			break
		}
		applied++
	}
	return applied
}
