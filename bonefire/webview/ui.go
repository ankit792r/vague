package webview

import (
	"context"
	"log/slog"
	"os"

	backbone "vague/backbone/clients"
	"vague/backbone/process"
)

// UI is the desktop client: IPC connection plus webview window.
type UI struct {
	FrameID uint64
	ctx     context.Context
	client  *backbone.Client
}

// Run connects to the server, attaches a UI surface, and blocks in the webview.
func Run(ctx context.Context, files ...string) error {
	conn, err := backbone.ClientConnect()
	if err != nil {
		return err
	}

	defer func() {
		if err := conn.UiDetach(ctx); err != nil {
			slog.Error("ui detach", "error", err)
		}
		conn.Close()
	}()

	workDir, err := clientWorkDir()
	if err != nil {
		return err
	}

	attachResult, err := conn.UiAttach(ctx, process.UiAttachParams{
		Files:   files,
		WorkDir: workDir,
	})
	if err != nil {
		return err
	}

	ui := &UI{
		FrameID: attachResult.FrameID,
		ctx:     ctx,
		client:  conn,
	}

	return ui.open()
}

func clientWorkDir() (string, error) {
	workDir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	if workDir == "" || workDir == "/" {
		if home, err := os.UserHomeDir(); err == nil {
			return home, nil
		}
	}

	return workDir, nil
}
