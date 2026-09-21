package editor_test

import (
	"vague/bonefire/editor"
	"testing"

	"vague/bonefire/workspace"
)

func TestNormalKeyMovesCursor(t *testing.T) {
	t.Parallel()

	ws := workspace.New()
	ed := ws.Editor
	buf := ed.Scratch("*scratch*")
	buf.Text.SetBytes([]byte("ab\ncd"))

	frame, err := ws.NewFrame(80, 10, "")
	if err != nil {
		t.Fatal(err)
	}

	win := ws.Windows[frame.ActiveWindowID]

	if err := ws.HandleInput(frame.ID, "l"); err != nil {
		t.Fatal(err)
	}

	point := buf.Text.PointOf(editor.WindowCursorForTest(win))
	if point.Line != 0 || point.Col != 1 {
		t.Fatalf("after l: got line=%d col=%d, want 0,1", point.Line, point.Col)
	}

	if err := ws.HandleInput(frame.ID, "j"); err != nil {
		t.Fatal(err)
	}

	point = buf.Text.PointOf(editor.WindowCursorForTest(win))
	if point.Line != 1 || point.Col != 1 {
		t.Fatalf("after j: got line=%d col=%d, want 1,1", point.Line, point.Col)
	}
}

func TestGoToTopAndBottom(t *testing.T) {
	t.Parallel()

	ws := workspace.New()
	ed := ws.Editor
	buf := ed.Scratch("*scratch*")
	buf.Text.SetBytes([]byte("one\ntwo\nthree"))

	frame, err := ws.NewFrame(80, 10, "")
	if err != nil {
		t.Fatal(err)
	}

	win := ws.Windows[frame.ActiveWindowID]

	editor.SetWindowCursorForTest(buf, win, buf.Text.LineStart(2))

	if err := ws.HandleInput(frame.ID, "g"); err != nil {
		t.Fatal(err)
	}
	if err := ws.HandleInput(frame.ID, "g"); err != nil {
		t.Fatal(err)
	}

	point := buf.Text.PointOf(editor.WindowCursorForTest(win))
	if point.Line != 0 {
		t.Fatalf("after gg: line = %d, want 0", point.Line)
	}

	if err := ws.HandleInput(frame.ID, "G"); err != nil {
		t.Fatal(err)
	}

	point = buf.Text.PointOf(editor.WindowCursorForTest(win))
	if point.Line != 2 {
		t.Fatalf("after G: line = %d, want 2", point.Line)
	}
}

func TestInsertAtFirstNonBlank(t *testing.T) {
	t.Parallel()

	ws := workspace.New()
	ed := ws.Editor
	buf := ed.Scratch("*scratch*")
	buf.Text.SetBytes([]byte("  hi"))

	frame, _ := ws.NewFrame(80, 10, "")
	win := ws.Windows[frame.ActiveWindowID]

	editor.SetWindowCursorForTest(buf, win, buf.Text.LineEnd(0))

	if err := ws.HandleInput(frame.ID, "I"); err != nil {
		t.Fatal(err)
	}

	if ed.Mode != editor.InsertMode {
		t.Fatal("expected insert mode after I")
	}

	point := buf.Text.PointOf(editor.WindowCursorForTest(win))
	if point.Col != 2 {
		t.Fatalf("cursor after I = col %d, want 2", point.Col)
	}

	if err := ws.HandleInput(frame.ID, "X"); err != nil {
		t.Fatal(err)
	}

	if string(buf.Text.Bytes()) != "  Xhi" {
		t.Fatalf("got %q, want %q", buf.Text.Bytes(), "  Xhi")
	}
}

func TestLineStartAndEndKeys(t *testing.T) {
	t.Parallel()

	ws := workspace.New()
	ed := ws.Editor
	buf := ed.Scratch("*scratch*")
	buf.Text.SetBytes([]byte("hello"))

	frame, _ := ws.NewFrame(80, 10, "")
	win := ws.Windows[frame.ActiveWindowID]

	if err := ws.HandleInput(frame.ID, "$"); err != nil {
		t.Fatal(err)
	}

	point := buf.Text.PointOf(editor.WindowCursorForTest(win))
	if point.Col != 4 {
		t.Fatalf("after $: col = %d, want 4", point.Col)
	}

	if err := ws.HandleInput(frame.ID, "0"); err != nil {
		t.Fatal(err)
	}

	point = buf.Text.PointOf(editor.WindowCursorForTest(win))
	if point.Col != 0 {
		t.Fatalf("after 0: col = %d, want 0", point.Col)
	}
}

func TestAppendEndOfLineEntersInsertAtEnd(t *testing.T) {
	t.Parallel()

	ws := workspace.New()
	ed := ws.Editor
	buf := ed.Scratch("*scratch*")
	buf.Text.SetBytes([]byte("hi"))

	frame, _ := ws.NewFrame(80, 10, "")
	win := ws.Windows[frame.ActiveWindowID]

	if err := ws.HandleInput(frame.ID, "A"); err != nil {
		t.Fatal(err)
	}

	if ed.Mode != editor.InsertMode {
		t.Fatal("expected insert mode after A")
	}

	if got := editor.WindowCursorForTest(win); got != 2 {
		t.Fatalf("cursor after A = %d, want 2", got)
	}

	if err := ws.HandleInput(frame.ID, "!"); err != nil {
		t.Fatal(err)
	}

	if string(buf.Text.Bytes()) != "hi!" {
		t.Fatalf("got %q, want %q", buf.Text.Bytes(), "hi!")
	}
}

func TestAppendMovesCursorPastLastChar(t *testing.T) {
	t.Parallel()

	ws := workspace.New()
	ed := ws.Editor
	buf := ed.Scratch("*scratch*")
	buf.Text.SetBytes([]byte("hi"))

	frame, _ := ws.NewFrame(80, 10, "")
	win := ws.Windows[frame.ActiveWindowID]

	_ = ws.HandleInput(frame.ID, "l")
	_ = ws.HandleInput(frame.ID, "a")

	if got := editor.WindowCursorForTest(win); got != 2 {
		t.Fatalf("cursor after a = %d, want 2", got)
	}
}

func TestInsertModeEditsScratchBuffer(t *testing.T) {
	t.Parallel()

	ws := workspace.New()
	ed := ws.Editor
	buf := ed.Scratch("*scratch*")
	buf.Text.SetBytes([]byte("hi"))

	frame, err := ws.NewFrame(80, 10, "")
	if err != nil {
		t.Fatal(err)
	}

	if err := ws.HandleInput(frame.ID, "l"); err != nil {
		t.Fatal(err)
	}

	if err := ws.HandleInput(frame.ID, "a"); err != nil {
		t.Fatal(err)
	}

	if ed.Mode != editor.InsertMode {
		t.Fatal("expected insert mode")
	}

	if err := ws.HandleInput(frame.ID, "!"); err != nil {
		t.Fatal(err)
	}

	if err := ws.HandleInput(frame.ID, "<Space>"); err != nil {
		t.Fatal(err)
	}

	if err := ws.HandleInput(frame.ID, "x"); err != nil {
		t.Fatal(err)
	}

	if string(buf.Text.Bytes()) != "hi! x" {
		t.Fatalf("got %q, want %q", buf.Text.Bytes(), "hi! x")
	}

	if err := ws.HandleInput(frame.ID, "<Esc>"); err != nil {
		t.Fatal(err)
	}

	if ed.Mode != editor.NormalMode {
		t.Fatal("expected normal mode after Esc")
	}
}
