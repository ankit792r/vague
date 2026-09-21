package editor_test

import (
	"vague/bonefire/editor"
	"vague/bonefire/workspace"
	"errors"
	"testing"
)

func TestQuitFrameRefusesModifiedBuffer(t *testing.T) {
	t.Parallel()

	ws := workspace.New()
	ed := ws.Editor
	buf := ed.Scratch("*scratch*")
	buf.Text.SetBytes([]byte("hi"))

	frame, err := ws.NewFrame(80, 10, "")
	if err != nil {
		t.Fatal(err)
	}

	if _, err := buf.Insert(buf.Text.Len(), []byte("!")); err != nil {
		t.Fatal(err)
	}

	err = ws.QuitFrame(frame.ID, false)
	if !errors.Is(err, editor.ErrNotSaved) {
		t.Fatalf("got %v, want editor.ErrNotSaved", err)
	}

	if _, ok := ws.Frames[frame.ID]; !ok {
		t.Fatal("frame should still exist after refused quit")
	}
}

func TestQuitFrameForceCloses(t *testing.T) {
	t.Parallel()

	ws := workspace.New()
	ed := ws.Editor
	buf := ed.Scratch("*scratch*")
	buf.Text.SetBytes([]byte("hi"))

	frame, err := ws.NewFrame(80, 10, "")
	if err != nil {
		t.Fatal(err)
	}

	if _, err := buf.Insert(buf.Text.Len(), []byte("!")); err != nil {
		t.Fatal(err)
	}

	if err := ws.QuitFrame(frame.ID, true); err != nil {
		t.Fatal(err)
	}

	if _, ok := ws.Frames[frame.ID]; ok {
		t.Fatal("frame should be gone after :q!")
	}
}

func TestQuitFrameCleanBuffer(t *testing.T) {
	t.Parallel()

	ws := workspace.New()
	ed := ws.Editor
	ed.Scratch("*scratch*")

	frame, err := ws.NewFrame(80, 10, "")
	if err != nil {
		t.Fatal(err)
	}

	if err := ws.QuitFrame(frame.ID, false); err != nil {
		t.Fatal(err)
	}

	if _, ok := ws.Frames[frame.ID]; ok {
		t.Fatal("frame should be gone after :q")
	}
}
