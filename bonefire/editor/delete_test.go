package editor

import "testing"

func TestDeleteChar(t *testing.T) {
	t.Parallel()

	ed := NewEditor()
	buf := ed.Scratch("*scratch*")
	buf.Text.SetBytes([]byte("abcd"))

	frame, err := ed.NewFrame(80, 10)
	if err != nil {
		t.Fatal(err)
	}

	if err := ed.HandleInput(frame.ID, "l"); err != nil {
		t.Fatal(err)
	}

	if err := ed.HandleInput(frame.ID, "x"); err != nil {
		t.Fatal(err)
	}

	if string(buf.Text.Bytes()) != "acd" {
		t.Fatalf("after x got %q, want acd", buf.Text.Bytes())
	}
}

func TestDeleteCharOnEmptyLineNoOp(t *testing.T) {
	t.Parallel()

	ed := NewEditor()
	buf := ed.Scratch("*scratch*")
	buf.Text.SetBytes([]byte("hello\n\nworld"))

	frame, err := ed.NewFrame(80, 10)
	if err != nil {
		t.Fatal(err)
	}

	if err := ed.HandleInput(frame.ID, "j"); err != nil {
		t.Fatal(err)
	}

	if err := ed.HandleInput(frame.ID, "x"); err != nil {
		t.Fatal(err)
	}

	if string(buf.Text.Bytes()) != "hello\n\nworld" {
		t.Fatalf("x on empty line should not change buffer, got %q", buf.Text.Bytes())
	}
}

func TestDeleteLine(t *testing.T) {
	t.Parallel()

	ed := NewEditor()
	buf := ed.Scratch("*scratch*")
	buf.Text.SetBytes([]byte("keep\nremove\nstay"))

	frame, err := ed.NewFrame(80, 10)
	if err != nil {
		t.Fatal(err)
	}

	if err := ed.HandleInput(frame.ID, "j"); err != nil {
		t.Fatal(err)
	}

	if err := ed.HandleInput(frame.ID, "d"); err != nil {
		t.Fatal(err)
	}

	if err := ed.HandleInput(frame.ID, "d"); err != nil {
		t.Fatal(err)
	}

	if string(buf.Text.Bytes()) != "keep\nstay" {
		t.Fatalf("after dd got %q, want keep\\nstay", buf.Text.Bytes())
	}
}

func TestDeleteCharUndo(t *testing.T) {
	t.Parallel()

	ed := NewEditor()
	buf := ed.Scratch("*scratch*")
	buf.Text.SetBytes([]byte("ab"))

	frame, err := ed.NewFrame(80, 10)
	if err != nil {
		t.Fatal(err)
	}

	if err := ed.HandleInput(frame.ID, "x"); err != nil {
		t.Fatal(err)
	}

	if string(buf.Text.Bytes()) != "b" {
		t.Fatalf("after x got %q", buf.Text.Bytes())
	}

	if err := ed.HandleInput(frame.ID, "u"); err != nil {
		t.Fatal(err)
	}

	if string(buf.Text.Bytes()) != "ab" {
		t.Fatalf("after undo got %q, want ab", buf.Text.Bytes())
	}
}

func TestDeleteLineUndo(t *testing.T) {
	t.Parallel()

	ed := NewEditor()
	buf := ed.Scratch("*scratch*")
	buf.Text.SetBytes([]byte("one\ntwo"))

	frame, err := ed.NewFrame(80, 10)
	if err != nil {
		t.Fatal(err)
	}

	if err := ed.HandleInput(frame.ID, "d"); err != nil {
		t.Fatal(err)
	}

	if err := ed.HandleInput(frame.ID, "d"); err != nil {
		t.Fatal(err)
	}

	if string(buf.Text.Bytes()) != "two" {
		t.Fatalf("after dd got %q", buf.Text.Bytes())
	}

	if err := ed.HandleInput(frame.ID, "u"); err != nil {
		t.Fatal(err)
	}

	if string(buf.Text.Bytes()) != "one\ntwo" {
		t.Fatalf("after undo got %q, want one\\ntwo", buf.Text.Bytes())
	}
}

func TestPendingDThenMotionClears(t *testing.T) {
	t.Parallel()

	ed := NewEditor()
	buf := ed.Scratch("*scratch*")
	buf.Text.SetBytes([]byte("ab"))

	frame, err := ed.NewFrame(80, 10)
	if err != nil {
		t.Fatal(err)
	}

	win := ed.Windows[frame.ActiveWindowID]

	if err := ed.HandleInput(frame.ID, "d"); err != nil {
		t.Fatal(err)
	}

	if err := ed.HandleInput(frame.ID, "l"); err != nil {
		t.Fatal(err)
	}

	if string(buf.Text.Bytes()) != "ab" {
		t.Fatalf("d then l should not delete, got %q", buf.Text.Bytes())
	}

	point := buf.Text.PointOf(windowCursor(win))
	if point.Col != 1 {
		t.Fatalf("cursor col = %d, want 1", point.Col)
	}
}
