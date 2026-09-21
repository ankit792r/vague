package editor_test

import (
	"testing"

	"vague/bonefire/workspace"
)

func TestRepeatDeleteChar(t *testing.T) {
	t.Parallel()

	ws := workspace.New()
	buf := ws.Editor.Scratch("*scratch*")
	buf.Text.SetBytes([]byte("abcd\n"))

	frame, err := ws.NewFrame(80, 10, "")
	if err != nil {
		t.Fatal(err)
	}

	if err := ws.HandleInput(frame.ID, "x"); err != nil {
		t.Fatal(err)
	}
	if string(buf.Text.Bytes()) != "bcd\n" {
		t.Fatalf("after x got %q", buf.Text.Bytes())
	}

	if err := ws.HandleInput(frame.ID, "."); err != nil {
		t.Fatal(err)
	}
	if string(buf.Text.Bytes()) != "cd\n" {
		t.Fatalf("after . got %q", buf.Text.Bytes())
	}
}

func TestRepeatOperatorMotion(t *testing.T) {
	t.Parallel()

	ws := workspace.New()
	buf := ws.Editor.Scratch("*scratch*")
	buf.Text.SetBytes([]byte("one\ntwo\nthree\n"))

	frame, err := ws.NewFrame(80, 10, "")
	if err != nil {
		t.Fatal(err)
	}

	if err := ws.HandleInput(frame.ID, "d"); err != nil {
		t.Fatal(err)
	}
	if err := ws.HandleInput(frame.ID, "d"); err != nil {
		t.Fatal(err)
	}
	if string(buf.Text.Bytes()) != "two\nthree\n" {
		t.Fatalf("after dd got %q", buf.Text.Bytes())
	}

	if err := ws.HandleInput(frame.ID, "."); err != nil {
		t.Fatal(err)
	}
	if string(buf.Text.Bytes()) != "three\n" {
		t.Fatalf("after . got %q", buf.Text.Bytes())
	}
}

func TestRedrawIncludesPosition(t *testing.T) {
	t.Parallel()

	ws := workspace.New()
	buf := ws.Editor.Scratch("*scratch*")
	buf.Text.SetBytes([]byte("ab\n"))

	frame, err := ws.NewFrame(80, 10, "")
	if err != nil {
		t.Fatal(err)
	}

	if err := ws.HandleInput(frame.ID, "l"); err != nil {
		t.Fatal(err)
	}

	redraw, ok := ws.RenderRedraw(frame.ID)
	if !ok {
		t.Fatal("expected redraw")
	}
	if redraw.Position.Line != 1 || redraw.Position.Column != 2 {
		t.Fatalf("position = %d:%d, want 1:2", redraw.Position.Line, redraw.Position.Column)
	}
}
