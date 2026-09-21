package webview

import (
	"encoding/json"
	"fmt"
	"vague/backbone/process"
)

type HostReply struct {
	ID     uint64          `json:"id"`
	Result json.RawMessage `json:"result,omitempty"`
	Error  string          `json:"error,omitempty"`
}

func (u *UI) HostRequest(id uint64, method string, params json.RawMessage) HostReply {
	reply := HostReply{ID: id}

	switch method {
	case process.MethodInput:
		var p process.InputParams
		if err := json.Unmarshal(params, &p); err != nil {
			reply.Error = err.Error()
			return reply
		}
		if err := u.client.Input(p.Keys); err != nil {
			reply.Error = err.Error()
		}
		return reply

	case process.MethodExecute:
		var p process.ExecuteParams
		if err := json.Unmarshal(params, &p); err != nil {
			reply.Error = err.Error()
			return reply
		}
		raw, err := u.client.Raw(u.ctx, p)
		if err != nil {
			reply.Error = err.Error()
			return reply
		}
		reply.Result = raw
		return reply

	case process.MethodUiAttach:
		res, err := u.client.UiAttach(u.ctx, process.UiAttachParams{})
		if err != nil {
			reply.Error = err.Error()
			return reply
		}
		result, _ := json.Marshal(res)
		reply.Result = result
		return reply

	case process.MethodUiReady:
		var p process.UiReadyParams
		if err := json.Unmarshal(params, &p); err != nil {
			reply.Error = err.Error()
			return reply
		}

		res, err := u.client.UiReady(u.ctx, p)
		if err != nil {
			reply.Error = err.Error()
			return reply
		}

		result, _ := json.Marshal(res)
		reply.Result = result
		return reply

	default:
		reply.Error = fmt.Sprintf("unknown method %q", method)
		return reply
	}
}
