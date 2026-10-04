package client

import "vague/backbone/process/protocols"

func RegisterFrameHandlers(b *Bridge) {
	b.Register("frame_attached", b.frameAttach)
	b.Register("frame_detached", b.frameDetach)
}

func (b *Bridge) frameAttach(params MethodParams) (any, error) {
	var result protocols.FrameAttachResult
	if err := b.cli.call(b.ctx, protocols.FrameAttached, params, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (b *Bridge) frameDetach(params MethodParams) (any, error) {
	var result protocols.FrameDetachResult
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
