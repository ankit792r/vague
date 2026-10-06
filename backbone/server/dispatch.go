package server

import (
	"context"
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

	var params protocols.CommandParams
	if err := msg.DecodeParams(&params); err != nil {
		sess.Reply(msg.ID, nil, err)
		return
	}

	switch msg.Method {
	case frame.MethodFrameAttached:
		frame.FrameAttached(ctx, sess, &params, s.runtime, msg.ID)
	}
}
