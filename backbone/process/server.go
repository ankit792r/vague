package process

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"time"
)

const serverStartupTimeout = 10 * time.Second

// Ping reports whether the vague server accepts connections.
func Ping() error {
	conn, err := Dial()
	if err != nil {
		return err
	}

	return conn.Close()
}

// EnsureServer starts a detached server when none is running and waits until
// it accepts connections.
func EnsureServer(ctx context.Context) error {
	err := Ping()
	if err == nil {
		return nil
	}

	if !errors.Is(err, ErrNotRunning) {
		return err
	}

	if err := StartDetachedServer(); err != nil {
		return err
	}

	return WaitUntilRunning(ctx, serverStartupTimeout)
}

// StartDetachedServer launches `vague server-start` in the background.
func StartDetachedServer() error {
	exe, err := os.Executable()
	if err != nil {
		return fmt.Errorf("resolve executable: %w", err)
	}

	logPath, err := ServerLogPath()
	if err != nil {
		return err
	}

	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("open server log: %w", err)
	}

	cmd := exec.Command(exe, "server-start")
	cmd.Stdout = logFile
	cmd.Stderr = logFile
	cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
	if home, err := os.UserHomeDir(); err == nil {
		cmd.Dir = home
	}

	if err := cmd.Start(); err != nil {
		logFile.Close()
		return fmt.Errorf("start server: %w", err)
	}

	if err := logFile.Close(); err != nil {
		return fmt.Errorf("close server log: %w", err)
	}

	return nil
}

// WaitUntilRunning polls until Ping succeeds or the timeout elapses.
func WaitUntilRunning(ctx context.Context, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)

	for {
		if err := Ping(); err == nil {
			return nil
		}

		if time.Now().After(deadline) {
			logPath, _ := ServerLogPath()
			if logPath != "" {
				return fmt.Errorf("timed out waiting for vague server (see %s)", logPath)
			}

			return fmt.Errorf("timed out waiting for vague server")
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(50 * time.Millisecond):
		}
	}
}

// ServerLogPath is where a detached server writes its output.
func ServerLogPath() (string, error) {
	runtimeDir := os.Getenv("XDG_RUNTIME_DIR")
	if runtimeDir == "" {
		return "", fmt.Errorf("XDG_RUNTIME_DIR is not set")
	}

	return filepath.Join(runtimeDir, "vague-server.log"), nil
}
