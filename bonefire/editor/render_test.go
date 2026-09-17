package editor

import "testing"

func TestRenderRedrawWrapsWithoutPadding(t *testing.T) {
	t.Parallel()

	ed := NewEditor()
	buf := ed.Scratch("*scratch*")
	buf.Text.SetBytes([]byte("abcdefghijklmn"))

	frame, err := ed.NewFrame(13, 4)
	if err != nil {
		t.Fatal(err)
	}

	redraw, ok := ed.RenderRedraw(frame.ID)
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

	ed := NewEditor()
	buf := ed.Scratch("*scratch*")
	buf.Text.SetBytes([]byte("abcdefghijklmn"))

	frame, err := ed.NewFrame(13, 4)
	if err != nil {
		t.Fatal(err)
	}

	win := ed.Windows[frame.ActiveWindowID]
	win.WindowOptions.Wrap = false

	redraw, ok := ed.RenderRedraw(frame.ID)
	if !ok {
		t.Fatal("expected redraw")
	}

	if len(redraw.Lines) != 1 || redraw.Lines[0] != "abcdefghijklm" {
		t.Fatalf("got %v, want truncated line", redraw.Lines)
	}
}
