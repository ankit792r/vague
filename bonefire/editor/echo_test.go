package editor_test

import (
	"vague/bonefire/editor"
	"testing"

	"vague/bonefire/workspace"
)

func TestSetEchoIncludedInRedraw(t *testing.T) {
	t.Parallel()

	ws := workspace.New()
	ed := ws.Editor
	ed.Scratch("*scratch*")

	frame, err := ws.NewFrame(80, 10, "")
	if err != nil {
		t.Fatal(err)
	}

	if err := ws.SetEcho(frame.ID, `"foo" written`, editor.EchoInfo); err != nil {
		t.Fatal(err)
	}

	redraw, ok := ws.RenderRedraw(frame.ID)
	if !ok {
		t.Fatal("expected redraw")
	}

	if redraw.Echo == nil || redraw.Echo.Message != `"foo" written` {
		t.Fatalf("echo = %#v, want written message", redraw.Echo)
	}
}

func TestClearEchoOnRedraw(t *testing.T) {
	t.Parallel()

	ws := workspace.New()
	ed := ws.Editor
	ed.Scratch("*scratch*")

	frame, err := ws.NewFrame(80, 10, "")
	if err != nil {
		t.Fatal(err)
	}

	_ = ws.SetEcho(frame.ID, "wrap on", editor.EchoInfo)
	ws.ClearEcho(frame.ID)

	redraw, ok := ws.RenderRedraw(frame.ID)
	if !ok {
		t.Fatal("expected redraw")
	}

	if redraw.Echo != nil {
		t.Fatalf("echo = %#v, want nil", redraw.Echo)
	}
}

func TestWriteEchoMessageUsesPath(t *testing.T) {
	t.Parallel()

	buf := editor.NewEditor().Scratch("*scratch*")
	buf.Path = "/tmp/example.txt"
	buf.Name = "example.txt"

	if got := editor.WriteEchoMessage(buf); got != `"/tmp/example.txt" written` {
		t.Fatalf("got %q", got)
	}
}
