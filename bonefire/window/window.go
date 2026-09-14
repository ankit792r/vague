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

	// cursor is a marker rather than a bare offset so it stays correct when
	// the buffer is edited elsewhere, including by another window showing
	// the same buffer.
	// cursor *text.Marker

	WindowOptions WindowOptions
}
