package editor

import (
	"testing"

	frame "vague/bonefire/frame"
	"vague/bonefire/text"
	"vague/bonefire/window"
)

func TestSetLocalMarkAndJump(t *testing.T) {
	ed := NewEditor()
	buf := ed.Scratch("x")
	win := &window.Window{Cursor: buf.Text.AddMarker(0, text.GravityRight)}
	fr := &frame.Frame{Dirty: true, Width: 80, Height: 24}

	if err := ed.HandleInput(fr, win, buf, "m"); err != nil {
		t.Fatal(err)
	}
	if err := ed.HandleInput(fr, win, buf, "a"); err != nil {
		t.Fatal(err)
	}

	setWindowCursor(buf, win, buf.Text.Len())
	if err := ed.HandleInput(fr, win, buf, "`"); err != nil {
		t.Fatal(err)
	}
	if err := ed.HandleInput(fr, win, buf, "a"); err != nil {
		t.Fatal(err)
	}

	if windowCursor(win) != 0 {
		t.Fatalf("expected jump to mark a at 0, got %d", windowCursor(win))
	}
}

func TestLineMarkJump(t *testing.T) {
	ed := NewEditor()
	buf := ed.Scratch("x")
	_, _ = buf.Insert(0, []byte("  hello\n"))
	win := &window.Window{Cursor: buf.Text.AddMarker(7, text.GravityRight)}
	fr := &frame.Frame{Dirty: true, Width: 80, Height: 24}

	setWindowCursor(buf, win, 7)
	if err := ed.HandleInput(fr, win, buf, "m"); err != nil {
		t.Fatal(err)
	}
	if err := ed.HandleInput(fr, win, buf, "b"); err != nil {
		t.Fatal(err)
	}

	setWindowCursor(buf, win, 0)
	if err := ed.HandleInput(fr, win, buf, "'"); err != nil {
		t.Fatal(err)
	}
	if err := ed.HandleInput(fr, win, buf, "b"); err != nil {
		t.Fatal(err)
	}

	want := firstNonBlankOnLine(buf.Text, 0)
	if windowCursor(win) != want {
		t.Fatalf("line jump: got %d want %d", windowCursor(win), want)
	}
}

func TestMarksOnLineLabel(t *testing.T) {
	ed := NewEditor()
	buf := ed.Scratch("x")
	win := &window.Window{Cursor: buf.Text.AddMarker(0, text.GravityRight)}
	fr := &frame.Frame{Dirty: true, Width: 80, Height: 24}

	_ = ed.HandleInput(fr, win, buf, "m")
	_ = ed.HandleInput(fr, win, buf, "a")
	_ = ed.HandleInput(fr, win, buf, "m")
	_ = ed.HandleInput(fr, win, buf, "b")

	label := ed.marksOnLineLabel(buf, 0)
	if label != "ab" {
		t.Fatalf("marks label: got %q want ab", label)
	}
}
