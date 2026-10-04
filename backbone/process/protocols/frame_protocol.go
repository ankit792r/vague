package protocols

// FrameAttachParams carries optional startup files for a new client surface.
type FrameAttachParams struct {
	Files   []string `json:"files,omitempty"`
	WorkDir string   `json:"work_dir,omitempty"`
	Y       uint     `json:"y,omitempty"`
	X       uint     `json:"x,omitempty"`
}

// FrameAttachResult tells the client which server frame it owns.
type FrameAttachResult struct {
	SessionID uint64 `json:"session_id"`
	FrameID   uint64 `json:"frame_id"`
}

// FrameUpdateParams carries updated layout or dimension of frame
type FrameUpdateParams struct {
	Height uint `json:"height,omitempty"`
	Width  uint `json:"width,omitempty"`
	Y      uint `json:"y,omitempty"`
	X      uint `json:"x,omitempty"`
}

// FrameUpdateParams carries updated layout or dimension of frame
type FrameUpdateResult struct{}

type FrameDetachResult struct{}
