package frame

import (
	"encoding/json"
	"fmt"
	"vague/bonefire"
	"vague/bonefire/buffer"
	"vague/bonefire/window"
)

func init() {
	fmt.Println("Registering frame_attached handler")
	bonefire.RegisterV2("frame_attached", FrameAttached)
}

type FrameAttachedParams struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

func FrameAttached(dsp *bonefire.DispatcherV2, params *json.RawMessage) (any, error) {
	_, sess := dsp.Ctx, dsp.Sess
	fmt.Println("Params", params)
	var parsedParams FrameAttachedParams
	if err := json.Unmarshal(*params, &parsedParams); err != nil {
		fmt.Println("Error unmarshalling params", err)
		return nil, fmt.Errorf("frame_attached: %w", err)
	}

	fmt.Println("Params", parsedParams)

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
		return nil, err
	}

	return res, nil
}
