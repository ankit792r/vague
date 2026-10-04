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
	if err := msg.DecodeParams(&params); err != nil {
		sess.Reply(msg.ID, nil, err)
		return
	}

	// sess.Reply(msg.ID, "{}", nil)

	switch params.Name {
	case protocols.FrameAttached:
		res, err := frame.UiAttached(ctx, sess, &params, s.runtime)
		sess.Reply(msg.ID, res, err)
	}
}
