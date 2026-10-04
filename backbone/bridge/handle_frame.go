package bridge

import (
	"fmt"
)

func RegisterFrameHandler(b *Bridge) {
	b.Register("frame_attached", handlerFrameAttached)
	b.Register("frame_detached", handlerFrameDetach)
}

func handlerFrameAttached(params MethodParams) (any, error) {
	fmt.Println("got frame attach request")
	return nil, nil
}

func handlerFrameDetach(params MethodParams) (any, error) {
	fmt.Println("got frame dettach request")
	return nil, nil
}
