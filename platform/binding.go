package platform

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

func (f *Frame) HostRequest(id uint64, method string, params json.RawMessage) (HostReply, error) {
	reply := HostReply{ID: id}

	switch method {
	case process.MethodInput:
		var p process.InputParams
		if err := json.Unmarshal(params, &p); err != nil {
			reply.Error = err.Error()
			return reply, nil
		}
		if err := f.client.Input(p.Keys); err != nil {
			reply.Error = err.Error()
		}
		return reply, nil

	case process.MethodExecute:
		var p process.ExecuteParams
		if err := json.Unmarshal(params, &p); err != nil {
			reply.Error = err.Error()
			return reply, nil
		}
		raw, err := f.client.Raw(f.ctx, p)
		if err != nil {
			reply.Error = err.Error()
			return reply, nil
		}
		reply.Result = raw
		return reply, nil

	case process.MethodFrameAttach:
		res, err := f.client.FrameAttach(f.ctx, process.AttachParams{})
		if err != nil {
			reply.Error = err.Error()
			return reply, nil
		}
		result, _ := json.Marshal(res)
		reply.Result = result
		return reply, nil

	default:
		reply.Error = fmt.Sprintf("unknown method %q", method)
		return reply, nil
	}
}

func (f *Frame) IpcBinding(id uint64, method string, params json.RawMessage) {

	fmt.Println(method)

}
