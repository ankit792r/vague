package editor

import "testing"

func TestNormalKeyMovesCursor(t *testing.T) {
	t.Parallel()

	ed := NewEditor()
	buf := ed.Scratch("*scratch*")
	buf.Text.SetBytes([]byte("ab\ncd"))

	frame, err := ed.NewFrame(80, 10)
	if err != nil {
		t.Fatal(err)
	}

	win := ed.Windows[frame.ActiveWindowID]

	if err := ed.HandleInput(frame.ID, "l"); err != nil {
		t.Fatal(err)
	}

	point := buf.Text.PointOf(windowCursor(win))
	if point.Line != 0 || point.Col != 1 {
		t.Fatalf("after l: got line=%d col=%d, want 0,1", point.Line, point.Col)
	}

	if err := ed.HandleInput(frame.ID, "j"); err != nil {
		t.Fatal(err)
	}

	point = buf.Text.PointOf(windowCursor(win))
	if point.Line != 1 || point.Col != 1 {
		t.Fatalf("after j: got line=%d col=%d, want 1,1", point.Line, point.Col)
	}
}

func TestAppendEndOfLineEntersInsertAtEnd(t *testing.T) {
	t.Parallel()

	ed := NewEditor()
	buf := ed.Scratch("*scratch*")
	buf.Text.SetBytes([]byte("hi"))

	frame, _ := ed.NewFrame(80, 10)
	win := ed.Windows[frame.ActiveWindowID]

	if err := ed.HandleInput(frame.ID, "A"); err != nil {
		t.Fatal(err)
	}

	if ed.Mode != InsertMode {
		t.Fatal("expected insert mode after A")
	}

	if got := windowCursor(win); got != 2 {
		t.Fatalf("cursor after A = %d, want 2", got)
	}

	if err := ed.HandleInput(frame.ID, "!"); err != nil {
		t.Fatal(err)
	}

	if string(buf.Text.Bytes()) != "hi!" {
		t.Fatalf("got %q, want %q", buf.Text.Bytes(), "hi!")
	}
}

func TestAppendMovesCursorPastLastChar(t *testing.T) {
	t.Parallel()

	ed := NewEditor()
	buf := ed.Scratch("*scratch*")
	buf.Text.SetBytes([]byte("hi"))

	frame, _ := ed.NewFrame(80, 10)
	win := ed.Windows[frame.ActiveWindowID]

	_ = ed.HandleInput(frame.ID, "l")
	_ = ed.HandleInput(frame.ID, "a")

	if got := windowCursor(win); got != 2 {
		t.Fatalf("cursor after a = %d, want 2", got)
	}
}

func TestInsertModeEditsScratchBuffer(t *testing.T) {
	t.Parallel()

	ed := NewEditor()
	buf := ed.Scratch("*scratch*")
	buf.Text.SetBytes([]byte("hi"))

	frame, err := ed.NewFrame(80, 10)
	if err != nil {
		t.Fatal(err)
	}

	if err := ed.HandleInput(frame.ID, "l"); err != nil {
		t.Fatal(err)
	}

	if err := ed.HandleInput(frame.ID, "a"); err != nil {
		t.Fatal(err)
	}

	if ed.Mode != InsertMode {
		t.Fatal("expected insert mode")
	}

	if err := ed.HandleInput(frame.ID, "!"); err != nil {
		t.Fatal(err)
	}

	if err := ed.HandleInput(frame.ID, "<Space>"); err != nil {
		t.Fatal(err)
	}

	if err := ed.HandleInput(frame.ID, "x"); err != nil {
		t.Fatal(err)
	}

	if string(buf.Text.Bytes()) != "hi! x" {
		t.Fatalf("got %q, want %q", buf.Text.Bytes(), "hi! x")
	}

	if err := ed.HandleInput(frame.ID, "<Esc>"); err != nil {
		t.Fatal(err)
	}

	if ed.Mode != NormalMode {
		t.Fatal("expected normal mode after Esc")
	}
}
