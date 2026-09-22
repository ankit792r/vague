package window

// ReplaceLeaf swaps a leaf for another node (typically a new split).
func (root *Node) ReplaceLeaf(windowID uint64, replacement *Node) bool {
	if root == nil {
		return false
	}
	if root.IsLeaf() {
		if root.WindowID != windowID {
			return false
		}
		weight := root.Weight
		*root = *replacement
		root.Weight = weight
		return true
	}
	for i, child := range root.Children {
		if child.IsLeaf() && child.WindowID == windowID {
			replacement.Weight = child.Weight
			root.Children[i] = replacement
			return true
		}
		if child.ReplaceLeaf(windowID, replacement) {
			return true
		}
	}
	return false
}

// SplitLeaf divides the window into two panes; vertical is :vsplit (side by side).
func SplitLeaf(root *Node, windowID, newWindowID uint64, vertical bool) bool {
	axis := Column
	if vertical {
		axis = Row
	}
	split := &Node{
		Axis: axis,
		Children: []*Node{
			{WindowID: windowID, Weight: 1},
			{WindowID: newWindowID, Weight: 1},
		},
	}
	return root.ReplaceLeaf(windowID, split)
}

// PaneRect is a window viewport in terminal cells.
type PaneRect struct {
	WindowID uint64
	X, Y     int
	Width    int
	Height   int
}

// LayoutPanes assigns cell rectangles from the split tree.
func LayoutPanes(root *Node, cols, rows int) []PaneRect {
	if root == nil || cols <= 0 || rows <= 0 {
		return nil
	}
	var out []PaneRect
	layoutNode(root, 0, 0, cols, rows, &out)
	return out
}

func layoutNode(n *Node, x, y, w, h int, out *[]PaneRect) {
	if n == nil || w <= 0 || h <= 0 {
		return
	}
	if n.IsLeaf() {
		*out = append(*out, PaneRect{WindowID: n.WindowID, X: x, Y: y, Width: w, Height: h})
		return
	}
	total := 0.0
	for _, c := range n.Children {
		total += c.Weight
	}
	if total <= 0 {
		total = float64(len(n.Children))
		for _, c := range n.Children {
			c.Weight = 1
		}
	}
	switch n.Axis {
	case Row:
		cursor := x
		for i, c := range n.Children {
			cw := int(float64(w) * c.Weight / total)
			if i == len(n.Children)-1 {
				cw = x + w - cursor
			}
			layoutNode(c, cursor, y, cw, h, out)
			cursor += cw
		}
	default:
		cursor := y
		for i, c := range n.Children {
			ch := int(float64(h) * c.Weight / total)
			if i == len(n.Children)-1 {
				ch = y + h - cursor
			}
			layoutNode(c, x, cursor, w, ch, out)
			cursor += ch
		}
	}
}

// OnlyWindow collapses the tree to a single leaf.
func OnlyWindow(root *Node, keepWindowID uint64) *Node {
	if root == nil {
		return nil
	}
	if root.IsLeaf() {
		if root.WindowID == keepWindowID {
			return root
		}
		root.WindowID = keepWindowID
		return root
	}
	leaf := root.find(keepWindowID)
	if leaf == nil {
		return root
	}
	leaf.Weight = 1
	return &Node{WindowID: keepWindowID, Weight: 1}
}

// NeighborWindow finds an adjacent window in the given direction (h/j/k/l).
func NeighborWindow(root *Node, activeID uint64, dir rune) uint64 {
	panes := LayoutPanes(root, 1000, 1000)
	if len(panes) <= 1 {
		return 0
	}
	var active *PaneRect
	byID := map[uint64]PaneRect{}
	for i := range panes {
		byID[panes[i].WindowID] = panes[i]
		if panes[i].WindowID == activeID {
			active = &panes[i]
		}
	}
	if active == nil {
		return 0
	}
	acx := active.X + active.Width/2
	acy := active.Y + active.Height/2
	bestID := uint64(0)
	bestDist := 1<<31 - 1
	for id, p := range byID {
		if id == activeID {
			continue
		}
		pcx := p.X + p.Width/2
		pcy := p.Y + p.Height/2
		switch dir {
		case 'h':
			if pcx >= acx {
				continue
			}
		case 'l':
			if pcx <= acx {
				continue
			}
		case 'j':
			if pcy <= acy {
				continue
			}
		case 'k':
			if pcy >= acy {
				continue
			}
		default:
			continue
		}
		dx := pcx - acx
		dy := pcy - acy
		d := dx*dx + dy*dy
		if d < bestDist {
			bestDist = d
			bestID = id
		}
	}
	return bestID
}

// AdjustWeight resizes the active pane relative to a neighbor in the split.
func AdjustWeight(root *Node, activeID uint64, dir rune, delta float64) bool {
	parent := root.parentOf(root.find(activeID))
	if parent == nil || parent.IsLeaf() {
		return false
	}
	idx := -1
	for i, c := range parent.Children {
		if c.find(activeID) != nil {
			idx = i
			break
		}
	}
	if idx < 0 {
		return false
	}
	neighbor := -1
	switch parent.Axis {
	case Row:
		if dir == 'l' && idx > 0 {
			neighbor = idx - 1
		}
		if dir == 'r' || dir == '>' {
			if idx+1 < len(parent.Children) {
				neighbor = idx + 1
			}
		}
		if dir == '<' {
			if idx > 0 {
				neighbor = idx - 1
			}
		}
		if dir == '+' || dir == '-' {
			if dir == '+' && idx+1 < len(parent.Children) {
				neighbor = idx + 1
			}
			if dir == '-' && idx > 0 {
				neighbor = idx - 1
			}
		}
	case Column:
		if dir == 'k' && idx > 0 {
			neighbor = idx - 1
		}
		if dir == 'j' || dir == '+' {
			if idx+1 < len(parent.Children) {
				neighbor = idx + 1
			}
		}
		if dir == '-' {
			if idx > 0 {
				neighbor = idx - 1
			}
		}
	}
	if neighbor < 0 {
		return false
	}
	a := parent.Children[idx]
	b := parent.Children[neighbor]
	if delta > 0 && b.Weight > delta {
		a.Weight += delta
		b.Weight -= delta
		return true
	}
	return false
}

// MaximizePane gives the active window all space along the split axis.
func MaximizePane(root *Node, activeID uint64, horizontal bool) bool {
	parent := root.parentOf(root.find(activeID))
	if parent == nil {
		return false
	}
	if horizontal {
		if parent.Axis != Row {
			return false
		}
	} else {
		if parent.Axis != Column {
			return false
		}
	}
	idx := -1
	for i, c := range parent.Children {
		if c.find(activeID) != nil {
			idx = i
			break
		}
	}
	if idx < 0 {
		return false
	}
	for i, c := range parent.Children {
		if i == idx {
			c.Weight = 100
		} else {
			c.Weight = 1
		}
	}
	return true
}
