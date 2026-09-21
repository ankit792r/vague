package editor_test

import (
	"fmt"
	"strings"
	"testing"

	"vague/bonefire/display"
	"vague/bonefire/editor"
	"vague/bonefire/text"
	"vague/bonefire/window"
	"vague/bonefire/workspace"
)

func TestEnsureCursorVisibleScrollsDown(t *testing.T) {
	t.Parallel()

	var lines strings.Builder
	for i := range 30 {
		if i > 0 {
			lines.WriteByte('\n')
		}
		fmt.Fprintf(&lines, "line %d", i)
	}

	tex := text.New([]byte(lines.String()))
	view := editor.LayoutViewForTest(tex, 80, 0, false)

	win := window.NewWindow(1, 1)
	frame := &display.Frame{Height: 10, Width: 80}

	point := text.Point{Line: 29, Col: 0}
	editor.EnsureCursorVisibleForTest(win, frame, view, point)

	if win.TopLine != 20 {
		t.Fatalf("TopLine = %d, want 20", win.TopLine)
	}
}

func TestEnsureCursorVisibleScrollsUp(t *testing.T) {
	t.Parallel()

	var lines strings.Builder
	for i := range 30 {
		if i > 0 {
			lines.WriteByte('\n')
		}
		fmt.Fprintf(&lines, "line %d", i)
	}

	tex := text.New([]byte(lines.String()))
	view := editor.LayoutViewForTest(tex, 80, 0, false)

	win := window.NewWindow(1, 1)
	win.TopLine = 20
	frame := &display.Frame{Height: 10, Width: 80}

	point := text.Point{Line: 5, Col: 0}
	editor.EnsureCursorVisibleForTest(win, frame, view, point)

	if win.TopLine != 5 {
		t.Fatalf("TopLine = %d, want 5", win.TopLine)
	}
}

func TestRenderRedrawScrollsToCursor(t *testing.T) {
	t.Parallel()

	ws := workspace.New()
	ed := ws.Editor
	buf := ed.Scratch("*scratch*")

	var lines strings.Builder
	for i := range 30 {
		if i > 0 {
			lines.WriteByte('\n')
		}
		fmt.Fprintf(&lines, "line %d", i)
	}
	buf.Text.SetBytes([]byte(lines.String()))

	frame, err := ws.NewFrame(80, 10, "")
	if err != nil {
		t.Fatal(err)
	}

	win := ws.Windows[frame.ActiveWindowID]
	editor.SetWindowCursorForTest(buf, win, buf.Text.LineEnd(29))

	redraw, ok := ws.RenderRedraw(frame.ID)
	if !ok {
		t.Fatal("expected redraw")
	}

	if !redraw.Cursor.Visible {
		t.Fatal("cursor not visible")
	}

	if redraw.Cursor.Row != 9 {
		t.Fatalf("cursor row = %d, want 9", redraw.Cursor.Row)
	}

	if len(redraw.Lines) != 10 {
		t.Fatalf("got %d lines, want 10", len(redraw.Lines))
	}

	if redraw.Lines[0] != "line 20" {
		t.Fatalf("first line = %q, want %q", redraw.Lines[0], "line 20")
	}
}

func TestHandleInputScrollsViewport(t *testing.T) {
	t.Parallel()

	ws := workspace.New()
	ed := ws.Editor
	buf := ed.Scratch("*scratch*")

	var lines strings.Builder
	for i := range 30 {
		if i > 0 {
			lines.WriteByte('\n')
		}
		fmt.Fprintf(&lines, "line %d", i)
	}
	buf.Text.SetBytes([]byte(lines.String()))

	frame, err := ws.NewFrame(80, 10, "")
	if err != nil {
		t.Fatal(err)
	}

	win := ws.Windows[frame.ActiveWindowID]

	for range 29 {
		if err := ws.HandleInput(frame.ID, "j"); err != nil {
			t.Fatal(err)
		}
	}

	ws.InvalidateFrame(frame.ID)
	redraw, ok := ws.RenderRedraw(frame.ID)
	if !ok {
		t.Fatal("expected redraw")
	}

	if !redraw.Cursor.Visible {
		t.Fatal("cursor not visible after j to bottom")
	}

	if redraw.Cursor.Row != 9 {
		t.Fatalf("cursor row = %d, want 9", redraw.Cursor.Row)
	}

	for range 10 {
		if err := ws.HandleInput(frame.ID, "k"); err != nil {
			t.Fatal(err)
		}
	}

	ws.InvalidateFrame(frame.ID)
	redraw, ok = ws.RenderRedraw(frame.ID)
	if !ok {
		t.Fatal("expected redraw")
	}

	if !redraw.Cursor.Visible {
		t.Fatal("cursor not visible after k up")
	}

	if redraw.Lines[0] != "line 19" {
		t.Fatalf("first line = %q, want %q", redraw.Lines[0], "line 19")
	}

	if win.TopLine != 19 {
		t.Fatalf("TopLine = %d, want 19", win.TopLine)
	}
}
