package process_test

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"vague/backbone/process"
)

func TestPingWhenServerNotRunning(t *testing.T) {
	t.Parallel()

	if os.Getenv("XDG_RUNTIME_DIR") == "" {
		t.Skip("XDG_RUNTIME_DIR is not set")
	}

	if err := process.Ping(); err == nil {
		t.Skip("vague server is running")
	}

	err := process.Ping()
	if !errors.Is(err, process.ErrNotRunning) {
		t.Fatalf("Ping() = %v, want ErrNotRunning", err)
	}
}

func TestWaitUntilRunningTimesOut(t *testing.T) {
	t.Parallel()

	if os.Getenv("XDG_RUNTIME_DIR") == "" {
		t.Skip("XDG_RUNTIME_DIR is not set")
	}

	if process.Ping() == nil {
		t.Skip("vague server is running")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	err := process.WaitUntilRunning(ctx, 100*time.Millisecond)
	if err == nil {
		t.Fatal("expected timeout error")
	}
}

func TestServerLogPath(t *testing.T) {
	t.Parallel()

	if os.Getenv("XDG_RUNTIME_DIR") == "" {
		t.Skip("XDG_RUNTIME_DIR is not set")
	}

	path, err := process.ServerLogPath()
	if err != nil {
		t.Fatal(err)
	}

	if filepath.Base(path) != "vague-server.log" {
		t.Fatalf("got %q", path)
	}
}
