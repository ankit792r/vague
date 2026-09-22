package editor_test

import (
	"testing"

	"vague/bonefire/editor"
	"vague/bonefire/workspace"
)

func TestCountMotionDown(t *testing.T) {
	t.Parallel()

	ws := workspace.New()
	buf := ws.Editor.Scratch("*scratch*")
	buf.Text.SetBytes([]byte("a\nb\nc"))

	frame, err := ws.NewFrame(80, 10, "")
	if err != nil {
		t.Fatal(err)
	}
	if err := ws.SetWindowWrap(frame.ID, false); err != nil {
		t.Fatal(err)
	}

	for _, key := range []string{"3", "j"} {
		if err := ws.HandleInput(frame.ID, key); err != nil {
			t.Fatal(err)
		}
	}

	win := ws.Windows[frame.ActiveWindowID]
	line := buf.Text.PointOf(editor.WindowCursorForTest(win)).Line
	if line != 2 {
		t.Fatalf("cursor line = %d, want 2", line)
	}
}

func TestCountDeleteLines(t *testing.T) {
	t.Parallel()

	ws := workspace.New()
	buf := ws.Editor.Scratch("*scratch*")
	buf.Text.SetBytes([]byte("one\ntwo\nthree"))

	frame, err := ws.NewFrame(80, 10, "")
	if err != nil {
		t.Fatal(err)
	}

	for _, key := range []string{"2", "d", "d"} {
		if err := ws.HandleInput(frame.ID, key); err != nil {
			t.Fatal(err)
		}
	}

	if string(buf.Text.Bytes()) != "three" {
		t.Fatalf("after 2dd got %q", buf.Text.Bytes())
	}
}

func TestCountDeleteChar(t *testing.T) {
	t.Parallel()

	ws := workspace.New()
	buf := ws.Editor.Scratch("*scratch*")
	buf.Text.SetBytes([]byte("abcdef\n"))

	frame, err := ws.NewFrame(80, 10, "")
	if err != nil {
		t.Fatal(err)
	}

	for _, key := range []string{"3", "x"} {
		if err := ws.HandleInput(frame.ID, key); err != nil {
			t.Fatal(err)
		}
	}

	if string(buf.Text.Bytes()) != "def\n" {
		t.Fatalf("after 3x got %q", buf.Text.Bytes())
	}
}

func TestBareZeroStillLineStart(t *testing.T) {
	t.Parallel()

	ws := workspace.New()
	buf := ws.Editor.Scratch("*scratch*")
	buf.Text.SetBytes([]byte("hello\n"))

	frame, err := ws.NewFrame(80, 10, "")
	if err != nil {
		t.Fatal(err)
	}

	if err := ws.HandleInput(frame.ID, "l"); err != nil {
		t.Fatal(err)
	}
	if err := ws.HandleInput(frame.ID, "0"); err != nil {
		t.Fatal(err)
	}

	win := ws.Windows[frame.ActiveWindowID]
	if editor.WindowCursorForTest(win) != 0 {
		t.Fatalf("cursor = %d, want 0", editor.WindowCursorForTest(win))
	}
}

func TestCountGotoLine(t *testing.T) {
	t.Parallel()

	ws := workspace.New()
	buf := ws.Editor.Scratch("*scratch*")
	buf.Text.SetBytes([]byte("a\nb\nc"))

	frame, err := ws.NewFrame(80, 10, "")
	if err != nil {
		t.Fatal(err)
	}

	for _, key := range []string{"2", "G"} {
		if err := ws.HandleInput(frame.ID, key); err != nil {
			t.Fatal(err)
		}
	}

	win := ws.Windows[frame.ActiveWindowID]
	line := buf.Text.PointOf(editor.WindowCursorForTest(win)).Line
	if line != 1 {
		t.Fatalf("cursor line = %d, want 1", line)
	}
}
