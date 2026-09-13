package platform

import (
	"context"
	backbone "vague/backbone/clients"
)

type Frame struct {
	ctx    context.Context
	client *backbone.Client
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
