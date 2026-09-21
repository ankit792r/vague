package editor_test

import (
	"testing"

	"vague/bonefire/editor"
	"vague/bonefire/workspace"
)

func TestDeleteWord(t *testing.T) {
	t.Parallel()

	ws := workspace.New()
	buf := ws.Editor.Scratch("*scratch*")
	buf.Text.SetBytes([]byte("hello world\n"))

	frame, err := ws.NewFrame(80, 10, "")
	if err != nil {
		t.Fatal(err)
	}

	if err := ws.HandleInput(frame.ID, "d"); err != nil {
		t.Fatal(err)
	}
	if err := ws.HandleInput(frame.ID, "w"); err != nil {
		t.Fatal(err)
	}

	if string(buf.Text.Bytes()) != "world\n" {
		t.Fatalf("after dw got %q", buf.Text.Bytes())
	}
}

func TestDeleteToEOLWithD(t *testing.T) {
	t.Parallel()

	ws := workspace.New()
	buf := ws.Editor.Scratch("*scratch*")
	buf.Text.SetBytes([]byte("hello world\n"))

	frame, err := ws.NewFrame(80, 10, "")
	if err != nil {
		t.Fatal(err)
	}

	if err := ws.HandleInput(frame.ID, "0"); err != nil {
		t.Fatal(err)
	}
	if err := ws.HandleInput(frame.ID, "D"); err != nil {
		t.Fatal(err)
	}

	if string(buf.Text.Bytes()) != "\n" {
		t.Fatalf("after D got %q", buf.Text.Bytes())
	}
}

func TestDeleteLineStillWorks(t *testing.T) {
	t.Parallel()

	ws := workspace.New()
	buf := ws.Editor.Scratch("*scratch*")
	buf.Text.SetBytes([]byte("one\ntwo\n"))

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

	if string(buf.Text.Bytes()) != "two\n" {
		t.Fatalf("after dd got %q", buf.Text.Bytes())
	}
}

func TestYankWordThenPaste(t *testing.T) {
	t.Parallel()

	ws := workspace.New()
	buf := ws.Editor.Scratch("*scratch*")
	buf.Text.SetBytes([]byte("hello world\n"))

	frame, err := ws.NewFrame(80, 10, "")
	if err != nil {
		t.Fatal(err)
	}

	if err := ws.HandleInput(frame.ID, "y"); err != nil {
		t.Fatal(err)
	}
	if err := ws.HandleInput(frame.ID, "w"); err != nil {
		t.Fatal(err)
	}
	if err := ws.HandleInput(frame.ID, "$"); err != nil {
		t.Fatal(err)
	}
	if err := ws.HandleInput(frame.ID, "p"); err != nil {
		t.Fatal(err)
	}

	if string(buf.Text.Bytes()) != "hello worldhello \n" {
		t.Fatalf("after yw+p got %q", buf.Text.Bytes())
	}
}

func TestChangeWord(t *testing.T) {
	t.Parallel()

	ws := workspace.New()
	buf := ws.Editor.Scratch("*scratch*")
	buf.Text.SetBytes([]byte("hello world\n"))

	frame, err := ws.NewFrame(80, 10, "")
	if err != nil {
		t.Fatal(err)
	}

	if err := ws.HandleInput(frame.ID, "c"); err != nil {
		t.Fatal(err)
	}
	if err := ws.HandleInput(frame.ID, "w"); err != nil {
		t.Fatal(err)
	}
	if ws.Editor.Mode != editor.InsertMode {
		t.Fatal("expected insert mode after cw")
	}

	if err := ws.HandleInput(frame.ID, "hi"); err != nil {
		t.Fatal(err)
	}
	if err := ws.HandleInput(frame.ID, "<Esc>"); err != nil {
		t.Fatal(err)
	}

	if string(buf.Text.Bytes()) != "hi world\n" {
		t.Fatalf("after cw got %q", buf.Text.Bytes())
	}
}

func TestChangeLine(t *testing.T) {
	t.Parallel()

	ws := workspace.New()
	buf := ws.Editor.Scratch("*scratch*")
	buf.Text.SetBytes([]byte("remove me\nkeep\n"))

	frame, err := ws.NewFrame(80, 10, "")
	if err != nil {
		t.Fatal(err)
	}

	if err := ws.HandleInput(frame.ID, "c"); err != nil {
		t.Fatal(err)
	}
	if err := ws.HandleInput(frame.ID, "c"); err != nil {
		t.Fatal(err)
	}
	if ws.Editor.Mode != editor.InsertMode {
		t.Fatal("expected insert mode after cc")
	}

	if err := ws.HandleInput(frame.ID, "new"); err != nil {
		t.Fatal(err)
	}
	if err := ws.HandleInput(frame.ID, "<Esc>"); err != nil {
		t.Fatal(err)
	}

	if string(buf.Text.Bytes()) != "new\nkeep\n" {
		t.Fatalf("after cc got %q", buf.Text.Bytes())
	}
}

func TestJoinLines(t *testing.T) {
	t.Parallel()

	ws := workspace.New()
	buf := ws.Editor.Scratch("*scratch*")
	buf.Text.SetBytes([]byte("hello\nworld\n"))

	frame, err := ws.NewFrame(80, 10, "")
	if err != nil {
		t.Fatal(err)
	}

	if err := ws.HandleInput(frame.ID, "J"); err != nil {
		t.Fatal(err)
	}

	if string(buf.Text.Bytes()) != "hello world\n" {
		t.Fatalf("after J got %q", buf.Text.Bytes())
	}
}

func TestDeletePutsRegisterForPaste(t *testing.T) {
	t.Parallel()

	ws := workspace.New()
	buf := ws.Editor.Scratch("*scratch*")
	buf.Text.SetBytes([]byte("foo bar\n"))

	frame, err := ws.NewFrame(80, 10, "")
	if err != nil {
		t.Fatal(err)
	}

	if err := ws.HandleInput(frame.ID, "d"); err != nil {
		t.Fatal(err)
	}
	if err := ws.HandleInput(frame.ID, "w"); err != nil {
		t.Fatal(err)
	}
	if err := ws.HandleInput(frame.ID, "$"); err != nil {
		t.Fatal(err)
	}
	if err := ws.HandleInput(frame.ID, "p"); err != nil {
		t.Fatal(err)
	}

	if string(buf.Text.Bytes()) != "barfoo \n" {
		t.Fatalf("after dw+p got %q", buf.Text.Bytes())
	}
}
