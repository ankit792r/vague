package editor

import (
	"errors"
	"testing"
)

func TestQuitFrameRefusesModifiedBuffer(t *testing.T) {
	t.Parallel()

	ed := NewEditor()
	buf := ed.Scratch("*scratch*")
	buf.Text.SetBytes([]byte("hi"))

	frame, err := ed.NewFrame(80, 10, "")
	if err != nil {
		t.Fatal(err)
	}

	if _, err := buf.Insert(buf.Text.Len(), []byte("!")); err != nil {
		t.Fatal(err)
	}

	err = ed.QuitFrame(frame.ID, false)
	if !errors.Is(err, ErrNotSaved) {
		t.Fatalf("got %v, want ErrNotSaved", err)
	}

	if _, ok := ed.Frames[frame.ID]; !ok {
		t.Fatal("frame should still exist after refused quit")
	}
}

func TestQuitFrameForceCloses(t *testing.T) {
	t.Parallel()

	ed := NewEditor()
	buf := ed.Scratch("*scratch*")
	buf.Text.SetBytes([]byte("hi"))

	frame, err := ed.NewFrame(80, 10, "")
	if err != nil {
		t.Fatal(err)
	}

	if _, err := buf.Insert(buf.Text.Len(), []byte("!")); err != nil {
		t.Fatal(err)
	}

	if err := ed.QuitFrame(frame.ID, true); err != nil {
		t.Fatal(err)
	}

	if _, ok := ed.Frames[frame.ID]; ok {
		t.Fatal("frame should be gone after :q!")
	}
}

func TestQuitFrameCleanBuffer(t *testing.T) {
	t.Parallel()

	ed := NewEditor()
	ed.Scratch("*scratch*")

	frame, err := ed.NewFrame(80, 10, "")
	if err != nil {
		t.Fatal(err)
	}

	if err := ed.QuitFrame(frame.ID, false); err != nil {
		t.Fatal(err)
	}

	if _, ok := ed.Frames[frame.ID]; ok {
		t.Fatal("frame should be gone after :q")
	}
}
