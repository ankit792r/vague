// Dispatch the incomming messages to handler
package backbone

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"vague/backbone/process"
	"vague/backbone/session"
	"vague/bonefire/editor"
	frame "vague/bonefire/frame"
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

	case process.MethodUiAttach:
		var params process.UiAttachParams
		if err := msg.DecodeParams(&params); err != nil {
			sess.Reply(msg.ID, nil, err)
			return
		}
		result, err := s.handleUiAttach(ctx, sess, params)
		sess.Reply(msg.ID, result, err)

	case process.MethodUiDetach:
		err := s.handleUiDetach(ctx, sess)
		sess.Reply(msg.ID, nil, err)

	case process.MethodUiReady:
		var params process.UiReadyParams
		if err := msg.DecodeParams(&params); err != nil {
			sess.Reply(msg.ID, nil, err)
			return
		}
		result, err := s.handleUiReady(ctx, sess, params)
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
			result, err := s.runtime.Do(ctx, func(ws *workspace.Workspace) (any, error) {
				buf, err := ws.ReloadCurrentBuffer(frameID, params.Bang)
				if err != nil {
					return nil, err
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
		}

		path := strings.Join(params.Args, " ")
		result, err := s.runtime.Do(ctx, func(ws *workspace.Workspace) (any, error) {
			buf, err := ws.EditFile(frameID, path)
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

		path := argPath(params.Args)
		result, err := s.executeWrite(ctx, sess, frameID, path, params.Bang)
		if err != nil {
			return fail(err)
		}
		return result, nil

	case "wq":
		if frameID == 0 {
			return fail(fmt.Errorf("session is not attached to a frame"))
		}

		path := argPath(params.Args)
		if _, err := s.executeWrite(ctx, sess, frameID, path, params.Bang); err != nil {
			return fail(err)
		}
		if err := s.executeQuit(ctx, sess, frameID, false); err != nil {
			return fail(err)
		}
		return map[string]any{"quit": true}, nil

	case "x":
		if frameID == 0 {
			return fail(fmt.Errorf("session is not attached to a frame"))
		}

		path := argPath(params.Args)
		modified, err := s.bufferModified(ctx, frameID)
		if err != nil {
			return fail(err)
		}
		if modified {
			if _, err := s.executeWrite(ctx, sess, frameID, path, params.Bang); err != nil {
				return fail(err)
			}
		}
		if err := s.executeQuit(ctx, sess, frameID, false); err != nil {
			return fail(err)
		}
		return map[string]any{"quit": true}, nil

	case "quit", "q":
		if frameID == 0 {
			return fail(fmt.Errorf("session is not attached to a frame"))
		}

		if err := s.executeQuit(ctx, sess, frameID, params.Bang); err != nil {
			return fail(err)
		}

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

	case "number", "nonumber":
		if frameID == 0 {
			return fail(fmt.Errorf("session is not attached to a frame"))
		}

		number := params.Name == "number"
		msg := "number off"
		if number {
			msg = "number on"
		}

		result, err := s.runtime.Do(ctx, func(ws *workspace.Workspace) (any, error) {
			if err := ws.SetWindowNumber(frameID, number); err != nil {
				return nil, err
			}
			if err := ws.SetEcho(frameID, msg, editor.EchoInfo); err != nil {
				return nil, err
			}
			return map[string]any{"number": number}, nil
		})
		if err != nil {
			return fail(err)
		}

		s.pushRedraw(ctx, sess, frameID)
		return result, nil

	case "set":
		if frameID == 0 {
			return fail(fmt.Errorf("session is not attached to a frame"))
		}
		if len(params.Args) == 0 {
			return fail(fmt.Errorf("set: option required"))
		}

		opt := strings.ToLower(strings.Join(params.Args, " "))
		switch opt {
		case "number", "nu":
			params.Name = "number"
		case "nonumber", "nonu":
			params.Name = "nonumber"
		case "wrap":
			params.Name = "wrap"
		case "nowrap":
			params.Name = "nowrap"
		case "ignorecase", "ic", "noignorecase", "noic",
			"smartcase", "scs", "nosmartcase", "noscs",
			"hlsearch", "hls", "nohlsearch", "nohls",
			"incsearch", "noincsearch",
			"wrapscan", "ws", "nowrapscan", "nows":
			_, err := s.runtime.Do(ctx, func(ws *workspace.Workspace) (any, error) {
				msg, err := ws.RunExLine(frameID, "set "+opt)
				if err != nil {
					return nil, err
				}
				if msg != "" {
					_ = ws.SetEcho(frameID, msg, editor.EchoInfo)
				}
				return nil, nil
			})
			if err != nil {
				return fail(err)
			}
			s.pushRedraw(ctx, sess, frameID)
			return nil, nil
		default:
			return fail(fmt.Errorf("Unknown option: %s", opt))
		}
		return s.handleExecute(ctx, sess, params)

	case "goto", "go":
		if frameID == 0 {
			return fail(fmt.Errorf("session is not attached to a frame"))
		}

		line := params.Count
		if line <= 0 && len(params.Args) > 0 {
			var err error
			line, err = parsePositiveInt(params.Args[0])
			if err != nil {
				return fail(editor.ErrInvalidLine)
			}
		}
		if line <= 0 {
			return fail(editor.ErrInvalidLine)
		}

		_, err := s.runtime.Do(ctx, func(ws *workspace.Workspace) (any, error) {
			if err := ws.GoToLine(frameID, line); err != nil {
				return nil, err
			}
			if err := ws.SetEcho(frameID, fmt.Sprintf("line %d", line), editor.EchoInfo); err != nil {
				return nil, err
			}
			return nil, nil
		})
		if err != nil {
			return fail(err)
		}

		s.pushRedraw(ctx, sess, frameID)
		return nil, nil

	case "only":
		if frameID == 0 {
			return fail(fmt.Errorf("session is not attached to a frame"))
		}

		_, err := s.runtime.Do(ctx, func(ws *workspace.Workspace) (any, error) {
			if err := ws.OnlyWindow(frameID); err != nil {
				return nil, err
			}
			return nil, nil
		})
		if err != nil {
			return fail(err)
		}

		s.pushRedraw(ctx, sess, frameID)
		return nil, nil

	case "search":
		if frameID == 0 {
			return fail(fmt.Errorf("session is not attached to a frame"))
		}

		pattern := strings.Join(params.Args, " ")
		forward := !params.Bang

		_, err := s.runtime.Do(ctx, func(ws *workspace.Workspace) (any, error) {
			if err := ws.Search(frameID, pattern, forward); err != nil {
				return nil, err
			}
			return nil, nil
		})
		if err != nil {
			return fail(err)
		}

		s.pushRedraw(ctx, sess, frameID)
		return nil, nil

	case "search_begin":
		if frameID == 0 {
			return fail(fmt.Errorf("session is not attached to a frame"))
		}
		forward := !params.Bang
		_, err := s.runtime.Do(ctx, func(ws *workspace.Workspace) (any, error) {
			return nil, ws.BeginIncsearch(frameID, forward)
		})
		if err != nil {
			return fail(err)
		}
		s.pushRedraw(ctx, sess, frameID)
		return nil, nil

	case "search_preview":
		if frameID == 0 {
			return fail(fmt.Errorf("session is not attached to a frame"))
		}
		pattern := strings.Join(params.Args, " ")
		_, err := s.runtime.Do(ctx, func(ws *workspace.Workspace) (any, error) {
			return nil, ws.PreviewIncsearch(frameID, pattern)
		})
		if err != nil {
			return fail(err)
		}
		s.pushRedraw(ctx, sess, frameID)
		return nil, nil

	case "search_cancel":
		if frameID == 0 {
			return fail(fmt.Errorf("session is not attached to a frame"))
		}
		_, err := s.runtime.Do(ctx, func(ws *workspace.Workspace) (any, error) {
			return nil, ws.CancelIncsearch(frameID)
		})
		if err != nil {
			return fail(err)
		}
		s.pushRedraw(ctx, sess, frameID)
		return nil, nil

	case "complete":
		if frameID == 0 {
			return fail(fmt.Errorf("session is not attached to a frame"))
		}
		if len(params.Args) < 2 {
			return fail(fmt.Errorf("complete: kind and prefix required"))
		}
		kind := params.Args[0]
		prefix := strings.Join(params.Args[1:], " ")
		result, err := s.runtime.Do(ctx, func(ws *workspace.Workspace) (any, error) {
			candidates, err := ws.Complete(frameID, kind, prefix)
			if err != nil {
				return nil, err
			}
			return map[string]any{"candidates": candidates}, nil
		})
		if err != nil {
			return fail(err)
		}
		return result, nil

	case "register_get":
		if frameID == 0 {
			return fail(fmt.Errorf("session is not attached to a frame"))
		}
		regName := `"`
		if len(params.Args) > 0 {
			regName = params.Args[0]
		}
		result, err := s.runtime.Do(ctx, func(ws *workspace.Workspace) (any, error) {
			if _, _, _, err := ws.FrameContext(frameID); err != nil {
				return nil, err
			}
			text, ok := ws.Editor.GetRegister(regName)
			if !ok {
				return map[string]any{"text": ""}, nil
			}
			return map[string]any{"text": text}, nil
		})
		if err != nil {
			return fail(err)
		}
		return result, nil

	case "register_set":
		if frameID == 0 {
			return fail(fmt.Errorf("session is not attached to a frame"))
		}
		if len(params.Args) < 2 {
			return fail(fmt.Errorf("register_set: name and text required"))
		}
		name := params.Args[0]
		text := strings.Join(params.Args[1:], " ")
		_, err := s.runtime.Do(ctx, func(ws *workspace.Workspace) (any, error) {
			if _, _, _, err := ws.FrameContext(frameID); err != nil {
				return nil, err
			}
			if err := ws.Editor.SetRegister(name, text, false); err != nil {
				return nil, err
			}
			return nil, nil
		})
		if err != nil {
			return fail(err)
		}
		s.pushRedraw(ctx, sess, frameID)
		return map[string]any{"ok": true}, nil

	case "clear_echo":
		if frameID == 0 {
			return nil, nil
		}
		_, err := s.runtime.Do(ctx, func(ws *workspace.Workspace) (any, error) {
			ws.ClearEcho(frameID)
			return nil, nil
		})
		if err != nil {
			return fail(err)
		}
		s.pushRedraw(ctx, sess, frameID)
		return nil, nil

	case "help":
		if frameID == 0 {
			return fail(fmt.Errorf("session is not attached to a frame"))
		}
		_, err := s.runtime.Do(ctx, func(ws *workspace.Workspace) (any, error) {
			if err := ws.SetEcho(frameID, editor.HelpReference, editor.EchoInfo); err != nil {
				return nil, err
			}
			return nil, nil
		})
		if err != nil {
			return fail(err)
		}
		s.pushRedraw(ctx, sess, frameID)
		return nil, nil

	case "bnext", "bn":
		if frameID == 0 {
			return fail(fmt.Errorf("session is not attached to a frame"))
		}

		result, err := s.runtime.Do(ctx, func(ws *workspace.Workspace) (any, error) {
			buf, err := ws.SwitchToNextBuffer(frameID, params.Bang)
			if err != nil {
				return nil, err
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

	case "bprev", "bp":
		if frameID == 0 {
			return fail(fmt.Errorf("session is not attached to a frame"))
		}

		result, err := s.runtime.Do(ctx, func(ws *workspace.Workspace) (any, error) {
			buf, err := ws.SwitchToPrevBuffer(frameID, params.Bang)
			if err != nil {
				return nil, err
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

	case "buffer", "b":
		if frameID == 0 {
			return fail(fmt.Errorf("session is not attached to a frame"))
		}
		if len(params.Args) == 0 {
			return fail(fmt.Errorf("buffer: name or number required"))
		}

		spec := strings.Join(params.Args, " ")
		result, err := s.runtime.Do(ctx, func(ws *workspace.Workspace) (any, error) {
			buf, err := ws.SwitchToBuffer(frameID, spec, params.Bang)
			if err != nil {
				return nil, err
			}
			if err := ws.SetEcho(frameID, fmt.Sprintf(`"%s"`, buf.Name), editor.EchoInfo); err != nil {
				return nil, err
			}
			return editor.BufferInfo(buf, true), nil
		})
		if err != nil {
			if errors.Is(err, editor.ErrBufferNotFound) {
				return fail(fmt.Errorf("Buffer not found"))
			}
			return fail(err)
		}

		s.pushRedraw(ctx, sess, frameID)
		return result, nil

	case "buffers", "ls":
		if frameID == 0 {
			return fail(fmt.Errorf("session is not attached to a frame"))
		}

		_, err := s.runtime.Do(ctx, func(ws *workspace.Workspace) (any, error) {
			msg, err := ws.BufferListMessage(frameID)
			if err != nil {
				return nil, err
			}
			if err := ws.SetEcho(frameID, msg, editor.EchoInfo); err != nil {
				return nil, err
			}
			return nil, nil
		})
		if err != nil {
			return fail(err)
		}

		s.pushRedraw(ctx, sess, frameID)
		return nil, nil

	case "ex":
		if frameID == 0 {
			return fail(fmt.Errorf("session is not attached to a frame"))
		}
		line := strings.Join(params.Args, " ")
		return s.runExLine(ctx, sess, frameID, line)

	case "confirm_answer":
		if frameID == 0 {
			return fail(fmt.Errorf("session is not attached to a frame"))
		}
		if len(params.Args) == 0 {
			return fail(fmt.Errorf("confirm_answer: y/n/a/q/l required"))
		}
		_, err := s.runtime.Do(ctx, func(ws *workspace.Workspace) (any, error) {
			return nil, ws.AnswerConfirm(frameID, params.Args[0], false)
		})
		if err != nil {
			return fail(err)
		}
		s.pushRedraw(ctx, sess, frameID)
		return nil, nil

	case "wa", "wall":
		if frameID == 0 {
			return fail(fmt.Errorf("session is not attached to a frame"))
		}
		_, err := s.runtime.Do(ctx, func(ws *workspace.Workspace) (any, error) {
			return nil, ws.WriteAll(frameID, params.Bang)
		})
		if err != nil {
			return fail(err)
		}
		s.pushRedraw(ctx, sess, frameID)
		return nil, nil

	case "qa", "qall":
		if frameID == 0 {
			return fail(fmt.Errorf("session is not attached to a frame"))
		}
		_, err := s.runtime.Do(ctx, func(ws *workspace.Workspace) (any, error) {
			return nil, ws.QuitAll(params.Bang)
		})
		if err != nil {
			return fail(err)
		}
		return map[string]any{"quit": true}, nil

	case "wqa", "wqall":
		if frameID == 0 {
			return fail(fmt.Errorf("session is not attached to a frame"))
		}
		_, err := s.runtime.Do(ctx, func(ws *workspace.Workspace) (any, error) {
			return nil, ws.WQAll(frameID, params.Bang)
		})
		if err != nil {
			return fail(err)
		}
		return map[string]any{"quit": true}, nil

	case "xa", "xall":
		if frameID == 0 {
			return fail(fmt.Errorf("session is not attached to a frame"))
		}
		_, err := s.runtime.Do(ctx, func(ws *workspace.Workspace) (any, error) {
			return nil, ws.XAll(frameID, params.Bang)
		})
		if err != nil {
			return fail(err)
		}
		return map[string]any{"quit": true}, nil

	default:
		return fail(fmt.Errorf("Unknown command: %s", params.Name))
	}
}

func (s *Server) runExLine(ctx context.Context, sess *session.Session, frameID uint64, line string) (any, error) {
	fail := func(err error) (any, error) {
		if frameID != 0 && err != nil {
			s.pushEcho(ctx, sess, frameID, err.Error(), editor.EchoError)
		}
		return nil, err
	}

	line = strings.TrimSpace(line)
	if strings.HasPrefix(strings.ToLower(line), "find ") || strings.HasPrefix(strings.ToLower(line), "sf ") {
		parts := strings.Fields(line)
		if len(parts) >= 2 {
			pattern := parts[1]
			_, err := s.runtime.Do(ctx, func(ws *workspace.Workspace) (any, error) {
				return nil, ws.ExFindFile(frameID, pattern)
			})
			if err != nil {
				return fail(err)
			}
			s.pushRedraw(ctx, sess, frameID)
			return nil, nil
		}
	}

	echo, err := s.runtime.Do(ctx, func(ws *workspace.Workspace) (any, error) {
		msg, runErr := ws.RunExLine(frameID, line)
		if runErr != nil {
			if errors.Is(runErr, editor.ErrConfirmPending) {
				return msg, nil
			}
			return nil, runErr
		}
		if msg != "" {
			if setErr := ws.SetEcho(frameID, msg, editor.EchoInfo); setErr != nil {
				return nil, setErr
			}
		}
		return msg, nil
	})
	if err != nil {
		return fail(err)
	}
	s.pushRedraw(ctx, sess, frameID)
	return echo, nil
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

func (s *Server) handleUiAttach(ctx context.Context, sess *session.Session, params process.UiAttachParams) (process.UiAttachResult, error) {
	select {
	case <-ctx.Done():
		return process.UiAttachResult{}, ctx.Err()
	default:
	}

	result, err := s.runtime.Do(ctx, func(ws *workspace.Workspace) (any, error) {
		var frame *frame.Frame
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

		return process.UiAttachResult{
			SessionID: sess.Id,
			FrameID:   frame.ID,
		}, nil
	})
	if err != nil {
		return process.UiAttachResult{}, err
	}

	attached := result.(process.UiAttachResult)
	sess.SetFrameID(attached.FrameID)

	return attached, nil
}

func (s *Server) handleUiDetach(ctx context.Context, sess *session.Session) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	sess.SetFrameID(0)
	return nil
}

func (s *Server) handleUiReady(ctx context.Context, sess *session.Session, params process.UiReadyParams) (process.UiReadyResult, error) {
	select {
	case <-ctx.Done():
		return process.UiReadyResult{}, ctx.Err()
	default:
	}

	frameID := sess.FrameID()
	if frameID == 0 {
		result, err := s.runtime.Do(ctx, func(ws *workspace.Workspace) (any, error) {
			frame, err := ws.NewFrame(params.Height, params.Width, "")
			if err != nil {
				return nil, err
			}

			ws.InvalidateFrame(frame.ID)

			return process.UiReadyResult{
				SessionID: sess.Id,
				FrameID:   frame.ID,
			}, nil
		})
		if err != nil {
			return process.UiReadyResult{}, err
		}

		ready := result.(process.UiReadyResult)
		sess.SetFrameID(ready.FrameID)
		s.pushRedraw(ctx, sess, ready.FrameID)

		return ready, nil
	}

	_, err := s.runtime.Do(ctx, func(ws *workspace.Workspace) (any, error) {
		if err := ws.ResizeFrame(frameID, params.Width, params.Height); err != nil {
			return nil, err
		}
		ws.InvalidateFrame(frameID)
		return nil, nil
	})
	if err != nil {
		return process.UiReadyResult{}, err
	}

	s.pushRedraw(ctx, sess, frameID)

	return process.UiReadyResult{
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

func argPath(args []string) string {
	if len(args) > 0 {
		return args[0]
	}
	return ""
}

func (s *Server) bufferModified(ctx context.Context, frameID uint64) (bool, error) {
	result, err := s.runtime.Do(ctx, func(ws *workspace.Workspace) (any, error) {
		_, _, buf, err := ws.FrameContext(frameID)
		if err != nil {
			return nil, err
		}
		return buf.Modified(), nil
	})
	if err != nil {
		return false, err
	}
	return result.(bool), nil
}

func (s *Server) executeWrite(ctx context.Context, sess *session.Session, frameID uint64, path string, bang bool) (any, error) {
	result, err := s.runtime.Do(ctx, func(ws *workspace.Workspace) (any, error) {
		_, _, buf, err := ws.FrameContext(frameID)
		if err != nil {
			return nil, err
		}

		if err := ws.WriteFile(frameID, path, bang); err != nil {
			return nil, editor.WriteFileError(path, err)
		}

		if err := ws.SetEcho(frameID, editor.WriteEchoMessage(buf), editor.EchoInfo); err != nil {
			return nil, err
		}

		return editor.BufferInfo(buf, true), nil
	})
	if err != nil {
		return nil, err
	}

	s.pushRedraw(ctx, sess, frameID)
	return result, nil
}

func (s *Server) executeQuit(ctx context.Context, sess *session.Session, frameID uint64, bang bool) error {
	_, err := s.runtime.Do(ctx, func(ws *workspace.Workspace) (any, error) {
		return nil, ws.QuitFrame(frameID, bang)
	})
	if err != nil {
		return err
	}

	sess.SetFrameID(0)
	sess.Notify(process.MethodQuit, process.QuitParams{FrameID: frameID})
	return nil
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

func parsePositiveInt(s string) (int, error) {
	n, err := strconv.Atoi(strings.TrimSpace(s))
	if err != nil || n <= 0 {
		return 0, fmt.Errorf("invalid number")
	}
	return n, nil
}
