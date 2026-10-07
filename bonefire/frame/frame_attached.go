package frame

import (
	"context"
	"vague/backbone/process/protocols"
	"vague/backbone/runtime"
	"vague/backbone/session"
	"vague/bonefire/window"
)

const MethodFrameAttached = "frame_attached"

func FrameAttached(ctx context.Context, sess *session.Session, params *protocols.CommandParams, rt *runtime.Runtime, messID uint64) {
	res, err := rt.Do(ctx, func(editor, frames, windows string) (any, error) {
		state := sess.NewSessionState()

		state.NextWindowId += 1
		state.Root = window.NewWindow(state.NextWindowId, sess.Id)

		sess.State = state
		return state, nil
	})

	if err != nil {
		sess.Reply(messID, nil, err)
		return
	}

	sess.Reply(messID, res, nil)
}
