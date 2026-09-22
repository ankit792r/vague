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

func TestSearchHighlightInRedraw(t *testing.T) {
	t.Parallel()

	ws := workspace.New()
	buf := ws.Editor.Scratch("*scratch*")
	buf.Text.SetBytes([]byte("hello world"))

	frame, err := ws.NewFrame(80, 10, "")
	if err != nil {
		t.Fatal(err)
	}

	if err := ws.Search(frame.ID, "world", true); err != nil {
		t.Fatal(err)
	}

	redraw, ok := ws.RenderRedraw(frame.ID)
	if !ok {
		t.Fatal("expected redraw")
	}
	if redraw.SearchMatch == nil || !redraw.SearchMatch.Visible {
		t.Fatal("expected search_match highlight")
	}
	if redraw.SearchMatch.Start.Column != 6 || redraw.SearchMatch.End.Column != 10 {
		t.Fatalf("search_match cols = %d..%d",
			redraw.SearchMatch.Start.Column, redraw.SearchMatch.End.Column)
	}
}

func TestSearchHighlightMovesWithN(t *testing.T) {
	t.Parallel()

	ws := workspace.New()
	buf := ws.Editor.Scratch("*scratch*")
	buf.Text.SetBytes([]byte("foo foo foo"))

	frame, err := ws.NewFrame(80, 10, "")
	if err != nil {
		t.Fatal(err)
	}

	if err := ws.Search(frame.ID, "foo", true); err != nil {
		t.Fatal(err)
	}
	if err := ws.HandleInput(frame.ID, "n"); err != nil {
		t.Fatal(err)
	}

	redraw, ok := ws.RenderRedraw(frame.ID)
	if !ok {
		t.Fatal("expected redraw")
	}
	if redraw.SearchMatch == nil {
		t.Fatal("expected search_match")
	}
	if redraw.SearchMatch.Start.Column != 8 {
		t.Fatalf("search_match start col = %d, want 8", redraw.SearchMatch.Start.Column)
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

func TestHlsearchHighlightsInRedraw(t *testing.T) {
	t.Parallel()

	ws := workspace.New()
	buf := ws.Editor.Scratch("*scratch*")
	buf.Text.SetBytes([]byte("aa bb aa"))

	frame, err := ws.NewFrame(80, 10, "")
	if err != nil {
		t.Fatal(err)
	}

	if _, err := ws.RunExLine(frame.ID, "set hlsearch"); err != nil {
		t.Fatal(err)
	}
	if err := ws.Search(frame.ID, "aa", true); err != nil {
		t.Fatal(err)
	}

	redraw, ok := ws.RenderRedraw(frame.ID)
	if !ok {
		t.Fatal("expected redraw")
	}
	if len(redraw.SearchHighlights) == 0 {
		t.Fatal("expected search_highlights with hlsearch on")
	}
}
