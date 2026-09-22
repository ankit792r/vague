package display

import "vague/bonefire/window"

// TabPage is one tab's window layout.
type TabPage struct {
	Root           *window.Node
	ActiveWindowID uint64
	Label          string
}
