package editor

import (
	"os"
	"path/filepath"
	"testing"
)

func TestOpenFileLoadsIntoFrame(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "sample.go")
	content := []byte("package main\n\nfunc main() {}\n")

	if err := os.WriteFile(path, content, 0644); err != nil {
		t.Fatal(err)
	}

	ed := NewEditor()
	frame, err := ed.NewFrame(80, 10)
	if err != nil {
		t.Fatal(err)
	}

	buf, err := ed.OpenFile(frame.ID, path, false)
	if err != nil {
		t.Fatal(err)
	}

	if buf.Name != "sample.go" {
		t.Fatalf("name = %q, want sample.go", buf.Name)
	}

	redraw, ok := ed.RenderRedraw(frame.ID)
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

	ed := NewEditor()
	frame, err := ed.NewFrame(80, 10)
	if err != nil {
		t.Fatal(err)
	}

	if _, err := ed.OpenFile(frame.ID, path, false); err != nil {
		t.Fatal(err)
	}

	if err := ed.HandleInput(frame.ID, "i"); err != nil {
		t.Fatal(err)
	}

	if err := ed.HandleInput(frame.ID, "h"); err != nil {
		t.Fatal(err)
	}

	if err := ed.HandleInput(frame.ID, "i"); err != nil {
		t.Fatal(err)
	}

	if err := ed.WriteFile(frame.ID, "", false); err != nil {
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

	ed := NewEditor()
	frame, err := ed.NewFrame(80, 10)
	if err != nil {
		t.Fatal(err)
	}

	first, err := ed.OpenFile(frame.ID, path, false)
	if err != nil {
		t.Fatal(err)
	}

	second, err := ed.OpenFile(frame.ID, path, false)
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

	ed := NewEditor()
	frame, err := ed.NewFrame(80, 10)
	if err != nil {
		t.Fatal(err)
	}

	buf, err := ed.OpenFile(frame.ID, path, false)
	if err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(path, []byte("v2\n"), 0644); err != nil {
		t.Fatal(err)
	}

	if _, err := ed.OpenFile(frame.ID, path, false); err != nil {
		t.Fatal(err)
	}

	if string(buf.Text.Bytes()) != "v1\n" {
		t.Fatalf("without bang got %q, want v1\\n", buf.Text.Bytes())
	}

	if _, err := ed.OpenFile(frame.ID, path, true); err != nil {
		t.Fatal(err)
	}

	if string(buf.Text.Bytes()) != "v2\n" {
		t.Fatalf("with bang got %q, want v2\\n", buf.Text.Bytes())
	}
}
