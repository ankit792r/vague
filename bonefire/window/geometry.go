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

func (a Axis) String() string {
	if a == Column {
		return "column"
	}
	return "row"
}
