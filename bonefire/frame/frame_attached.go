package frame

import (
	"context"
	"fmt"
	"vague/backbone/process/protocols"
	"vague/backbone/runtime"
	"vague/backbone/session"
	"vague/bonefire/buffer"
	"vague/bonefire/window"
)

const MethodFrameAttached = "frame_attached"

func FrameAttached(ctx context.Context, sess *session.Session, params *protocols.CommandParams, rt *runtime.Runtime, messID uint64) {
	fmt.Println("FrameAttached", params)
	res, err := rt.Do(ctx, func(editor, frames, windows string) (any, error) {
		state := sess.NewSessionState()
		defer func() {
			sess.NextBufferId += 1
			sess.NextWindowId += 1
			sess.State = state
		}()

		root := window.NewWindow(sess.NextWindowId, sess.Id)
		buff := buffer.NewScratch(sess.NextBufferId, "Scratch")

		state.Root = root
		state.BufferMap[root.Id] = buff

		// Now here we have to calculate the lines for window
		fullView := window.LayoutView(buff.Text, 100, 100, false)

		root.Lines = fullView.Lines

		return state, nil
	})

	if err != nil {
		sess.Reply(messID, nil, err)
		return
	}

	sess.Reply(messID, res, nil)
}
