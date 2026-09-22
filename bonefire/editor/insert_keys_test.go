package editor_test

import (
	"strings"
	"testing"

	"vague/bonefire/editor"
	"vague/bonefire/workspace"
)

func TestInsertCtrlWAndCtrlU(t *testing.T) {
	t.Parallel()
	ws := workspace.New()
	buf := ws.Editor.Scratch("*scratch*")
	buf.Text.SetBytes([]byte("hello world"))
	frame, _ := ws.NewFrame(80, 10, "")
	win := ws.Windows[frame.ActiveWindowID]
	editor.SetWindowCursorForTest(buf, win, buf.Text.Len())

	_ = ws.HandleInput(frame.ID, "i")
	_ = ws.HandleInput(frame.ID, "<C-w>")
	if got := string(buf.Text.Bytes()); got != "hello " {
		t.Fatalf("after C-w: %q", got)
	}

	buf.Text.SetBytes([]byte("  abcd"))
	editor.SetWindowCursorForTest(buf, win, 5)
	_ = ws.HandleInput(frame.ID, "i")
	_ = ws.HandleInput(frame.ID, "<C-u>")
	if got := string(buf.Text.Bytes()); got != "d" {
		t.Fatalf("after C-u: %q want %q", got, "d")
	}
}

func TestInsertCtrlAAndCtrlE(t *testing.T) {
	t.Parallel()
	ws := workspace.New()
	buf := ws.Editor.Scratch("*scratch*")
	buf.Text.SetBytes([]byte("hello"))
	frame, _ := ws.NewFrame(80, 10, "")
	win := ws.Windows[frame.ActiveWindowID]
	editor.SetWindowCursorForTest(buf, win, 3)
	_ = ws.HandleInput(frame.ID, "i")
	_ = ws.HandleInput(frame.ID, "<C-a>")
	if col := buf.Text.PointOf(editor.WindowCursorForTest(win)).Col; col != 0 {
		t.Fatalf("C-a col = %d", col)
	}
	_ = ws.HandleInput(frame.ID, "<C-e>")
	if col := buf.Text.PointOf(editor.WindowCursorForTest(win)).Col; col != 5 {
		t.Fatalf("C-e col = %d want 5", col)
	}
}

func TestInsertCtrlONormalOnce(t *testing.T) {
	t.Parallel()
	ws := workspace.New()
	buf := ws.Editor.Scratch("*scratch*")
	buf.Text.SetBytes([]byte("hello"))
	frame, _ := ws.NewFrame(80, 10, "")
	win := ws.Windows[frame.ActiveWindowID]
	editor.SetWindowCursorForTest(buf, win, 0)
	_ = ws.HandleInput(frame.ID, "i")
	_ = ws.HandleInput(frame.ID, "<C-o>")
	if ws.Editor.Mode != editor.NormalMode {
		t.Fatal("expected normal during C-o")
	}
	_ = ws.HandleInput(frame.ID, "l")
	if ws.Editor.Mode != editor.InsertMode {
		t.Fatal("expected insert after one normal key")
	}
	if col := buf.Text.PointOf(editor.WindowCursorForTest(win)).Col; col != 1 {
		t.Fatalf("col = %d", col)
	}
}

func yankLineToRegA(t *testing.T, ws *workspace.Workspace, frameID uint64) {
	t.Helper()
	for _, k := range []string{`"`, "a", "y", "y"} {
		if err := ws.HandleInput(frameID, k); err != nil {
			t.Fatal(err)
		}
	}
}

func TestInsertCtrlRRegister(t *testing.T) {
	t.Parallel()
	ws := workspace.New()
	buf := ws.Editor.Scratch("*scratch*")
	buf.Text.SetBytes([]byte("hi"))
	frame, _ := ws.NewFrame(80, 10, "")
	yankLineToRegA(t, ws, frame.ID)
	_ = ws.HandleInput(frame.ID, "A")
	_ = ws.HandleInput(frame.ID, "<C-r>")
	_ = ws.HandleInput(frame.ID, "a")
	if got := string(buf.Text.Bytes()); got != "hihi" {
		t.Fatalf("C-r a: %q", got)
	}
}

func TestInsertCtrlTIndent(t *testing.T) {
	t.Parallel()
	ws := workspace.New()
	buf := ws.Editor.Scratch("*scratch*")
	buf.Text.SetBytes([]byte("x"))
	frame, _ := ws.NewFrame(80, 10, "")
	_ = ws.HandleInput(frame.ID, "i")
	_ = ws.HandleInput(frame.ID, "<C-t>")
	if got := string(buf.Text.Bytes()); got != "  x" {
		t.Fatalf("C-t: %q", got)
	}
}

func TestInsertCompletionCtrlXNP(t *testing.T) {
	t.Parallel()
	ws := workspace.New()
	buf := ws.Editor.Scratch("*scratch*")
	buf.Text.SetBytes([]byte("alpha alps"))
	frame, _ := ws.NewFrame(80, 10, "")
	win := ws.Windows[frame.ActiveWindowID]
	editor.SetWindowCursorForTest(buf, win, buf.Text.Len())
	_ = ws.HandleInput(frame.ID, "i")
	_ = ws.HandleInput(frame.ID, " al")
	_ = ws.HandleInput(frame.ID, "<C-x>")
	_ = ws.HandleInput(frame.ID, "<C-n>")
	got := string(buf.Text.Bytes())
	if !strings.Contains(got, "alpha") || !strings.Contains(got, "alps") {
		t.Fatalf("completion: %q", got)
	}
}

func TestReplaceModeAndVisualReplace(t *testing.T) {
	t.Parallel()
	ws := workspace.New()
	buf := ws.Editor.Scratch("*scratch*")
	buf.Text.SetBytes([]byte("abc"))
	frame, _ := ws.NewFrame(80, 10, "")
	win := ws.Windows[frame.ActiveWindowID]
	editor.SetWindowCursorForTest(buf, win, 0)
	_ = ws.HandleInput(frame.ID, "R")
	_ = ws.HandleInput(frame.ID, "x")
	_ = ws.HandleInput(frame.ID, "y")
	if got := string(buf.Text.Bytes()); got != "xyc" {
		t.Fatalf("replace mode: %q", got)
	}

	_ = ws.HandleInput(frame.ID, "<Esc>")

	buf.Text.SetBytes([]byte("abc"))
	editor.SetWindowCursorForTest(buf, win, 0)
	_ = ws.HandleInput(frame.ID, "v")
	_ = ws.HandleInput(frame.ID, "l")
	_ = ws.HandleInput(frame.ID, "r")
	_ = ws.HandleInput(frame.ID, "z")
	if got := string(buf.Text.Bytes()); got != "zzc" {
		t.Fatalf("visual r: %q", got)
	}
}

func TestSetPasteInsertRegister(t *testing.T) {
	t.Parallel()
	ws := workspace.New()
	buf := ws.Editor.Scratch("*scratch*")
	buf.Text.SetBytes([]byte("hi"))
	frame, _ := ws.NewFrame(80, 10, "")
	_, _ = ws.RunExLine(frame.ID, "set paste")
	yankLineToRegA(t, ws, frame.ID)
	_ = ws.HandleInput(frame.ID, "A")
	_ = ws.HandleInput(frame.ID, "<C-r>")
	_ = ws.HandleInput(frame.ID, "a")
	if got := string(buf.Text.Bytes()); got != "hihi" {
		t.Fatalf("paste insert: %q", got)
	}
}
