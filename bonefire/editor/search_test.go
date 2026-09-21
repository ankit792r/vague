package editor_test

import (
	"testing"

	"vague/bonefire/editor"
	"vague/bonefire/workspace"
)

func TestSearchForward(t *testing.T) {
	t.Parallel()

	ws := workspace.New()
	buf := ws.Editor.Scratch("*scratch*")
	buf.Text.SetBytes([]byte("foo bar foo\n"))

	frame, err := ws.NewFrame(80, 10, "")
	if err != nil {
		t.Fatal(err)
	}

	if err := ws.Search(frame.ID, "foo", true); err != nil {
		t.Fatal(err)
	}
	if got := editor.WindowCursorForTest(ws.Windows[frame.ActiveWindowID]); got != 8 {
		t.Fatalf("first search cursor = %d, want 8", got)
	}

	if err := ws.Search(frame.ID, "foo", true); err != nil {
		t.Fatal(err)
	}
	if got := editor.WindowCursorForTest(ws.Windows[frame.ActiveWindowID]); got != 0 {
		t.Fatalf("wrapped search cursor = %d, want 0", got)
	}
}

func TestSearchWrapForward(t *testing.T) {
	t.Parallel()

	ws := workspace.New()
	buf := ws.Editor.Scratch("*scratch*")
	buf.Text.SetBytes([]byte("foo bar\n"))

	frame, err := ws.NewFrame(80, 10, "")
	if err != nil {
		t.Fatal(err)
	}

	if err := ws.HandleInput(frame.ID, "$"); err != nil {
		t.Fatal(err)
	}
	if err := ws.Search(frame.ID, "foo", true); err != nil {
		t.Fatal(err)
	}
	if got := editor.WindowCursorForTest(ws.Windows[frame.ActiveWindowID]); got != 0 {
		t.Fatalf("wrapped search cursor = %d, want 0", got)
	}
}

func TestSearchBackward(t *testing.T) {
	t.Parallel()

	ws := workspace.New()
	buf := ws.Editor.Scratch("*scratch*")
	buf.Text.SetBytes([]byte("foo bar foo\n"))

	frame, err := ws.NewFrame(80, 10, "")
	if err != nil {
		t.Fatal(err)
	}

	if err := ws.HandleInput(frame.ID, "$"); err != nil {
		t.Fatal(err)
	}
	if err := ws.Search(frame.ID, "foo", false); err != nil {
		t.Fatal(err)
	}
	if got := editor.WindowCursorForTest(ws.Windows[frame.ActiveWindowID]); got != 8 {
		t.Fatalf("backward search cursor = %d, want 8", got)
	}
}

func TestRepeatSearchNAndN(t *testing.T) {
	t.Parallel()

	ws := workspace.New()
	buf := ws.Editor.Scratch("*scratch*")
	buf.Text.SetBytes([]byte("aa aa aa\n"))

	frame, err := ws.NewFrame(80, 10, "")
	if err != nil {
		t.Fatal(err)
	}

	if err := ws.Search(frame.ID, "aa", true); err != nil {
		t.Fatal(err)
	}
	first := editor.WindowCursorForTest(ws.Windows[frame.ActiveWindowID])

	if err := ws.HandleInput(frame.ID, "n"); err != nil {
		t.Fatal(err)
	}
	second := editor.WindowCursorForTest(ws.Windows[frame.ActiveWindowID])
	if second == first {
		t.Fatal("n did not move cursor")
	}

	if err := ws.HandleInput(frame.ID, "N"); err != nil {
		t.Fatal(err)
	}
	back := editor.WindowCursorForTest(ws.Windows[frame.ActiveWindowID])
	if back != first {
		t.Fatalf("N cursor = %d, want %d", back, first)
	}
}

func TestSearchNotFound(t *testing.T) {
	t.Parallel()

	ws := workspace.New()
	buf := ws.Editor.Scratch("*scratch*")
	buf.Text.SetBytes([]byte("hello\n"))

	frame, err := ws.NewFrame(80, 10, "")
	if err != nil {
		t.Fatal(err)
	}

	err = ws.Search(frame.ID, "zzz", true)
	if err != editor.ErrPatternNotFound {
		t.Fatalf("err = %v, want ErrPatternNotFound", err)
	}
}
