package editor

import (
	"testing"

	frame "vague/bonefire/frame"
	"vague/bonefire/text"
	"vague/bonefire/window"
)

func TestIncsearchConfirmSetsPattern(t *testing.T) {
	ed := NewEditor()
	buf := ed.Scratch("x")
	_, _ = buf.Insert(0, []byte("abc abc"))
	win := &window.Window{Cursor: buf.Text.AddMarker(0, text.GravityRight)}
	fm := &frame.Frame{Width: 80, Height: 24, Dirty: true}

	ed.BeginIncsearch(win, buf, true)
	_ = ed.PreviewIncsearch(fm, win, buf, "abc")
	if err := ed.Search(fm, win, buf, "abc", true); err != nil {
		t.Fatal(err)
	}
	if ed.searchPattern != "abc" {
		t.Fatalf("pattern = %q", ed.searchPattern)
	}
	if ed.incsearchActive {
		t.Fatal("incsearch still active after confirm")
	}
}
