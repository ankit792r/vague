package frame

import (
	"context"
	"vague/backbone/process/protocols"
	"vague/backbone/runtime"
	"vague/backbone/session"
)

const MethodFrameAttached = "frame_attached"

// This will modify the session.State
func FrameAttached(ctx context.Context, sess *session.Session, params *protocols.CommandParams, rt *runtime.Runtime, messID uint64) {
	res, err := rt.Do(ctx, func(editor, frames, windows string) (any, error) {
		sess.State.FrameID = 10
		return sess.State, nil
	})

	if err != nil {
		sess.Reply(messID, nil, err)
		return
	}

	sess.Reply(messID, res, nil)
}
