// Dispatch the incomming messages to handler
package backbone

import (
	"context"
	"fmt"
	"vague/backbone/process"
	"vague/backbone/session"
	"vague/bonefire/editor"
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

	case process.MethodFrameAttach:
		var params process.AttachParams
		if err := msg.DecodeParams(&params); err != nil {
			sess.Reply(msg.ID, nil, err)
			return
		}
		result, err := s.handleFrameAttach(ctx, sess, params)
		sess.Reply(msg.ID, result, err)

	case process.MethodFrameDetach:
		err := s.handleFrameDetach(ctx, sess)
		sess.Reply(msg.ID, nil, err)

	case process.MethodFrameReady:
		var params process.FrameReadyParams
		if err := msg.DecodeParams(&params); err != nil {
			sess.Reply(msg.ID, nil, err)
			return
		}
		result, err := s.handleFrameReady(ctx, sess, params)
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

	frameID := sess.FrameID()
	switch params.Name {
	case "wrap", "nowrap":
		if frameID == 0 {
			return nil, fmt.Errorf("session is not attached to a frame")
		}

		wrap := params.Name == "wrap"
		result, err := s.runtime.Do(ctx, func(ed *editor.Editor) (any, error) {
			if err := ed.SetWindowWrap(frameID, wrap); err != nil {
				return nil, err
			}
			return map[string]any{"wrap": wrap}, nil
		})
		if err != nil {
			return nil, err
		}

		s.pushRedraw(ctx, sess, frameID)
		return result, nil
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

func (s *Server) handleFrameAttach(ctx context.Context, sess *session.Session, params process.AttachParams) (process.AttachResult, error) {
	select {
	case <-ctx.Done():
		return process.AttachResult{}, ctx.Err()
	default:
	}

	_ = params

	result, err := s.runtime.Do(ctx, func(ed *editor.Editor) (any, error) {
		frame, err := ed.NewFrame(0, 0)
		if err != nil {
			return nil, err
		}

		ed.InvalidateFrame(frame.ID)

		return process.AttachResult{
			SessionID: sess.Id,
			FrameID:   frame.ID,
		}, nil
	})
	if err != nil {
		return process.AttachResult{}, err
	}

	attached := result.(process.AttachResult)
	sess.SetFrameID(attached.FrameID)

	return attached, nil
}

func (s *Server) handleFrameDetach(ctx context.Context, sess *session.Session) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	sess.SetFrameID(0)
	return nil
}

func (s *Server) handleFrameReady(ctx context.Context, sess *session.Session, params process.FrameReadyParams) (process.FrameReadyResult, error) {
	select {
	case <-ctx.Done():
		return process.FrameReadyResult{}, ctx.Err()
	default:
	}

	frameID := sess.FrameID()
	if frameID == 0 {
		result, err := s.runtime.Do(ctx, func(ed *editor.Editor) (any, error) {
			frame, err := ed.NewFrame(params.Height, params.Widht)
			if err != nil {
				return nil, err
			}

			ed.InvalidateFrame(frame.ID)

			return process.FrameReadyResult{
				SessionID: sess.Id,
				FrameID:   frame.ID,
			}, nil
		})
		if err != nil {
			return process.FrameReadyResult{}, err
		}

		ready := result.(process.FrameReadyResult)
		sess.SetFrameID(ready.FrameID)
		s.pushRedraw(ctx, sess, ready.FrameID)

		return ready, nil
	}

	_, err := s.runtime.Do(ctx, func(ed *editor.Editor) (any, error) {
		if err := ed.ResizeFrame(frameID, params.Widht, params.Height); err != nil {
			return nil, err
		}
		ed.InvalidateFrame(frameID)
		return nil, nil
	})
	if err != nil {
		return process.FrameReadyResult{}, err
	}

	s.pushRedraw(ctx, sess, frameID)

	return process.FrameReadyResult{
		SessionID: sess.Id,
		FrameID:   frameID,
	}, nil
}

func (s *Server) pushRedraw(ctx context.Context, sess *session.Session, frameID uint64) {
	result, err := s.runtime.Do(ctx, func(ed *editor.Editor) (any, error) {
		redraw, ok := ed.RenderRedraw(frameID)
		if !ok {
			return nil, fmt.Errorf("frame %d has no pending redraw", frameID)
		}

		return redraw, nil
	})
	if err != nil {
		return
	}

	sess.Notify(process.MethodRedraw, result.(process.Redraw))
}
