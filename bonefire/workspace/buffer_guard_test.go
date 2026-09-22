package workspace

import (
	"errors"
	"testing"

	"vague/bonefire/editor"
)

func TestWithBufferLeaveBlocksModified(t *testing.T) {
	ws := New()
	frame, err := ws.NewFrame(80, 10, "")
	if err != nil {
		t.Fatal(err)
	}
	_, _, buf, _ := ws.FrameContext(frame.ID)
	_, _ = buf.Insert(buf.Text.Len(), []byte("x"))

	err = ws.withBufferLeave(frame.ID, false, func() error { return nil })
	if !errors.Is(err, editor.ErrNotSaved) {
		t.Fatalf("want ErrNotSaved, got %v", err)
	}
}

func TestHiddenAllowsLeaveModified(t *testing.T) {
	ws := New()
	_, _ = ws.Editor.ApplySetFileOption("hidden", true)
	frame, err := ws.NewFrame(80, 10, "")
	if err != nil {
		t.Fatal(err)
	}
	_, _, buf, _ := ws.FrameContext(frame.ID)
	_, _ = buf.Insert(buf.Text.Len(), []byte("x"))

	err = ws.withBufferLeave(frame.ID, false, func() error { return nil })
	if err != nil {
		t.Fatalf("hidden: %v", err)
	}
}

func TestWriteAll(t *testing.T) {
	ws := New()
	dir := t.TempDir()
	path := dir + "/f.txt"
	buf, err := ws.Editor.LoadBufferPathAt("", path)
	if err != nil {
		t.Fatal(err)
	}
	_, _ = buf.Insert(0, []byte("hi"))
	frame, _ := ws.NewFrame(80, 10, dir)
	if err := ws.WriteAll(frame.ID, false); err != nil {
		t.Fatal(err)
	}
	if buf.Modified() {
		t.Fatal("still modified after wa")
	}
}
