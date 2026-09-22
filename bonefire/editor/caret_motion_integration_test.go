package editor_test

import (
	"testing"

	"vague/bonefire/editor"
	"vague/bonefire/workspace"
)

func TestCaretMotionKeys(t *testing.T) {
	ws := workspace.New()
	buf := ws.Editor.Scratch("*scratch*")
	buf.Text.SetBytes([]byte("  alpha\n    beta\n"))

	frame, err := ws.NewFrame(80, 10, "")
	if err != nil {
		t.Fatal(err)
	}

	win := ws.Windows[frame.ActiveWindowID]
	editor.SetWindowCursorForTest(buf, win, 0)

	if err := ws.HandleInput(frame.ID, "^"); err != nil {
		t.Fatal(err)
	}
	p := buf.Text.PointOf(editor.WindowCursorForTest(win))
	if p.Line != 0 || p.Col != 2 {
		t.Fatalf("^ got line=%d col=%d, want 0,2", p.Line, p.Col)
	}

	editor.SetWindowCursorForTest(buf, win, 0)
	if err := ws.HandleInput(frame.ID, "_"); err != nil {
		t.Fatal(err)
	}
	p = buf.Text.PointOf(editor.WindowCursorForTest(win))
	if p.Line != 0 || p.Col != 2 {
		t.Fatalf("_ got line=%d col=%d, want 0,2", p.Line, p.Col)
	}

	editor.SetWindowCursorForTest(buf, win, 0)
	for _, key := range []string{"2", "_"} {
		if err := ws.HandleInput(frame.ID, key); err != nil {
			t.Fatal(err)
		}
	}
	p = buf.Text.PointOf(editor.WindowCursorForTest(win))
	if p.Line != 1 || p.Col != 4 {
		t.Fatalf("2_ got line=%d col=%d, want 1,4", p.Line, p.Col)
	}

	editor.SetWindowCursorForTest(buf, win, buf.Text.LineEnd(1))
	for _, key := range []string{"5", "|"} {
		if err := ws.HandleInput(frame.ID, key); err != nil {
			t.Fatal(err)
		}
	}
	p = buf.Text.PointOf(editor.WindowCursorForTest(win))
	if p.Col != 4 {
		t.Fatalf("5| got col=%d, want 4", p.Col)
	}
}
