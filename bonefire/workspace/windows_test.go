package workspace

import (
	"testing"

	"vague/bonefire/window"
)

func TestSplitAndOnly(t *testing.T) {
	ws := New()
	frame, err := ws.NewFrame(80, 24, "/tmp")
	if err != nil {
		t.Fatal(err)
	}
	if err := ws.SplitWindow(frame.ID, ""); err != nil {
		t.Fatal(err)
	}
	leaves := frame.Root.Leaves()
	if len(leaves) != 2 {
		t.Fatalf("want 2 windows, got %d", len(leaves))
	}
	if err := ws.OnlyWindow(frame.ID); err != nil {
		t.Fatal(err)
	}
	leaves = frame.Root.Leaves()
	if len(leaves) != 1 {
		t.Fatalf("want 1 window after only, got %d", len(leaves))
	}
}

func TestLayoutPanes(t *testing.T) {
	root := &window.Node{
		Axis: window.Column,
		Children: []*window.Node{
			{WindowID: 1, Weight: 1},
			{WindowID: 2, Weight: 1},
		},
	}
	panes := window.LayoutPanes(root, 80, 24)
	if len(panes) != 2 {
		t.Fatalf("want 2 panes, got %d", len(panes))
	}
	if panes[0].Height+panes[1].Height != 24 {
		t.Fatalf("heights %d+%d != 24", panes[0].Height, panes[1].Height)
	}
}
