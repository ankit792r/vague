package editor_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"vague/bonefire/editor"
	"vague/bonefire/workspace"
)

func TestGoToLine(t *testing.T) {
	t.Parallel()

	ws := workspace.New()
	buf := ws.Editor.Scratch("*scratch*")
	buf.Text.SetBytes([]byte("a\nb\nc\n"))

	frame, err := ws.NewFrame(80, 10, "")
	if err != nil {
		t.Fatal(err)
	}

	if err := ws.GoToLine(frame.ID, 2); err != nil {
		t.Fatal(err)
	}

	win := ws.Windows[frame.ActiveWindowID]
	point := buf.Text.PointOf(editor.WindowCursorForTest(win))
	if point.Line != 1 {
		t.Fatalf("cursor line = %d, want 1", point.Line)
	}

	if err := ws.GoToLine(frame.ID, 99); err == nil {
		t.Fatal("expected error for out of range line")
	} else if !errors.Is(err, editor.ErrInvalidLine) {
		t.Fatalf("err = %v, want ErrInvalidLine", err)
	}
}

func TestReloadCurrentBuffer(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "reload.txt")

	if err := os.WriteFile(path, []byte("v1\n"), 0644); err != nil {
		t.Fatal(err)
	}

	ws := workspace.New()
	frame, err := ws.NewFrame(80, 10, dir)
	if err != nil {
		t.Fatal(err)
	}

	buf, err := ws.OpenFile(frame.ID, path, false)
	if err != nil {
		t.Fatal(err)
	}

	if _, err := buf.Insert(buf.Text.Len(), []byte("edit")); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(path, []byte("v2\n"), 0644); err != nil {
		t.Fatal(err)
	}

	_, err = ws.ReloadCurrentBuffer(frame.ID, false)
	if !errors.Is(err, editor.ErrNotSaved) {
		t.Fatalf("reload without bang: err = %v, want ErrNotSaved", err)
	}

	reloaded, err := ws.ReloadCurrentBuffer(frame.ID, true)
	if err != nil {
		t.Fatal(err)
	}

	if string(reloaded.Text.Bytes()) != "v2\n" {
		t.Fatalf("got %q, want v2\\n", reloaded.Text.Bytes())
	}
}
