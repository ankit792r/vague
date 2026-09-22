package editor_test

import (
	"testing"

	"vague/bonefire/editor"
	"vague/bonefire/workspace"
)

func TestCharFindKeys(t *testing.T) {
	ws := workspace.New()
	buf := ws.Editor.Scratch("*scratch*")
	buf.Text.SetBytes([]byte("foo bar foo\n"))

	frame, err := ws.NewFrame(80, 10, "")
	if err != nil {
		t.Fatal(err)
	}

	win := ws.Windows[frame.ActiveWindowID]
	editor.SetWindowCursorForTest(buf, win, 0)

	if err := ws.HandleInput(frame.ID, "f"); err != nil {
		t.Fatal(err)
	}
	if err := ws.HandleInput(frame.ID, "o"); err != nil {
		t.Fatal(err)
	}

	if got := editor.WindowCursorForTest(win); got != 1 {
		t.Fatalf("fo cursor = %d, want 1", got)
	}

	if err := ws.HandleInput(frame.ID, ";"); err != nil {
		t.Fatal(err)
	}
	if got := editor.WindowCursorForTest(win); got != 2 {
		t.Fatalf("; cursor = %d, want 2", got)
	}

	editor.SetWindowCursorForTest(buf, win, 10)
	if err := ws.HandleInput(frame.ID, ","); err != nil {
		t.Fatal(err)
	}
	if got := editor.WindowCursorForTest(win); got != 9 {
		t.Fatalf(", cursor = %d, want 9", got)
	}
}

func TestCharFindWithCount(t *testing.T) {
	ws := workspace.New()
	buf := ws.Editor.Scratch("*scratch*")
	buf.Text.SetBytes([]byte("ooo"))

	frame, err := ws.NewFrame(80, 10, "")
	if err != nil {
		t.Fatal(err)
	}

	win := ws.Windows[frame.ActiveWindowID]
	editor.SetWindowCursorForTest(buf, win, 0)

	for _, key := range []string{"3", "f", "o"} {
		if err := ws.HandleInput(frame.ID, key); err != nil {
			t.Fatal(err)
		}
	}

	if got := editor.WindowCursorForTest(win); got != 2 {
		t.Fatalf("3fo cursor = %d, want 2", got)
	}
}
