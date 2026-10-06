package session

// CursorPos is where the client should draw the caret.
type CursorPos struct {
	Row     int  `json:"row"`
	Column  int  `json:"column"`
	Visible bool `json:"visible"`
}

// BufferPosition is the 1-based cursor position in the buffer for the status line.
type BufferPosition struct {
	Line   int `json:"line"`
	Column int `json:"column"`
}

// RedrawBuffer identifies the buffer being drawn.
type RedrawBuffer struct {
	ID       uint64 `json:"id"`
	Name     string `json:"name"`
	Modified bool   `json:"modified,omitempty"`
}

type FrameState struct {
	FrameID    uint64   `json:"frame_id"`
	Windows    []Window `json:"windows,omitempty"`
	Theme      string   `json:"theme,omitempty"`
	ShowCmd    string   `json:"showcmd,omitempty"`
	StatusLine string   `json:"statusline,omitempty"`
}

// RedrawPane is one window in a split layout.
type Window struct {
	WindowID       uint64         `json:"window_id"`
	X              int            `json:"x"`
	Y              int            `json:"y"`
	Columns        int            `json:"columns"`
	Rows           int            `json:"rows"`
	Active         bool           `json:"active,omitempty"`
	Wrap           bool           `json:"wrap"`
	Number         bool           `json:"number,omitempty"`
	Buffer         RedrawBuffer   `json:"buffer"`
	Lines          []string       `json:"lines"`
	LineNumbers    []int          `json:"line_numbers,omitempty"`
	Cursor         CursorPos      `json:"cursor"`
	Mode           string         `json:"mode"`
	Position       BufferPosition `json:"position"`
	LineMarks      string         `json:"line_marks,omitempty"`
	RelativeNumber bool           `json:"relative_number,omitempty"`
	CursorLineRow  int            `json:"cursor_line_row,omitempty"`
	CursorColumn   int            `json:"cursor_column,omitempty"`
	CursorColumnOn bool           `json:"cursor_column_on,omitempty"`
}

// Creating empty state
func EmptyFrameState() *FrameState {
	return &FrameState{}
}
