package client

func RegisterFrameHandlers(b *Bridge) {
	b.Register("frame_attached", b.frameAttach)
	b.Register("frame_detached", b.frameDetach)
}

func (b *Bridge) frameAttach(params MethodParams) (any, error) {
	return nil, nil
}

func (b *Bridge) frameDetach(params MethodParams) (any, error) {
	return nil, nil
}


func RegisterInputHandlers(b *Bridge) {
}


func RegisterCommandHandlers(b *Bridge) {
}
