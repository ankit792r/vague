package window

// WindowOptions are settings that belong to a viewport rather than to the
// text, so two windows on one buffer can differ.
type WindowOptions struct {
	// Number shows line numbers.
	Number bool

	// Wrap soft-wraps long lines instead of scrolling horizontally.
	Wrap bool

	// ScrollOff keeps this many lines visible above and below the cursor.
	ScrollOff int
}

func DefaultWindowOptions() WindowOptions {
	return WindowOptions{Number: false, Wrap: true, ScrollOff: 0}
}

type Window struct {
	Id      uint64
	FrameId uint64

	BufferId uint64

	Axis Axis

	Children []*Window `json:"children,omitempty"`

	// CursorLine and CursorCol are the cursor in buffer coordinates.
	// Col is a rune index within the line.
	CursorLine int
	CursorCol  int

	// DesiredCol is the rune column j/k aim for after vertical movement.
	DesiredCol int

	WindowOptions WindowOptions
}

func NewWindow(windowId, frameId uint64) *Window {
	return &Window{
		Id:            windowId,
		FrameId:       frameId,
		WindowOptions: DefaultWindowOptions(),
	}
}
