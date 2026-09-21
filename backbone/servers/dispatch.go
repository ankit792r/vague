// Dispatch the incomming messages to handler
package backbone

import (
	"context"
	"fmt"
	"vague/backbone/process"
	"vague/backbone/session"
	"vague/bonefire/display"
	"vague/bonefire/editor"
	"vague/bonefire/workspace"
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
	fail := func(err error) (any, error) {
		if frameID != 0 && err != nil {
			s.pushEcho(ctx, sess, frameID, err.Error(), editor.EchoError)
		}
		return nil, err
	}

	switch params.Name {
	case "edit", "e":
		if frameID == 0 {
			return fail(fmt.Errorf("session is not attached to a frame"))
		}
		if len(params.Args) == 0 {
			return fail(fmt.Errorf("edit: file name required"))
		}

		path := params.Args[0]
		result, err := s.runtime.Do(ctx, func(ws *workspace.Workspace) (any, error) {
			buf, err := ws.OpenFile(frameID, path, params.Bang)
			if err != nil {
				return nil, editor.OpenFileError(path, err)
			}
			if err := ws.SetEcho(frameID, fmt.Sprintf(`"%s"`, buf.Name), editor.EchoInfo); err != nil {
				return nil, err
			}
			return editor.BufferInfo(buf, true), nil
		})
		if err != nil {
			return fail(err)
		}

		s.pushRedraw(ctx, sess, frameID)
		return result, nil

	case "write", "w":
		if frameID == 0 {
			return fail(fmt.Errorf("session is not attached to a frame"))
		}

		path := ""
		if len(params.Args) > 0 {
			path = params.Args[0]
		}

		result, err := s.runtime.Do(ctx, func(ws *workspace.Workspace) (any, error) {
			_, _, buf, err := ws.FrameContext(frameID)
			if err != nil {
				return nil, err
			}

			if err := ws.WriteFile(frameID, path, params.Bang); err != nil {
				return nil, editor.WriteFileError(path, err)
			}

			if err := ws.SetEcho(frameID, editor.WriteEchoMessage(buf), editor.EchoInfo); err != nil {
				return nil, err
			}

			return editor.BufferInfo(buf, true), nil
		})
		if err != nil {
			return fail(err)
		}

		s.pushRedraw(ctx, sess, frameID)
		return result, nil

	case "quit", "q":
		if frameID == 0 {
			return fail(fmt.Errorf("session is not attached to a frame"))
		}

		_, err := s.runtime.Do(ctx, func(ws *workspace.Workspace) (any, error) {
			return nil, ws.QuitFrame(frameID, params.Bang)
		})
		if err != nil {
			return fail(err)
		}

		sess.SetFrameID(0)
		sess.Notify(process.MethodQuit, process.QuitParams{FrameID: frameID})

		return map[string]any{"quit": true}, nil

	case "wrap", "nowrap":
		if frameID == 0 {
			return fail(fmt.Errorf("session is not attached to a frame"))
		}

		wrap := params.Name == "wrap"
		msg := "wrap off"
		if wrap {
			msg = "wrap on"
		}

		result, err := s.runtime.Do(ctx, func(ws *workspace.Workspace) (any, error) {
			if err := ws.SetWindowWrap(frameID, wrap); err != nil {
				return nil, err
			}
			if err := ws.SetEcho(frameID, msg, editor.EchoInfo); err != nil {
				return nil, err
			}
			return map[string]any{"wrap": wrap}, nil
		})
		if err != nil {
			return fail(err)
		}

		s.pushRedraw(ctx, sess, frameID)
		return result, nil
	}

	return fail(fmt.Errorf("Unknown command: %s", params.Name))
}

func (s *Server) handleInput(ctx context.Context, sess *session.Session, params process.InputParams) {
	frameID := sess.FrameID()
	if frameID == 0 {
		return
	}

	_, err := s.runtime.Do(ctx, func(ws *workspace.Workspace) (any, error) {
		return nil, ws.HandleInput(frameID, params.Keys)
	})
	if err != nil {
		s.pushEcho(ctx, sess, frameID, err.Error(), editor.EchoError)
		return
	}

	s.pushRedraw(ctx, sess, frameID)
}

func (s *Server) handleFrameAttach(ctx context.Context, sess *session.Session, params process.AttachParams) (process.AttachResult, error) {
	select {
	case <-ctx.Done():
		return process.AttachResult{}, ctx.Err()
	default:
	}

	result, err := s.runtime.Do(ctx, func(ws *workspace.Workspace) (any, error) {
		var frame *display.Frame
		var err error

		if len(params.Files) > 0 && params.Files[0] != "" {
			frame, err = ws.NewFrame(0, 0, params.WorkDir, params.Files[0])
		} else {
			frame, err = ws.NewFrame(0, 0, params.WorkDir)
		}
		if err != nil {
			return nil, err
		}

		ws.InvalidateFrame(frame.ID)

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
		result, err := s.runtime.Do(ctx, func(ws *workspace.Workspace) (any, error) {
			frame, err := ws.NewFrame(params.Height, params.Widht, "")
			if err != nil {
				return nil, err
			}

			ws.InvalidateFrame(frame.ID)

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

	_, err := s.runtime.Do(ctx, func(ws *workspace.Workspace) (any, error) {
		if err := ws.ResizeFrame(frameID, params.Widht, params.Height); err != nil {
			return nil, err
		}
		ws.InvalidateFrame(frameID)
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
	result, err := s.runtime.Do(ctx, func(ws *workspace.Workspace) (any, error) {
		redraw, ok := ws.RenderRedraw(frameID)
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

func (s *Server) pushEcho(ctx context.Context, sess *session.Session, frameID uint64, message, kind string) {
	_, err := s.runtime.Do(ctx, func(ws *workspace.Workspace) (any, error) {
		return nil, ws.SetEcho(frameID, message, kind)
	})
	if err != nil {
		return
	}

	s.pushRedraw(ctx, sess, frameID)
}
