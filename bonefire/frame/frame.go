package frame

import (
	"context"
	"log/slog"
	"os"
	backbone "vague/backbone/clients"
	"vague/backbone/process"
	"vague/bonefire/window"
)

type Frame struct {
	Id     uint64
	ctx    context.Context
	client *backbone.Client

	RootWindow *window.Node

	Echo string

	NextWindowId uint64
	NextBufferId uint64
}

// TODO: take frame Id from FrameAttach event which is provided by server
func NewFrame(ctx context.Context, files ...string) error {
	conn, err := backbone.ClientConnect()
	if err != nil {
		return err
	}

	defer func() {
		if err := conn.FrameDetach(ctx); err != nil {
			slog.Error("Error Detach", "error", err)
		}
		conn.Close()
	}()

	workDir, err := clientWorkDir()
	if err != nil {
		return err
	}

	attachResult, err := conn.FrameAttach(ctx, process.AttachParams{
		Files:   files,
		WorkDir: workDir,
	})
	if err != nil {
		return err
	}

	frame := &Frame{
		Id:     attachResult.FrameID,
		ctx:    ctx,
		client: conn,

		NextWindowId: 1,
		NextBufferId: 1,
	}

	w := window.NewWindow(frame.NextWindowId, frame.Id)
	frame.RootWindow = w.LeafNode()

	// TODO: layout the root window

	// This to be called in last, since it will block the other execution
	return frame.BuildWebView()
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
