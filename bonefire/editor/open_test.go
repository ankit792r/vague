package editor

import "testing"

func TestOpenLineBelow(t *testing.T) {
	t.Parallel()

	ed := NewEditor()
	buf := ed.Scratch("*scratch*")
	buf.Text.SetBytes([]byte("hello"))

	frame, err := ed.NewFrame(80, 10, "")
	if err != nil {
		t.Fatal(err)
	}

	win := ed.Windows[frame.ActiveWindowID]

	if err := ed.HandleInput(frame.ID, "o"); err != nil {
		t.Fatal(err)
	}

	if ed.Mode != InsertMode {
		t.Fatal("expected insert mode after o")
	}

	if string(buf.Text.Bytes()) != "hello\n" {
		t.Fatalf("got %q, want hello\\n", buf.Text.Bytes())
	}

	point := buf.Text.PointOf(windowCursor(win))
	if point.Line != 1 || point.Col != 0 {
		t.Fatalf("cursor at line=%d col=%d, want line 1 col 0", point.Line, point.Col)
	}

	if err := ed.HandleInput(frame.ID, "x"); err != nil {
		t.Fatal(err)
	}

	if err := ed.HandleInput(frame.ID, "<Esc>"); err != nil {
		t.Fatal(err)
	}

	if string(buf.Text.Bytes()) != "hello\nx" {
		t.Fatalf("got %q, want hello\\nx", buf.Text.Bytes())
	}
}

func TestOpenLineAbove(t *testing.T) {
	t.Parallel()

	ed := NewEditor()
	buf := ed.Scratch("*scratch*")
	buf.Text.SetBytes([]byte("hello"))

	frame, err := ed.NewFrame(80, 10, "")
	if err != nil {
		t.Fatal(err)
	}

	win := ed.Windows[frame.ActiveWindowID]

	if err := ed.HandleInput(frame.ID, "O"); err != nil {
		t.Fatal(err)
	}

	if ed.Mode != InsertMode {
		t.Fatal("expected insert mode after O")
	}

	if string(buf.Text.Bytes()) != "\nhello" {
		t.Fatalf("got %q, want \\nhello", buf.Text.Bytes())
	}

	point := buf.Text.PointOf(windowCursor(win))
	if point.Line != 0 || point.Col != 0 {
		t.Fatalf("cursor at line=%d col=%d, want line 0 col 0", point.Line, point.Col)
	}
}

func TestOpenLineUndo(t *testing.T) {
	t.Parallel()

	ed := NewEditor()
	buf := ed.Scratch("*scratch*")
	buf.Text.SetBytes([]byte("ab"))

	frame, err := ed.NewFrame(80, 10, "")
	if err != nil {
		t.Fatal(err)
	}

	if err := ed.HandleInput(frame.ID, "o"); err != nil {
		t.Fatal(err)
	}

	if err := ed.HandleInput(frame.ID, "z"); err != nil {
		t.Fatal(err)
	}

	if err := ed.HandleInput(frame.ID, "<Esc>"); err != nil {
		t.Fatal(err)
	}

	if string(buf.Text.Bytes()) != "ab\nz" {
		t.Fatalf("before undo got %q", buf.Text.Bytes())
	}

	if err := ed.HandleInput(frame.ID, "u"); err != nil {
		t.Fatal(err)
	}

	if string(buf.Text.Bytes()) != "ab" {
		t.Fatalf("after undo got %q, want ab", buf.Text.Bytes())
	}
}
