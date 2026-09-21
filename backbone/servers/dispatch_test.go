package backbone_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	client "vague/backbone/clients"
	"vague/backbone/process"
	server "vague/backbone/servers"
)

func TestReadyPushesScratchBufferRedraw(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	listener, err := process.Listen()
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()

	srv := server.NewServer()
	go func() {
		_ = srv.Serve(ctx, listener)
	}()

	c, err := client.ClientConnect()
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	if _, err := c.UiAttach(ctx, process.UiAttachParams{}); err != nil {
		t.Fatal(err)
	}

	readyDone := make(chan struct{})
	go func() {
		if _, err := c.UiReady(ctx, process.UiReadyParams{Height: 24, Width: 80}); err != nil {
			t.Error(err)
		}
		close(readyDone)
	}()

	select {
	case note := <-c.Notifications():
		if note.Method != process.MethodRedraw {
			t.Fatalf("expected redraw, got %q", note.Method)
		}

		var redraw process.Redraw
		if err := json.Unmarshal(note.Params, &redraw); err != nil {
			t.Fatal(err)
		}

		if redraw.Buffer.Name != "*scratch*" {
			t.Fatalf("expected scratch buffer, got %q", redraw.Buffer.Name)
		}

		if len(redraw.Lines) == 0 {
			t.Fatal("expected redraw lines")
		}

		if redraw.Columns != 80 {
			t.Fatalf("expected 80 columns, got %d", redraw.Columns)
		}

		if !redraw.Wrap {
			t.Fatal("expected wrap enabled by default")
		}

		if !redraw.Full {
			t.Fatal("expected full redraw on first ready")
		}

	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for redraw")
	}

	<-readyDone
}
