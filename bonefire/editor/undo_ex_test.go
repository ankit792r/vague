package editor_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"vague/bonefire/buffer"
	"vague/bonefire/workspace"
)

func TestUndolistExCommand(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, "u.txt")
	if err := os.WriteFile(path, []byte("a"), 0o644); err != nil {
		t.Fatal(err)
	}
	ws := workspace.New()
	frame, err := ws.NewFrame(80, 10, dir)
	if err != nil {
		t.Fatal(err)
	}
	loaded, err := ws.Editor.LoadBufferPath(path)
	if err != nil {
		t.Fatal(err)
	}
	ws.Windows[frame.ActiveWindowID].BufferId = loaded.ID
	if err := ws.HandleInput(frame.ID, "A"); err != nil {
		t.Fatal(err)
	}
	if err := ws.HandleInput(frame.ID, "b"); err != nil {
		t.Fatal(err)
	}
	if err := ws.HandleInput(frame.ID, "<Esc>"); err != nil {
		t.Fatal(err)
	}
	out, err := ws.RunExLine(frame.ID, "undolist")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, ">") {
		t.Fatalf("undolist missing current marker: %q", out)
	}
}

func TestRedoExCommand(t *testing.T) {
	t.Parallel()
	ws := workspace.New()
	buf := ws.Editor.Scratch("*scratch*")
	buf.Text.SetBytes([]byte("ab"))
	frame, _ := ws.NewFrame(80, 10, "")
	_ = ws.HandleInput(frame.ID, "A")
	_ = ws.HandleInput(frame.ID, "X")
	_ = ws.HandleInput(frame.ID, "<Esc>")
	_ = ws.HandleInput(frame.ID, "u")
	msg, err := ws.RunExLine(frame.ID, "redo")
	if err != nil {
		t.Fatal(err)
	}
	if msg != "redo" {
		t.Fatalf("msg = %q", msg)
	}
	if got := string(buf.Text.Bytes()); got != "abX" {
		t.Fatalf("got %q", got)
	}
}

func TestChangesExCommand(t *testing.T) {
	t.Parallel()
	ws := workspace.New()
	buf := ws.Editor.Scratch("*scratch*")
	buf.Text.SetBytes([]byte("hi"))
	frame, _ := ws.NewFrame(80, 10, "")
	_ = ws.HandleInput(frame.ID, "A")
	_ = ws.HandleInput(frame.ID, "!")
	_ = ws.HandleInput(frame.ID, "<Esc>")
	out, err := ws.RunExLine(frame.ID, "changes")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "change mark") {
		t.Fatalf("changes output: %q", out)
	}
}

func TestPersistentUndoAcrossReload(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, "persist.txt")
	if err := os.WriteFile(path, []byte("start"), 0o644); err != nil {
		t.Fatal(err)
	}

	ws := workspace.New()
	frame, _ := ws.NewFrame(80, 10, dir)
	loaded, err := ws.Editor.LoadBufferPath(path)
	if err != nil {
		t.Fatal(err)
	}
	ws.Windows[frame.ActiveWindowID].BufferId = loaded.ID
	_, _, buf, err := ws.FrameContext(frame.ID)
	if err != nil {
		t.Fatal(err)
	}
	_ = ws.HandleInput(frame.ID, "A")
	_ = ws.HandleInput(frame.ID, "ed")
	_ = ws.HandleInput(frame.ID, "<Esc>")
	if err := buf.Save(); err != nil {
		t.Fatal(err)
	}
	if buf.History.Seq() < 1 {
		t.Fatal("expected undo seq")
	}

	reloaded, err := buffer.Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if reloaded.History.Seq() != buf.History.Seq() {
		t.Fatalf("seq = %d want %d", reloaded.History.Seq(), buf.History.Seq())
	}
}

func TestEarlierUndoes(t *testing.T) {
	t.Parallel()
	ws := workspace.New()
	buf := ws.Editor.Scratch("*scratch*")
	buf.Text.SetBytes([]byte("ab"))
	frame, _ := ws.NewFrame(80, 10, "")
	_ = ws.HandleInput(frame.ID, "A")
	_ = ws.HandleInput(frame.ID, "X")
	_ = ws.HandleInput(frame.ID, "<Esc>")
	if _, err := ws.RunExLine(frame.ID, "earlier"); err != nil {
		t.Fatal(err)
	}
	if got := string(buf.Text.Bytes()); got != "ab" {
		t.Fatalf("got %q", got)
	}
}
