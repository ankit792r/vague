package frame

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

func (f *Frame) HostRequest(id uint64, method string, params json.RawMessage) HostReply {
	reply := HostReply{ID: id}

	switch method {
	case process.MethodInput:
		var p process.InputParams
		if err := json.Unmarshal(params, &p); err != nil {
			reply.Error = err.Error()
			return reply
		}
		if err := f.client.Input(p.Keys); err != nil {
			reply.Error = err.Error()
		}
		return reply

	case process.MethodExecute:
		var p process.ExecuteParams
		if err := json.Unmarshal(params, &p); err != nil {
			reply.Error = err.Error()
			return reply
		}
		raw, err := f.client.Raw(f.ctx, p)
		if err != nil {
			reply.Error = err.Error()
			return reply
		}
		reply.Result = raw
		return reply

	case process.MethodFrameAttach:
		res, err := f.client.FrameAttach(f.ctx, process.AttachParams{})
		if err != nil {
			reply.Error = err.Error()
			return reply
		}
		result, _ := json.Marshal(res)
		reply.Result = result
		return reply

	case process.MethodFrameReady:
		var p process.FrameReadyParams
		if err := json.Unmarshal(params, &p); err != nil {
			reply.Error = err.Error()
			return reply
		}

		res, err := f.client.FrameReady(f.ctx, p)
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

func (f *Frame) IpcBinding(id uint64, method string, params json.RawMessage) {

	fmt.Println(method)

}
