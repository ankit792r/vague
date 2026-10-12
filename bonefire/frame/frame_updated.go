package frame

import (
	"encoding/json"
	"fmt"
	"vague/backbone/process"
	"vague/bonefire"
)

func init() {
	fmt.Println("Registering frame_updated handler")
	bonefire.RegisterV2("frame_updated", FrameUpdated)
}

type FrameUpdatedParams struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

func FrameUpdated(dsp *bonefire.DispatcherV2, params json.RawMessage) (any, error) {

	var parsedParams FrameUpdatedParams
	if err := process.DecodeInvokeArgs(params, &parsedParams); err != nil {
		fmt.Println("Error unmarshalling params", err)
		return nil, fmt.Errorf("frame_updated: %w", err)
	}

	res, err := dsp.Execute(func(editor, frames string) (any, error) {
		return nil, nil
	})

	if err != nil {
		return nil, err
	}

	return res, nil
}
