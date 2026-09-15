package frame

import (
	"context"
	backbone "vague/backbone/clients"
	window "vague/bonefire/window"
)

type Frame struct {
	Id     uint64
	ctx    context.Context
	client *backbone.Client

	RootWindow *window.Window

	Echo string
}

func NewFrame(ctx context.Context) error {
	conn, err := backbone.ClientConnect()
	if err != nil {
		return err
	}
	defer conn.Close()

	frame := &Frame{
		ctx:    ctx,
		client: conn,
	}

	// This to be called in last, since it will block the other execution
	return frame.BuildWebView()
}
