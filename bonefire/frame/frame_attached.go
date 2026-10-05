package frame

import (
	"context"
	"fmt"
	"vague/backbone/process/protocols"
	"vague/backbone/runtime"
	"vague/backbone/session"
)

func FrameAttached(ctx context.Context, sess *session.Session, params *protocols.CommandParams, rt *runtime.Runtime, messID uint64) {
	// res, err := rt.Do(ctx, func(editor, frames, windows string) (any, error) {
	// 	return protocols.FrameAttachResult{}, nil
	// })

	// if err != nil {
	// 	sess.Reply(messID, nil, err)
	// 	return
	// }

	fmt.Println("Reached in frame attached in server")

	sess.Reply(messID, sess.State, nil)
}
