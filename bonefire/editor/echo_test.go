package editor

import "testing"

func TestSetEchoIncludedInRedraw(t *testing.T) {
	t.Parallel()

	ed := NewEditor()
	ed.Scratch("*scratch*")

	frame, err := ed.NewFrame(80, 10)
	if err != nil {
		t.Fatal(err)
	}

	if err := ed.SetEcho(frame.ID, `"foo" written`, EchoInfo); err != nil {
		t.Fatal(err)
	}

	redraw, ok := ed.RenderRedraw(frame.ID)
	if !ok {
		t.Fatal("expected redraw")
	}

	if redraw.Echo == nil || redraw.Echo.Message != `"foo" written` {
		t.Fatalf("echo = %#v, want written message", redraw.Echo)
	}
}

func TestClearEchoOnRedraw(t *testing.T) {
	t.Parallel()

	ed := NewEditor()
	ed.Scratch("*scratch*")

	frame, err := ed.NewFrame(80, 10)
	if err != nil {
		t.Fatal(err)
	}

	_ = ed.SetEcho(frame.ID, "wrap on", EchoInfo)
	ed.ClearEcho(frame.ID)

	redraw, ok := ed.RenderRedraw(frame.ID)
	if !ok {
		t.Fatal("expected redraw")
	}

	if redraw.Echo != nil {
		t.Fatalf("echo = %#v, want nil", redraw.Echo)
	}
}

func TestWriteEchoMessageUsesPath(t *testing.T) {
	t.Parallel()

	buf := NewEditor().Scratch("*scratch*")
	buf.Path = "/tmp/example.txt"
	buf.Name = "example.txt"

	if got := WriteEchoMessage(buf); got != `"/tmp/example.txt" written` {
		t.Fatalf("got %q", got)
	}
}
