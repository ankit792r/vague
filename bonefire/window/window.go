package window

// type Window struct {
// 	ID             uint64         `json:"id"`
// 	X              int            `json:"x"`
// 	Y              int            `json:"y"`
// 	Columns        int            `json:"columns"`
// 	Rows           int            `json:"rows"`
// 	Active         bool           `json:"active,omitempty"`
// 	Wrap           bool           `json:"wrap"`
// 	Number         bool           `json:"number,omitempty"`
// 	Buffer         RedrawBuffer   `json:"buffer"`
// 	Lines          []string       `json:"lines"`
// 	LineNumbers    []int          `json:"line_numbers,omitempty"`
// 	Cursor         CursorPos      `json:"cursor"`
// 	Mode           string         `json:"mode"`
// 	Position       BufferPosition `json:"position"`
// 	LineMarks      string         `json:"line_marks,omitempty"`
// 	RelativeNumber bool           `json:"relative_number,omitempty"`
// 	CursorLineRow  int            `json:"cursor_line_row,omitempty"`
// 	CursorColumn   int            `json:"cursor_column,omitempty"`
// 	CursorColumnOn bool           `json:"cursor_column_on,omitempty"`
// }

type Window struct {
	Id       uint64 `json:"id"`
	FrameId  uint64 `json:"frame_id"`
	BufferId uint64 `json:"buffer_id"`

	Lines []string `json:"lines"`

	Children []*Window `json:"children,omitempty"`
}

func NewWindow(windowId, frameId uint64) *Window {
	return &Window{
		Id:      windowId,
		FrameId: frameId,
	}
}
