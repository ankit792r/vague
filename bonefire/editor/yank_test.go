package editor_test

import (
	"testing"

	"vague/bonefire/workspace"
)

func TestYankLineAndPasteBelow(t *testing.T) {
	t.Parallel()

	ws := workspace.New()
	ed := ws.Editor
	buf := ed.Scratch("*scratch*")
	buf.Text.SetBytes([]byte("alpha\nbeta\ngamma\n"))

	frame, err := ws.NewFrame(80, 10, "")
	if err != nil {
		t.Fatal(err)
	}

	if err := ws.HandleInput(frame.ID, "j"); err != nil {
		t.Fatal(err)
	}

	if err := ws.HandleInput(frame.ID, "y"); err != nil {
		t.Fatal(err)
	}
	if err := ws.HandleInput(frame.ID, "y"); err != nil {
		t.Fatal(err)
	}

	if string(buf.Text.Bytes()) != "alpha\nbeta\ngamma\n" {
		t.Fatalf("yy should not delete, got %q", buf.Text.Bytes())
	}

	if err := ws.HandleInput(frame.ID, "p"); err != nil {
		t.Fatal(err)
	}

	want := "alpha\nbeta\nbeta\ngamma\n"
	if string(buf.Text.Bytes()) != want {
		t.Fatalf("after p got %q, want %q", buf.Text.Bytes(), want)
	}
}

func TestYankLineAndPasteAbove(t *testing.T) {
	t.Parallel()

	ws := workspace.New()
	ed := ws.Editor
	buf := ed.Scratch("*scratch*")
	buf.Text.SetBytes([]byte("one\ntwo\nthree\n"))

	frame, err := ws.NewFrame(80, 10, "")
	if err != nil {
		t.Fatal(err)
	}

	if err := ws.HandleInput(frame.ID, "j"); err != nil {
		t.Fatal(err)
	}
	if err := ws.HandleInput(frame.ID, "y"); err != nil {
		t.Fatal(err)
	}
	if err := ws.HandleInput(frame.ID, "y"); err != nil {
		t.Fatal(err)
	}
	if err := ws.HandleInput(frame.ID, "P"); err != nil {
		t.Fatal(err)
	}

	want := "one\ntwo\ntwo\nthree\n"
	if string(buf.Text.Bytes()) != want {
		t.Fatalf("after P got %q, want %q", buf.Text.Bytes(), want)
	}
}

func TestPasteUndo(t *testing.T) {
	t.Parallel()

	ws := workspace.New()
	ed := ws.Editor
	buf := ed.Scratch("*scratch*")
	buf.Text.SetBytes([]byte("aa\nbb\n"))

	frame, err := ws.NewFrame(80, 10, "")
	if err != nil {
		t.Fatal(err)
	}

	if err := ws.HandleInput(frame.ID, "y"); err != nil {
		t.Fatal(err)
	}
	if err := ws.HandleInput(frame.ID, "y"); err != nil {
		t.Fatal(err)
	}
	if err := ws.HandleInput(frame.ID, "p"); err != nil {
		t.Fatal(err)
	}

	if err := ws.HandleInput(frame.ID, "u"); err != nil {
		t.Fatal(err)
	}

	if string(buf.Text.Bytes()) != "aa\nbb\n" {
		t.Fatalf("after undo got %q", buf.Text.Bytes())
	}
}

func TestPasteEmptyRegisterNoOp(t *testing.T) {
	t.Parallel()

	ws := workspace.New()
	ed := ws.Editor
	buf := ed.Scratch("*scratch*")
	buf.Text.SetBytes([]byte("hi\n"))

	frame, err := ws.NewFrame(80, 10, "")
	if err != nil {
		t.Fatal(err)
	}

	if err := ws.HandleInput(frame.ID, "p"); err != nil {
		t.Fatal(err)
	}

	if string(buf.Text.Bytes()) != "hi\n" {
		t.Fatalf("empty register paste changed buffer: %q", buf.Text.Bytes())
	}
}
