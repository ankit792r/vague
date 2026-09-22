package process

// Method names carried by requests and notifications.
const (
	// MethodExecute runs a registered editor command. Request.
	MethodExecute = "execute"

	// MethodInput delivers keys. Notification, client to server.
	MethodInput = "input"

	// MethodQuit tells a client its surface is gone and it should exit.
	// Notification, server to client.
	MethodQuit = "quit"

	// MethodUiAttach binds a client webview to a server display frame. Request.
	MethodUiAttach = "ui_attach"

	// MethodUiDetach releases the client surface. Request.
	MethodUiDetach = "ui_detach"

	// MethodUiReady is sent when the webview knows its viewport size. Request.
	MethodUiReady = "ui_ready"

	// MethodRedraw carries editor state to the client. Notification, server to client.
	MethodRedraw = "redraw"
)

// InputParams carries keys in Vim notation, for example "ihello<Esc>".
type InputParams struct {
	Keys string `json:"keys"`
}

// QuitParams names the frame that went away.
type QuitParams struct {
	FrameID uint64 `json:"frame_id"`
}

// ExecuteParams invokes a registered editor command.
type ExecuteParams struct {
	Name  string   `json:"name"`
	Args  []string `json:"args,omitempty"`
	Bang  bool     `json:"bang,omitempty"`
	Count int      `json:"count,omitempty"`
}

// BufferInfo mirrors the editor's buffer summary on the wire.
type BufferInfo struct {
	ID       uint64 `json:"id"`
	Name     string `json:"name"`
	Path     string `json:"path"`
	Modified bool   `json:"modified"`
	ReadOnly bool   `json:"readonly"`
	Current  bool   `json:"current"`
	Lines    int    `json:"lines"`
	Bytes    int    `json:"bytes"`
}

// CommandInfo describes one registered command.
type CommandInfo struct {
	Name    string `json:"name"`
	Summary string `json:"summary"`
	MinArgs int    `json:"min_args"`
	MaxArgs int    `json:"max_args"`
}

// UiAttachParams carries optional startup files for a new client surface.
type UiAttachParams struct {
	Files   []string `json:"files,omitempty"`
	WorkDir string   `json:"work_dir,omitempty"`
}

// UiAttachResult tells the client which server frame it owns.
type UiAttachResult struct {
	SessionID uint64 `json:"session_id"`
	FrameID   uint64 `json:"frame_id"`
}

// UiReadyParams carries the webview viewport size.
type UiReadyParams struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

// UiReadyResult confirms the surface is sized and active.
type UiReadyResult struct {
	SessionID uint64 `json:"session_id"`
	FrameID   uint64 `json:"frame_id"`
}

// CursorPos is where the client should draw the caret.
type CursorPos struct {
	Row     int  `json:"row"`
	Column  int  `json:"column"`
	Visible bool `json:"visible"`
}

// SelectionPoint is a cell in the viewport selection highlight.
type SelectionPoint struct {
	Row    int `json:"row"`
	Column int `json:"column"`
}

// Selection describes highlighted text in the viewport.
type Selection struct {
	Visible  bool           `json:"visible"`
	Linewise bool           `json:"linewise,omitempty"`
	Start    SelectionPoint `json:"start"`
	End      SelectionPoint `json:"end"`
}

// RedrawBuffer identifies the buffer being drawn.
type RedrawBuffer struct {
	ID       uint64 `json:"id"`
	Name     string `json:"name"`
	Modified bool   `json:"modified,omitempty"`
}

// StatusEcho carries a short message for the client status line.
type StatusEcho struct {
	Message string `json:"message"`
	Kind    string `json:"kind,omitempty"`
}

// BufferPosition is the 1-based cursor position in the buffer for the status line.
type BufferPosition struct {
	Line   int `json:"line"`
	Column int `json:"column"`
}

// Redraw brings a client's picture up to date.
type Redraw struct {
	FrameID       uint64         `json:"frame_id"`
	Columns       int            `json:"columns"`
	Rows          int            `json:"rows"`
	Wrap          bool           `json:"wrap"`
	Number        bool           `json:"number,omitempty"`
	GutterColumns int            `json:"gutter_columns,omitempty"`
	Full          bool           `json:"full,omitempty"`
	Buffer        RedrawBuffer   `json:"buffer"`
	Lines         []string       `json:"lines"`
	LineNumbers   []int          `json:"line_numbers,omitempty"`
	Cursor        CursorPos      `json:"cursor"`
	Selection        *Selection   `json:"selection,omitempty"`
	SearchMatch      *Selection   `json:"search_match,omitempty"`
	SearchHighlights []Selection  `json:"search_highlights,omitempty"`
	Mode             string       `json:"mode"`
	Position         BufferPosition `json:"position"`
	LineMarks        string         `json:"line_marks,omitempty"`
	Echo             *StatusEcho    `json:"echo,omitempty"`
}
