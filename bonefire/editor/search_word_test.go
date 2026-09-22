package editor_test

import (
	"testing"

	"vague/bonefire/editor"
	"vague/bonefire/workspace"
)

func TestStarHashSearch(t *testing.T) {
	ws := workspace.New()
	buf := ws.Editor.Scratch("*scratch*")
	buf.Text.SetBytes([]byte("foo bar foo\n"))

	frame, err := ws.NewFrame(80, 10, "")
	if err != nil {
		t.Fatal(err)
	}

	win := ws.Windows[frame.ActiveWindowID]
	editor.SetWindowCursorForTest(buf, win, 0)

	if err := ws.HandleInput(frame.ID, "*"); err != nil {
		t.Fatal(err)
	}
	if got := editor.WindowCursorForTest(win); got != 8 {
		t.Fatalf("* cursor = %d, want 8", got)
	}

	editor.SetWindowCursorForTest(buf, win, 8)
	if err := ws.HandleInput(frame.ID, "#"); err != nil {
		t.Fatal(err)
	}
	if got := editor.WindowCursorForTest(win); got != 0 {
		t.Fatalf("# cursor = %d, want 0", got)
	}
}
