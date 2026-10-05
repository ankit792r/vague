package server

import (
	"context"
	"fmt"
	"vague/backbone/process"
	"vague/backbone/process/protocols"
	"vague/backbone/session"
	"vague/bonefire/frame"
)

func (s *Server) dispatch(ctx context.Context, sess *session.Session, msg process.Message) {
	select {
	case <-ctx.Done():
		sess.Reply(msg.ID, nil, ctx.Err())
	default:
	}


		fmt.Println(msg)

	// frameID := sess.FrameID()

	// INFO: this function will send notification to ui on any failure
	// so instead of pushEcho we will return the Update and send Frame State
	// fail := func(err error) (any, error) {
	// 	if frameID != 0 && err != nil {
	// 		s.pushEcho(ctx, sess, frameID, err.Error(), editor.EchoError)
	// 	}
	// 	return nil, err
	// }

	var params protocols.CommandParams
	// FIXME: here is the issue, we dont getting the method name parsed correclty by DecodeParams method
	// instead we have to find and alternative way to pass Message
	if err := msg.DecodeParams(&params); err != nil {
		sess.Reply(msg.ID, nil, err)
		return
	}

	fmt.Println("params name", params)
	fmt.Println("protocol name", protocols.FrameAttached)

	switch params.Name {
	case protocols.FrameAttached:
		frame.FrameAttached(ctx, sess, &params, s.runtime, msg.ID)
	}
}
