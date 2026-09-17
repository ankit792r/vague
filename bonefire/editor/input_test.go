package editor

import "testing"

func TestNormalKeyMovesCursor(t *testing.T) {
	t.Parallel()

	ed := NewEditor()
	buf := ed.Scratch("*scratch*")
	buf.Text = "ab\ncd"

	frame, err := ed.NewFrame(80, 10)
	if err != nil {
		t.Fatal(err)
	}

	if err := ed.HandleInput(frame.ID, "l"); err != nil {
		t.Fatal(err)
	}

	win := ed.Windows[frame.ActiveWindowID]
	if win.CursorLine != 0 || win.CursorCol != 1 {
		t.Fatalf("after l: got line=%d col=%d, want 0,1", win.CursorLine, win.CursorCol)
	}

	if err := ed.HandleInput(frame.ID, "j"); err != nil {
		t.Fatal(err)
	}

	if win.CursorLine != 1 || win.CursorCol != 1 {
		t.Fatalf("after j: got line=%d col=%d, want 1,1", win.CursorLine, win.CursorCol)
	}
}
