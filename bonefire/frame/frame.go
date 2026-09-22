package display

import "vague/bonefire/window"

const (
	DefaultWidth  = 80
	DefaultHeight = 24
)

// StatusEcho is user-visible feedback shown on the command line until the next key.
type StatusEcho struct {
	Message string
	Kind    string
}

// Frame is one client surface: a tree of windows showing buffers.
type Frame struct {
	ID     uint64
	Width  int
	Height int
	Root   *window.Node

	ActiveWindowID uint64
	Tabs           []TabPage
	ActiveTab      int

	WorkDir        string
	Dirty          bool
	Echo           StatusEcho
}
