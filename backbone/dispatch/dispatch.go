package dispatch

import (
	"context"
	"vague/backbone/process"
	"vague/backbone/runtime"
	"vague/backbone/session"
	"vague/bonefire/frame"
)

type Dispatcher struct {
	Ctx  context.Context
	Sess *session.Session
	rt   *runtime.Runtime
}

func NewDispatcher(
	ctx context.Context,
	sess *session.Session,
	rt *runtime.Runtime,
) *Dispatcher {
	return &Dispatcher{
		Ctx:  ctx,
		Sess: sess,
		rt:   rt,
	}
}

func (d *Dispatcher) Dispatch(msg *process.Message) {
	select {
	case <-d.Ctx.Done():
		d.Sess.Reply(msg.ID, nil, d.Ctx.Err())
	default:
	}

	switch msg.Method {
	case FRAME_ATTACHED:
		frame.FrameAttached(d, msg)
	}
}

func (d *Dispatcher) Execute(dfn runtime.DoFunc) (any, error) {
	res, err := d.rt.Do(d.Ctx, dfn)
	return res, err
}
