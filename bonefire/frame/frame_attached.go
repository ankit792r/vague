package frame

import (
	"vague/backbone/process"
	"vague/bonefire"
	"vague/bonefire/buffer"
	"vague/bonefire/window"
)

func init() {
	bonefire.RegisterV2("frame_attached", FrameAttached)
}

func FrameAttached(dsp *bonefire.DispatcherV2, msg *process.Message) {
	_, sess := dsp.Ctx, dsp.Sess
	res, err := dsp.Execute(func(editor, frames, windows string) (any, error) {
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
		sess.Reply(msg.ID, nil, err)
		return
	}

	sess.Reply(msg.ID, res, nil)
}
