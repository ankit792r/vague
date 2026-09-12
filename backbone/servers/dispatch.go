// Dispatch the incomming messages to handler
package backbone

import (
	"context"
	"fmt"
	"vague/backbone/process"
	"vague/backbone/session"
)

func (s *Server) dispatch(ctx context.Context, sess *session.Session, msg process.Message) {
	switch msg.Kind {
	case process.KindRequest:
		s.dispatchRequest(ctx, sess, msg)
	case process.KindNotify:
		s.dispatchNotify(ctx, sess, msg)
	default:
		// Validate() already rejects unknown kinds; this is defensive.
		if msg.ID != 0 {
			sess.Reply(msg.ID, nil, fmt.Errorf("unexpected message kind %q", msg.Kind))
		}
	}
}

func (s *Server) dispatchRequest(ctx context.Context, sess *session.Session, msg process.Message) {
	switch msg.Method {
	case process.MethodExecute:
		var params process.ExecuteParams
		if err := msg.DecodeParams(&params); err != nil {
			sess.Reply(msg.ID, nil, err)
			return
		}
		result, err := s.handleExecute(ctx, sess, params)
		sess.Reply(msg.ID, result, err)
	default:
		sess.Reply(msg.ID, nil, fmt.Errorf("unknown method %q", msg.Method))
	}
}

func (s *Server) dispatchNotify(ctx context.Context, sess *session.Session, msg process.Message) {
	switch msg.Method {
	case process.MethodInput:
		var params process.InputParams
		if err := msg.DecodeParams(&params); err != nil {
			// Notifications have no reply channel; log or drop.
			return
		}
		s.handleInput(ctx, sess, params)
	default:
		// Unknown notification — ignore for now.
	}
}

// Stub until bonefire/editor exists.
func (s *Server) handleExecute(ctx context.Context, sess *session.Session, params process.ExecuteParams) (any, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	// Temporary echo so client-connect can prove the round-trip works.
	return map[string]any{
		"ok":    true,
		"name":  params.Name,
		"args":  params.Args,
		"bang":  params.Bang,
		"count": params.Count,
	}, nil
}

func (s *Server) handleInput(ctx context.Context, sess *session.Session, params process.InputParams) {
	// TODO: feed keys into the editor model.
	_ = ctx
	_ = sess
	_ = params
}
