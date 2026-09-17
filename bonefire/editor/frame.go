package editor

import "vague/bonefire/window"

const (
	defaultFrameWidth  = 80
	defaultFrameHeight = 24
)

// Frame is one client surface: a tree of windows showing buffers.
type Frame struct {
	ID     uint64
	Width  int
	Height int
	Root   *window.Node

	ActiveWindowID uint64
	dirty          bool
}
