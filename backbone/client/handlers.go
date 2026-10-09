package client

import (
	"vague/backbone/dispatch"
	"vague/backbone/process/protocols"
	"vague/backbone/session"
)

func RegisterFrameHandlers(b *Bridge) {
	b.Register("frame_attached", b.frameAttach)
	b.Register("frame_detached", b.frameDetach)
}

func (b *Bridge) frameAttach(params MethodParams) (any, error) {
	var result session.SessionState
	if err := b.cli.call(b.ctx, dispatch.FRAME_ATTACHED, params, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (b *Bridge) frameDetach(params MethodParams) (any, error) {
	var result session.SessionState
	if err := b.cli.call(b.ctx, protocols.FrameDetached, params, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func RegisterInputHandlers(b *Bridge) {
	// TODO: register input handlers
}

func RegisterCommandHandlers(b *Bridge) {
	// TODO: register command handlers
}
