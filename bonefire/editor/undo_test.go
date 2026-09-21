package editor_test

import (
	"testing"

	"vague/bonefire/workspace"
)

func TestUndoInsertSession(t *testing.T) {
	t.Parallel()

	ws := workspace.New()
	ed := ws.Editor
	buf := ed.Scratch("*scratch*")
	buf.Text.SetBytes([]byte("hi"))

	frame, err := ws.NewFrame(80, 10, "")
	if err != nil {
		t.Fatal(err)
	}

	if err := ws.HandleInput(frame.ID, "l"); err != nil {
		t.Fatal(err)
	}

	if err := ws.HandleInput(frame.ID, "l"); err != nil {
		t.Fatal(err)
	}

	if err := ws.HandleInput(frame.ID, "a"); err != nil {
		t.Fatal(err)
	}

	if err := ws.HandleInput(frame.ID, "!"); err != nil {
		t.Fatal(err)
	}

	if err := ws.HandleInput(frame.ID, "x"); err != nil {
		t.Fatal(err)
	}

	if err := ws.HandleInput(frame.ID, "<Esc>"); err != nil {
		t.Fatal(err)
	}

	if string(buf.Text.Bytes()) != "hi!x" {
		t.Fatalf("before undo got %q", buf.Text.Bytes())
	}

	if err := ws.HandleInput(frame.ID, "u"); err != nil {
		t.Fatal(err)
	}

	if string(buf.Text.Bytes()) != "hi" {
		t.Fatalf("after undo got %q, want hi", buf.Text.Bytes())
	}
}

func TestUndoRedoRoundTrip(t *testing.T) {
	t.Parallel()

	ws := workspace.New()
	ed := ws.Editor
	buf := ed.Scratch("*scratch*")
	buf.Text.SetBytes([]byte("ab"))

	frame, err := ws.NewFrame(80, 10, "")
	if err != nil {
		t.Fatal(err)
	}

	if err := ws.HandleInput(frame.ID, "l"); err != nil {
		t.Fatal(err)
	}

	if err := ws.HandleInput(frame.ID, "l"); err != nil {
		t.Fatal(err)
	}

	if err := ws.HandleInput(frame.ID, "a"); err != nil {
		t.Fatal(err)
	}

	if err := ws.HandleInput(frame.ID, "X"); err != nil {
		t.Fatal(err)
	}

	if err := ws.HandleInput(frame.ID, "<Esc>"); err != nil {
		t.Fatal(err)
	}

	if err := ws.HandleInput(frame.ID, "u"); err != nil {
		t.Fatal(err)
	}

	if string(buf.Text.Bytes()) != "ab" {
		t.Fatalf("after undo got %q", buf.Text.Bytes())
	}

	if err := ws.HandleInput(frame.ID, "<C-r>"); err != nil {
		t.Fatal(err)
	}

	if string(buf.Text.Bytes()) != "abX" {
		t.Fatalf("after redo got %q", buf.Text.Bytes())
	}
}
