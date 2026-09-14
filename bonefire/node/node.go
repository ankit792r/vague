package node

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
