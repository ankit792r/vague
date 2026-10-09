package frame

import (
	"vague/backbone/dispatch"
	"vague/backbone/process"
	"vague/bonefire/buffer"
	"vague/bonefire/window"
)

func FrameAttached(dsp *dispatch.Dispatcher, msg *process.Message) {
	_, sess := dsp.Ctx, dsp.Sess

	res, err := dsp.Execute(func(editor, frames, windows string) (any, error) {
		state := dsp.Sess.NewSessionState()
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
		sess.Reply(msg.ID, nil, err)
		return
	}

	sess.Reply(msg.ID, res, nil)
}
