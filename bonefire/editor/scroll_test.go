package editor

import (
	"fmt"
	"strings"
	"testing"

	"vague/bonefire/text"
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
	view := layoutView(tex, 80, 0, false)

	win := newWindow(1, 1, 1)
	frame := &Frame{Height: 10}

	point := text.Point{Line: 29, Col: 0}
	ensureCursorVisible(win, frame, view, point)

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
	view := layoutView(tex, 80, 0, false)

	win := newWindow(1, 1, 1)
	win.TopLine = 20
	frame := &Frame{Height: 10}

	point := text.Point{Line: 5, Col: 0}
	ensureCursorVisible(win, frame, view, point)

	if win.TopLine != 5 {
		t.Fatalf("TopLine = %d, want 5", win.TopLine)
	}
}

func TestRenderRedrawScrollsToCursor(t *testing.T) {
	t.Parallel()

	ed := NewEditor()
	buf := ed.Scratch("*scratch*")

	var lines strings.Builder
	for i := range 30 {
		if i > 0 {
			lines.WriteByte('\n')
		}
		fmt.Fprintf(&lines, "line %d", i)
	}
	buf.Text.SetBytes([]byte(lines.String()))

	frame, err := ed.NewFrame(80, 10, "")
	if err != nil {
		t.Fatal(err)
	}

	win := ed.Windows[frame.ActiveWindowID]
	setWindowCursor(buf, win, buf.Text.LineEnd(29))

	redraw, ok := ed.RenderRedraw(frame.ID)
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

	ed := NewEditor()
	buf := ed.Scratch("*scratch*")

	var lines strings.Builder
	for i := range 30 {
		if i > 0 {
			lines.WriteByte('\n')
		}
		fmt.Fprintf(&lines, "line %d", i)
	}
	buf.Text.SetBytes([]byte(lines.String()))

	frame, err := ed.NewFrame(80, 10, "")
	if err != nil {
		t.Fatal(err)
	}

	win := ed.Windows[frame.ActiveWindowID]

	for range 29 {
		if err := ed.HandleInput(frame.ID, "j"); err != nil {
			t.Fatal(err)
		}
	}

	ed.InvalidateFrame(frame.ID)
	redraw, ok := ed.RenderRedraw(frame.ID)
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
		if err := ed.HandleInput(frame.ID, "k"); err != nil {
			t.Fatal(err)
		}
	}

	ed.InvalidateFrame(frame.ID)
	redraw, ok = ed.RenderRedraw(frame.ID)
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
