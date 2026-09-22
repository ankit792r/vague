package editor_test

import (
	"strings"
	"testing"

	"vague/bonefire/workspace"
)

func TestBufferNextPrev(t *testing.T) {
	t.Parallel()

	ws := workspace.New()
	a := ws.Editor.Scratch("*a*")
	b := ws.Editor.Scratch("*b*")
	_ = b

	frame, err := ws.NewFrame(80, 10, "")
	if err != nil {
		t.Fatal(err)
	}

	next, err := ws.SwitchToNextBuffer(frame.ID, false)
	if err != nil {
		t.Fatal(err)
	}
	if next.ID != a.ID {
		t.Fatalf("next = %q, want *a*", next.Name)
	}

	prev, err := ws.SwitchToPrevBuffer(frame.ID, false)
	if err != nil {
		t.Fatal(err)
	}
	if prev.ID != b.ID {
		t.Fatalf("prev = %q, want *b*", prev.Name)
	}
}

func TestBufferByNumber(t *testing.T) {
	t.Parallel()

	ws := workspace.New()
	ws.Editor.Scratch("*a*")
	b := ws.Editor.Scratch("*b*")

	frame, err := ws.NewFrame(80, 10, "")
	if err != nil {
		t.Fatal(err)
	}

	got, err := ws.SwitchToBuffer(frame.ID, "2", false)
	if err != nil {
		t.Fatal(err)
	}
	if got.ID != b.ID {
		t.Fatalf("buffer 2 = %q", got.Name)
	}
}

func TestBufferListMessage(t *testing.T) {
	t.Parallel()

	ws := workspace.New()
	ws.Editor.Scratch("*a*")
	ws.Editor.Scratch("*b*")

	frame, err := ws.NewFrame(80, 10, "")
	if err != nil {
		t.Fatal(err)
	}

	msg, err := ws.BufferListMessage(frame.ID)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(msg, "*a*") || !strings.Contains(msg, "*b*") {
		t.Fatalf("list = %q", msg)
	}
	if !strings.Contains(msg, "%") {
		t.Fatalf("expected current buffer marker in %q", msg)
	}
}
