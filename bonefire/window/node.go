package window

// Axis is how a split arranges its children.
type Axis int

const (
	// Row places children side by side, so the dividers between them are
	// vertical. This is what :vsplit and split-window-right produce.
	Row Axis = iota

	// Column stacks children top to bottom, so the dividers are
	// horizontal. This is what :split and split-window-below produce.
	Column
)

type Node struct {
	// WindowID is non-zero on a leaf and zero on a split.
	WindowID uint64 `json:"window_id,omitempty"`

	Axis     Axis    `json:"axis,omitempty"`
	Children []*Node `json:"children,omitempty"`

	// Weight is this node's share of its parent's extent along the parent's
	// axis. Sizes are relative so that resizing the frame keeps the
	// proportions the user set up.
	Weight float64 `json:"weight"`
}

func (w *Window) LeafNode() *Node {
	return &Node{WindowID: w.Id, Weight: 1}
}

// IsLeaf reports whether the node holds a window.
func (n *Node) IsLeaf() bool { return n.WindowID != 0 }

// Leaves collects every window in tree order, which is the order
// other-window cycles through.
func (n *Node) Leaves() []*Node {
	if n == nil {
		return nil
	}
	if n.IsLeaf() {
		return []*Node{n}
	}

	var out []*Node
	for _, child := range n.Children {
		out = append(out, child.Leaves()...)
	}

	return out
}

// find returns the leaf holding a window.
func (n *Node) find(windowID uint64) *Node {
	if n == nil {
		return nil
	}
	if n.IsLeaf() {
		if n.WindowID == windowID {
			return n
		}
		return nil
	}

	for _, child := range n.Children {
		if found := child.find(windowID); found != nil {
			return found
		}
	}

	return nil
}

// parentOf returns the split containing a node, or nil for the root.
func (n *Node) parentOf(target *Node) *Node {
	if n == nil || n.IsLeaf() {
		return nil
	}

	for _, child := range n.Children {
		if child == target {
			return n
		}
		if found := child.parentOf(target); found != nil {
			return found
		}
	}

	return nil
}

// remove deletes a leaf, collapsing splits that end up with a single child
// so the tree never accumulates pointless nesting.
func removeLeaf(root *Node, leaf *Node) *Node {
	parent := root.parentOf(leaf)
	if parent == nil {
		// Removing the only window leaves an empty frame.
		return nil
	}

	index := indexOf(parent.Children, leaf)
	parent.Children = append(parent.Children[:index], parent.Children[index+1:]...)

	// Give the freed space to the remaining siblings.
	redistribute(parent.Children, leaf.Weight)

	if len(parent.Children) == 1 {
		survivor := parent.Children[0]
		weight := parent.Weight

		*parent = *survivor
		parent.Weight = weight
	}

	return root
}

// redistribute hands a departed child's weight to its siblings in
// proportion to what they already had.
func redistribute(siblings []*Node, freed float64) {
	if len(siblings) == 0 || freed <= 0 {
		return
	}

	total := 0.0
	for _, sibling := range siblings {
		total += sibling.Weight
	}

	if total <= 0 {
		for _, sibling := range siblings {
			sibling.Weight += freed / float64(len(siblings))
		}
		return
	}

	for _, sibling := range siblings {
		sibling.Weight += freed * sibling.Weight / total
	}
}

func indexOf(nodes []*Node, target *Node) int {
	for i, node := range nodes {
		if node == target {
			return i
		}
	}
	return -1
}

// balance resets every weight, making sibling windows equal size.
func balance(n *Node) {
	if n == nil || n.IsLeaf() {
		return
	}

	for _, child := range n.Children {
		child.Weight = 1
		balance(child)
	}
}

func insertNode(nodes []*Node, index int, node *Node) []*Node {
	if index < 0 {
		index = 0
	}
	if index > len(nodes) {
		index = len(nodes)
	}

	nodes = append(nodes, nil)
	copy(nodes[index+1:], nodes[index:])
	nodes[index] = node

	return nodes
}
