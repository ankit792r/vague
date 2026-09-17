package buffer_test

import (
	"os"
	"path/filepath"
	"testing"

	"vague/bonefire/buffer"
)

func TestLoadSaveRoundTripLF(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "hello.txt")
	original := []byte("hello\nworld\n")

	if err := os.WriteFile(path, original, 0644); err != nil {
		t.Fatal(err)
	}

	buf, err := buffer.Load(path)
	if err != nil {
		t.Fatal(err)
	}

	if buf.Name != "hello.txt" {
		t.Fatalf("name = %q, want hello.txt", buf.Name)
	}

	if buf.Modified() {
		t.Fatal("expected clean buffer after load")
	}

	if string(buf.Text.Bytes()) != "hello\nworld\n" {
		t.Fatalf("got %q", buf.Text.Bytes())
	}

	if err := buf.Save(); err != nil {
		t.Fatal(err)
	}

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	if string(got) != string(original) {
		t.Fatalf("saved %q, want %q", got, original)
	}
}

func TestLoadSavePreservesNoEOL(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "noeol.txt")
	original := []byte("hello")

	if err := os.WriteFile(path, original, 0644); err != nil {
		t.Fatal(err)
	}

	buf, err := buffer.Load(path)
	if err != nil {
		t.Fatal(err)
	}

	if !buf.NoEOL {
		t.Fatal("expected NoEOL")
	}

	if err := buf.Save(); err != nil {
		t.Fatal(err)
	}

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	if string(got) != string(original) {
		t.Fatalf("saved %q, want %q", got, original)
	}
}

func TestLoadSaveCRLF(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "dos.txt")
	original := []byte("a\r\nb\r\n")

	if err := os.WriteFile(path, original, 0644); err != nil {
		t.Fatal(err)
	}

	buf, err := buffer.Load(path)
	if err != nil {
		t.Fatal(err)
	}

	if buf.LineEnding != buffer.CRLF {
		t.Fatalf("line ending = %v, want CRLF", buf.LineEnding)
	}

	if err := buf.Save(); err != nil {
		t.Fatal(err)
	}

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	if string(got) != string(original) {
		t.Fatalf("saved %q, want %q", got, original)
	}
}

func TestModifiedAfterEdit(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "edit.txt")

	if err := os.WriteFile(path, []byte("hi\n"), 0644); err != nil {
		t.Fatal(err)
	}

	buf, err := buffer.Load(path)
	if err != nil {
		t.Fatal(err)
	}

	buf.NoteEdit()
	if !buf.Modified() {
		t.Fatal("expected modified after edit")
	}

	if err := buf.Save(); err != nil {
		t.Fatal(err)
	}

	if buf.Modified() {
		t.Fatal("expected clean after save")
	}
}

func TestSaveDetectsExternalChange(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "changed.txt")

	if err := os.WriteFile(path, []byte("v1\n"), 0644); err != nil {
		t.Fatal(err)
	}

	buf, err := buffer.Load(path)
	if err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(path, []byte("v2\n"), 0644); err != nil {
		t.Fatal(err)
	}

	if err := buf.Save(); err != buffer.ErrFileChanged {
		t.Fatalf("got %v, want ErrFileChanged", err)
	}

	if err := buf.SaveForce(); err != nil {
		t.Fatal(err)
	}

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	if string(got) != "v1\n" {
		t.Fatalf("force saved %q, want v1\\n", got)
	}
}

func TestSaveEmptyFile(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "empty.txt")

	if err := os.WriteFile(path, nil, 0644); err != nil {
		t.Fatal(err)
	}

	buf, err := buffer.Load(path)
	if err != nil {
		t.Fatal(err)
	}

	if err := buf.Save(); err != nil {
		t.Fatal(err)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}

	if info.Size() != 0 {
		t.Fatalf("expected empty file, got %d bytes", info.Size())
	}
}
