package editor_test

import (
	"testing"

	"vague/bonefire/workspace"
)

func TestRenderRedrawWrapsWithoutPadding(t *testing.T) {
	t.Parallel()

	ws := workspace.New()
	ed := ws.Editor
	buf := ed.Scratch("*scratch*")
	buf.Text.SetBytes([]byte("abcdefghijklmn"))

	frame, err := ws.NewFrame(13, 4, "")
	if err != nil {
		t.Fatal(err)
	}

	redraw, ok := ws.RenderRedraw(frame.ID)
	if !ok {
		t.Fatal("expected redraw")
	}

	want := []string{"abcdefghijklm", "n"}
	if len(redraw.Lines) != len(want) {
		t.Fatalf("got %v, want %v", redraw.Lines, want)
	}

	for i := range want {
		if redraw.Lines[i] != want[i] {
			t.Errorf("line %d = %q, want %q", i, redraw.Lines[i], want[i])
		}
	}

	if redraw.Wrap != true {
		t.Fatal("expected wrap enabled")
	}
}

func TestRenderRedrawNoWrapTruncates(t *testing.T) {
	t.Parallel()

	ws := workspace.New()
	ed := ws.Editor
	buf := ed.Scratch("*scratch*")
	buf.Text.SetBytes([]byte("abcdefghijklmn"))

	frame, err := ws.NewFrame(13, 4, "")
	if err != nil {
		t.Fatal(err)
	}

	win := ws.Windows[frame.ActiveWindowID]
	win.WindowOptions.Wrap = false

	redraw, ok := ws.RenderRedraw(frame.ID)
	if !ok {
		t.Fatal("expected redraw")
	}

	if len(redraw.Lines) != 1 || redraw.Lines[0] != "abcdefghijklm" {
		t.Fatalf("got %v, want truncated line", redraw.Lines)
	}
}

func TestRenderRedrawModified(t *testing.T) {
	t.Parallel()

	ws := workspace.New()
	ed := ws.Editor
	buf := ed.Scratch("*scratch*")
	buf.Text.SetBytes([]byte("hi"))

	frame, err := ws.NewFrame(80, 10, "")
	if err != nil {
		t.Fatal(err)
	}

	redraw, ok := ws.RenderRedraw(frame.ID)
	if !ok {
		t.Fatal("expected redraw")
	}

	if redraw.Buffer.Modified {
		t.Fatal("expected clean buffer before edit")
	}

	ws.InvalidateFrame(frame.ID)
	if _, err := buf.Insert(buf.Text.Len(), []byte("!")); err != nil {
		t.Fatal(err)
	}

	redraw, ok = ws.RenderRedraw(frame.ID)
	if !ok {
		t.Fatal("expected redraw after edit")
	}

	if !redraw.Buffer.Modified {
		t.Fatal("expected modified buffer after edit")
	}
}
