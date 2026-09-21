package editor_test

import (
	"vague/bonefire/workspace"
	"os"
	"path/filepath"
	"testing"
)

func TestNewFrameOpensExistingFile(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "hello.py")
	content := []byte("print('hi')\n")

	if err := os.WriteFile(path, content, 0644); err != nil {
		t.Fatal(err)
	}

	ws := workspace.New()
	ed := ws.Editor
	frame, err := ws.NewFrame(80, 10, dir, path)
	if err != nil {
		t.Fatal(err)
	}

	win := ws.Windows[frame.ActiveWindowID]
	buf := ed.Buffers[win.BufferId]

	if buf.Name != "hello.py" {
		t.Fatalf("name = %q, want hello.py", buf.Name)
	}

	if string(buf.Text.Bytes()) != string(content) {
		t.Fatalf("got %q, want %q", buf.Text.Bytes(), content)
	}

	redraw, ok := ws.RenderRedraw(frame.ID)
	if !ok {
		t.Fatal("expected redraw")
	}

	if redraw.Buffer.Name != "hello.py" {
		t.Fatalf("redraw buffer = %q, want hello.py", redraw.Buffer.Name)
	}
}

func TestNewFrameOpensMissingFileAsEmptyBuffer(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "new.py")

	ws := workspace.New()
	ed := ws.Editor
	frame, err := ws.NewFrame(80, 10, dir, path)
	if err != nil {
		t.Fatal(err)
	}

	win := ws.Windows[frame.ActiveWindowID]
	buf := ed.Buffers[win.BufferId]

	if buf.Name != "new.py" {
		t.Fatalf("name = %q, want new.py", buf.Name)
	}

	if buf.Text.Len() != 0 {
		t.Fatalf("expected empty buffer, got %q", buf.Text.Bytes())
	}

	if buf.Modified() {
		t.Fatal("expected unmodified new file buffer")
	}

	redraw, ok := ws.RenderRedraw(frame.ID)
	if !ok {
		t.Fatal("expected redraw")
	}

	if redraw.Buffer.Name != "new.py" {
		t.Fatalf("redraw buffer = %q, want new.py", redraw.Buffer.Name)
	}
}

func TestOpenFileLoadsIntoFrame(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "sample.go")
	content := []byte("package main\n\nfunc main() {}\n")

	if err := os.WriteFile(path, content, 0644); err != nil {
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

	if buf.Name != "sample.go" {
		t.Fatalf("name = %q, want sample.go", buf.Name)
	}

	redraw, ok := ws.RenderRedraw(frame.ID)
	if !ok {
		t.Fatal("expected redraw")
	}

	if redraw.Buffer.Name != "sample.go" {
		t.Fatalf("redraw buffer = %q, want sample.go", redraw.Buffer.Name)
	}
}

func TestWriteFilePersistsEdits(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "out.txt")

	ws := workspace.New()
	frame, err := ws.NewFrame(80, 10, dir)
	if err != nil {
		t.Fatal(err)
	}

	if _, err := ws.OpenFile(frame.ID, path, false); err != nil {
		t.Fatal(err)
	}

	if err := ws.HandleInput(frame.ID, "i"); err != nil {
		t.Fatal(err)
	}

	if err := ws.HandleInput(frame.ID, "h"); err != nil {
		t.Fatal(err)
	}

	if err := ws.HandleInput(frame.ID, "i"); err != nil {
		t.Fatal(err)
	}

	if err := ws.WriteFile(frame.ID, "", false); err != nil {
		t.Fatal(err)
	}

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	if string(got) != "hi\n" {
		t.Fatalf("saved %q, want hi\\n", got)
	}
}

func TestOpenFileReusesExistingBuffer(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "reuse.txt")

	if err := os.WriteFile(path, []byte("one\n"), 0644); err != nil {
		t.Fatal(err)
	}

	ws := workspace.New()
	frame, err := ws.NewFrame(80, 10, dir)
	if err != nil {
		t.Fatal(err)
	}

	first, err := ws.OpenFile(frame.ID, path, false)
	if err != nil {
		t.Fatal(err)
	}

	second, err := ws.OpenFile(frame.ID, path, false)
	if err != nil {
		t.Fatal(err)
	}

	if first.ID != second.ID {
		t.Fatalf("got buffer ids %d and %d, want same", first.ID, second.ID)
	}
}

func TestOpenFileForceReloads(t *testing.T) {
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

	if err := os.WriteFile(path, []byte("v2\n"), 0644); err != nil {
		t.Fatal(err)
	}

	if _, err := ws.OpenFile(frame.ID, path, false); err != nil {
		t.Fatal(err)
	}

	if string(buf.Text.Bytes()) != "v1\n" {
		t.Fatalf("without bang got %q, want v1\\n", buf.Text.Bytes())
	}

	if _, err := ws.OpenFile(frame.ID, path, true); err != nil {
		t.Fatal(err)
	}

	if string(buf.Text.Bytes()) != "v2\n" {
		t.Fatalf("with bang got %q, want v2\\n", buf.Text.Bytes())
	}
}

func TestWriteRelativePathUsesFrameWorkDir(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()

	ws := workspace.New()
	frame, err := ws.NewFrame(80, 10, dir)
	if err != nil {
		t.Fatal(err)
	}

	if err := ws.HandleInput(frame.ID, "i"); err != nil {
		t.Fatal(err)
	}

	if err := ws.HandleInput(frame.ID, "x"); err != nil {
		t.Fatal(err)
	}

	if err := ws.WriteFile(frame.ID, "hello.py", false); err != nil {
		t.Fatal(err)
	}

	got := filepath.Join(dir, "hello.py")
	if _, err := os.Stat(got); err != nil {
		t.Fatalf("expected %q to exist: %v", got, err)
	}
}

func TestWriteScratchBufferNoFileName(t *testing.T) {
	t.Parallel()

	ws := workspace.New()
	ws.Editor.Scratch("*scratch*")

	frame, err := ws.NewFrame(80, 10, "")
	if err != nil {
		t.Fatal(err)
	}

	err = ws.WriteFile(frame.ID, "", false)
	if err == nil {
		t.Fatal("expected error writing scratch buffer with no path")
	}
	if err.Error() != "No file name" {
		t.Fatalf("got %q, want %q", err.Error(), "No file name")
	}
}
