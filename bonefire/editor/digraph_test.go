package editor

import (
	"strings"
	"testing"

	frame "vague/bonefire/frame"
	"vague/bonefire/text"
	"vague/bonefire/window"
)

func TestInsertDigraph(t *testing.T) {
	ed := NewEditor()
	buf := ed.Scratch("x")
	win := &window.Window{Cursor: buf.Text.AddMarker(0, text.GravityRight)}
	fr := &frame.Frame{Dirty: true, Width: 80, Height: 24}
	ed.SetMode(InsertMode)

	if err := ed.HandleInput(fr, win, buf, "<C-k>"); err != nil {
		t.Fatal(err)
	}
	if err := ed.HandleInput(fr, win, buf, "e"); err != nil {
		t.Fatal(err)
	}
	if err := ed.HandleInput(fr, win, buf, "'"); err != nil {
		t.Fatal(err)
	}

	got := string(buf.Text.Bytes())
	if !strings.Contains(got, "\u00e9") {
		t.Fatalf("expected e acute, got %q", got)
	}
}

func TestExDigraphsLists(t *testing.T) {
	ed := NewEditor()
	out := ed.exDigraphs("")
	if !strings.Contains(out, "e'") {
		t.Fatalf("missing digraph listing: %q", out)
	}
}
