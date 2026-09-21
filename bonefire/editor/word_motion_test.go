package editor_test

import (
	"testing"

	"vague/bonefire/editor"
	"vague/bonefire/workspace"
)

func TestWordMotionKeys(t *testing.T) {
	t.Parallel()

	ws := workspace.New()
	ed := ws.Editor
	buf := ed.Scratch("*scratch*")
	buf.Text.SetBytes([]byte("hello world"))

	frame, err := ws.NewFrame(80, 10, "")
	if err != nil {
		t.Fatal(err)
	}

	win := ws.Windows[frame.ActiveWindowID]

	if err := ws.HandleInput(frame.ID, "w"); err != nil {
		t.Fatal(err)
	}

	point := buf.Text.PointOf(editor.WindowCursorForTest(win))
	if point.Col != 6 {
		t.Fatalf("after w: col = %d, want 6", point.Col)
	}

	if err := ws.HandleInput(frame.ID, "e"); err != nil {
		t.Fatal(err)
	}

	point = buf.Text.PointOf(editor.WindowCursorForTest(win))
	if point.Col != 10 {
		t.Fatalf("after e: col = %d, want 10", point.Col)
	}

	if err := ws.HandleInput(frame.ID, "b"); err != nil {
		t.Fatal(err)
	}

	point = buf.Text.PointOf(editor.WindowCursorForTest(win))
	if point.Col != 6 {
		t.Fatalf("after b: col = %d, want 6", point.Col)
	}
}
