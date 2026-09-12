package process

// Method names carried by requests and notifications.
const (
	// MethodExecute runs a registered editor command. Request.
	MethodExecute = "execute"

	// MethodInput delivers keys. Notification, client to server.
	//
	// The server's key handling arrives with the editing model; the
	// transport is defined here so the client does not have to change
	// shape later.
	MethodInput = "input"

	// MethodQuit tells a client its surface is gone and it should exit.
	// Notification, server to client.
	//
	// The editor decides this, not the client, because :q and :wq are
	// ordinary commands that happen to destroy a frame.
	MethodQuit = "quit"
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
//
// This is deliberately the whole RPC surface for editing. Key bindings, the
// ex command line, macros and remote callers all name a command and pass
// words, so adding a hardwired method per operation would just create a
// second way to do the same thing that the other callers cannot reach.
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
