package editor_test

import (
	"testing"

	"vague/bonefire/editor"
	"vague/bonefire/workspace"
)

func TestRenderRedrawLineNumbers(t *testing.T) {
	t.Parallel()

	ws := workspace.New()
	ed := ws.Editor
	buf := ed.Scratch("*scratch*")
	buf.Text.SetBytes([]byte("one\ntwo\nthree\n"))

	frame, err := ws.NewFrame(20, 10, "")
	if err != nil {
		t.Fatal(err)
	}

	win := ws.Windows[frame.ActiveWindowID]
	win.WindowOptions.Number = true

	ws.InvalidateFrame(frame.ID)
	redraw, ok := ws.RenderRedraw(frame.ID)
	if !ok {
		t.Fatal("expected redraw")
	}

	if !redraw.Number {
		t.Fatal("expected number on redraw")
	}
	if redraw.GutterColumns < 3 {
		t.Fatalf("gutter_columns = %d, want at least 3", redraw.GutterColumns)
	}
	if len(redraw.LineNumbers) != len(redraw.Lines) {
		t.Fatalf("line_numbers len %d != lines len %d", len(redraw.LineNumbers), len(redraw.Lines))
	}
	if len(redraw.LineNumbers) < 3 {
		t.Fatalf("expected at least 3 visible rows, got %d", len(redraw.LineNumbers))
	}
	if redraw.LineNumbers[0] != 1 || redraw.LineNumbers[1] != 2 || redraw.LineNumbers[2] != 3 {
		t.Fatalf("line_numbers = %v", redraw.LineNumbers)
	}

	// Gutter steals columns from wrap width.
	long := "abcdefghijklmnopqrstuvwxyz"
	buf.Text.SetBytes([]byte(long))
	ws.InvalidateFrame(frame.ID)
	redraw, ok = ws.RenderRedraw(frame.ID)
	if !ok {
		t.Fatal("expected redraw")
	}
	contentWidth := 20 - redraw.GutterColumns
	if contentWidth < 1 {
		contentWidth = 1
	}
	if len(redraw.Lines[0]) > contentWidth {
		t.Fatalf("line width %d > content width %d with gutter", len(redraw.Lines[0]), contentWidth)
	}
}

func TestGutterColumns(t *testing.T) {
	t.Parallel()

	if got := editor.GutterColumnsForTest(1); got != 3 {
		t.Fatalf("gutter for 1 line = %d, want 3", got)
	}
	if got := editor.GutterColumnsForTest(99); got != 3 {
		t.Fatalf("gutter for 99 lines = %d, want 3", got)
	}
	if got := editor.GutterColumnsForTest(100); got != 4 {
		t.Fatalf("gutter for 100 lines = %d, want 4", got)
	}
}
