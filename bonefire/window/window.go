package window

import "vague/bonefire/text"

// WindowOptions are settings that belong to a viewport rather than to the
// text, so two windows on one buffer can differ.
type WindowOptions struct {
	Number    bool
	Wrap      bool
	ScrollOff int

	Display DisplayOpts
}

func DefaultWindowOptions() WindowOptions {
	return WindowOptions{Number: false, Wrap: true, ScrollOff: 0, Display: DefaultDisplayOpts()}
}

type Window struct {
	Id      uint64
	FrameId uint64
	BufferId uint64

	Axis     Axis
	Children []*Window `json:"children,omitempty"`

	// Cursor is a marker into the buffer text so it survives edits.
	Cursor *text.Marker

	// DesiredCol is the byte column j/k aim for after vertical movement.
	DesiredCol int

	// TopLine is the index of the first visible visual row in the layout.
	TopLine int

	LocalWorkDir string

	WindowOptions WindowOptions
}

func NewWindow(windowId, frameId uint64) *Window {
	return &Window{
		Id:            windowId,
		FrameId:       frameId,
		WindowOptions: DefaultWindowOptions(),
	}
}
