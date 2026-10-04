package frame

import (
	"context"
	"fmt"
	"vague/backbone/process/protocols"
	"vague/backbone/runtime"
	"vague/backbone/session"
)

func UiAttached(ctx context.Context, sess *session.Session, params *protocols.CommandParams, rt *runtime.Runtime) (any, error) {
	fmt.Println("got request in UiAttached")
	res, err := rt.Do(ctx, func(editor, frames, windows string) (any, error) {
		return protocols.FrameAttachResult{}, nil
	})

	if err != nil {
		return protocols.FrameAttachResult{}, err
	}

	attached := res.(protocols.FrameAttachResult)
	return attached, nil
}
