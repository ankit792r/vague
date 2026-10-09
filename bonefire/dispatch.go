package bonefire

import (
	"context"
	"encoding/json"
	"fmt"
	"vague/backbone/process"
	"vague/backbone/runtime"
	"vague/backbone/session"
)

type HandlerV2 func(dsp *DispatcherV2, params *json.RawMessage) (any, error)

var handlersV2 = make(map[string]HandlerV2)

type DispatcherV2 struct {
	Ctx  context.Context
	Sess *session.Session
	Rt   *runtime.Runtime
}

func RegisterV2(method string, h HandlerV2) {
	if _, dup := handlersV2[method]; dup {
		panic("duplicate handler: " + method)
	}
	handlersV2[method] = h
}

func (d *DispatcherV2) Dispatch(msg *process.Message) {
	fmt.Println("Dispatching message", msg.Method)
	handler, ok := handlersV2[msg.Method]
	if !ok {
		panic("handler not found: " + msg.Method)
	}
	fmt.Println("Params", msg.Params)
	res, err := handler(d, &msg.Params)
	if err != nil {
		d.Sess.Reply(msg.ID, nil, err)
		return
	}
	d.Sess.Reply(msg.ID, res, nil)
}

func (d *DispatcherV2) Execute(dfn runtime.DoFunc) (any, error) {
	fmt.Println("Executing function", dfn)
	res, err := d.Rt.Do(d.Ctx, dfn)
	return res, err
}
