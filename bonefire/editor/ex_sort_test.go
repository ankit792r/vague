package editor

import (
	"testing"

	frame "vague/bonefire/frame"
	"vague/bonefire/text"
	"vague/bonefire/window"
)

func TestExSortRange(t *testing.T) {
	ed := NewEditor()
	buf := ed.Scratch("sort-test")
	buf.Text.SetBytes([]byte("c\nb\na"))
	f := &frame.Frame{Dirty: true}
	win := &window.Window{Cursor: buf.Text.AddMarker(0, text.GravityRight)}

	if err := ed.exSort(f, win, buf, "", exRange{wholeBuf: true}); err != nil {
		t.Fatal(err)
	}
	if got := string(buf.Text.Bytes()); got != "a\nb\nc" {
		t.Fatalf("got %q", got)
	}
}

func TestExUniqRange(t *testing.T) {
	ed := NewEditor()
	buf := ed.Scratch("uniq-test")
	buf.Text.SetBytes([]byte("a\na\nb\nb\na\n"))
	f := &frame.Frame{Dirty: true}
	win := &window.Window{Cursor: buf.Text.AddMarker(0, text.GravityRight)}

	if err := ed.exUniq(f, win, buf, exRange{wholeBuf: true}); err != nil {
		t.Fatal(err)
	}
	if got := string(buf.Text.Bytes()); got != "a\nb\na\n" {
		t.Fatalf("got %q", got)
	}
}
