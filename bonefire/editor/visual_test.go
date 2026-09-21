package editor_test

import (
	"testing"

	"vague/bonefire/editor"
	"vague/bonefire/workspace"
)

func TestVisualCharDelete(t *testing.T) {
	t.Parallel()

	ws := workspace.New()
	buf := ws.Editor.Scratch("*scratch*")
	buf.Text.SetBytes([]byte("hello world\n"))

	frame, err := ws.NewFrame(80, 10, "")
	if err != nil {
		t.Fatal(err)
	}

	if err := ws.HandleInput(frame.ID, "v"); err != nil {
		t.Fatal(err)
	}
	if ws.Editor.Mode != editor.VisualMode {
		t.Fatalf("mode = %v, want visual", ws.Editor.Mode)
	}

	if err := ws.HandleInput(frame.ID, "l"); err != nil {
		t.Fatal(err)
	}
	if err := ws.HandleInput(frame.ID, "l"); err != nil {
		t.Fatal(err)
	}
	if err := ws.HandleInput(frame.ID, "l"); err != nil {
		t.Fatal(err)
	}
	if err := ws.HandleInput(frame.ID, "l"); err != nil {
		t.Fatal(err)
	}
	if err := ws.HandleInput(frame.ID, "d"); err != nil {
		t.Fatal(err)
	}

	if string(buf.Text.Bytes()) != " world\n" {
		t.Fatalf("after visual d got %q", buf.Text.Bytes())
	}
	if ws.Editor.Mode != editor.NormalMode {
		t.Fatalf("mode = %v, want normal", ws.Editor.Mode)
	}
}

func TestVisualLineYank(t *testing.T) {
	t.Parallel()

	ws := workspace.New()
	buf := ws.Editor.Scratch("*scratch*")
	buf.Text.SetBytes([]byte("one\ntwo\nthree\n"))

	frame, err := ws.NewFrame(80, 10, "")
	if err != nil {
		t.Fatal(err)
	}

	if err := ws.HandleInput(frame.ID, "j"); err != nil {
		t.Fatal(err)
	}
	if err := ws.HandleInput(frame.ID, "V"); err != nil {
		t.Fatal(err)
	}
	if err := ws.HandleInput(frame.ID, "y"); err != nil {
		t.Fatal(err)
	}

	if string(buf.Text.Bytes()) != "one\ntwo\nthree\n" {
		t.Fatalf("after visual y got %q", buf.Text.Bytes())
	}

	if err := ws.HandleInput(frame.ID, "G"); err != nil {
		t.Fatal(err)
	}
	if err := ws.HandleInput(frame.ID, "p"); err != nil {
		t.Fatal(err)
	}

	if string(buf.Text.Bytes()) != "one\ntwo\nthree\ntwo\n" {
		t.Fatalf("after paste linewise yank got %q", buf.Text.Bytes())
	}
}

func TestVisualChange(t *testing.T) {
	t.Parallel()

	ws := workspace.New()
	buf := ws.Editor.Scratch("*scratch*")
	buf.Text.SetBytes([]byte("abcd\n"))

	frame, err := ws.NewFrame(80, 10, "")
	if err != nil {
		t.Fatal(err)
	}

	if err := ws.HandleInput(frame.ID, "l"); err != nil {
		t.Fatal(err)
	}
	if err := ws.HandleInput(frame.ID, "v"); err != nil {
		t.Fatal(err)
	}
	if err := ws.HandleInput(frame.ID, "c"); err != nil {
		t.Fatal(err)
	}
	if ws.Editor.Mode != editor.InsertMode {
		t.Fatalf("mode = %v, want insert", ws.Editor.Mode)
	}

	if err := ws.HandleInput(frame.ID, "X"); err != nil {
		t.Fatal(err)
	}
	if err := ws.HandleInput(frame.ID, "<Esc>"); err != nil {
		t.Fatal(err)
	}

	if string(buf.Text.Bytes()) != "aXcd\n" {
		t.Fatalf("after visual c got %q", buf.Text.Bytes())
	}
}

func TestVisualRedrawIncludesSelection(t *testing.T) {
	t.Parallel()

	ws := workspace.New()
	buf := ws.Editor.Scratch("*scratch*")
	buf.Text.SetBytes([]byte("abcd\n"))

	frame, err := ws.NewFrame(80, 10, "")
	if err != nil {
		t.Fatal(err)
	}

	if err := ws.HandleInput(frame.ID, "l"); err != nil {
		t.Fatal(err)
	}
	if err := ws.HandleInput(frame.ID, "v"); err != nil {
		t.Fatal(err)
	}

	redraw, ok := ws.RenderRedraw(frame.ID)
	if !ok {
		t.Fatal("expected redraw")
	}
	if redraw.Mode != "visual" {
		t.Fatalf("mode = %q", redraw.Mode)
	}
	if redraw.Selection == nil || !redraw.Selection.Visible {
		t.Fatal("expected visible selection")
	}
	if redraw.Selection.Start.Column != 1 || redraw.Selection.End.Column != 1 {
		t.Fatalf("selection cols = %d..%d", redraw.Selection.Start.Column, redraw.Selection.End.Column)
	}
	_ = buf
}
