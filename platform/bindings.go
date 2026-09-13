package platform

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	client "vague/backbone/clients"
	"vague/backbone/process"
)

// NativeReply is the value returned to the frontend's callNative Promise.
type NativeReply struct {
	ID     uint64          `json:"id"`
	Result json.RawMessage `json:"result,omitempty"`
	Error  string          `json:"error,omitempty"`
}

// Bindings is the webview → IPC bridge. The frontend calls the bound
// callNative(id, method, params); CallNative dispatches on method and
// forwards to the daemon over the Unix socket.
type Bindings struct {
	ctx context.Context

	mu     sync.Mutex
	client *client.Client
}

func NewBindings(ctx context.Context) *Bindings {
	if ctx == nil {
		ctx = context.Background()
	}
	return &Bindings{ctx: ctx}
}

// CallNative is bound as window.callNative. id is echoed so the
// frontend can correlate replies; method selects the IPC call.
func (b *Bindings) CallNative(id uint64, method string, params json.RawMessage) NativeReply {
	result, err := b.dispatch(method, params)
	if err != nil {
		fmt.Println("CallNative error", err)
		return NativeReply{ID: id, Error: err.Error()}
	}
	fmt.Println("CallNative", id, method, result)
	return NativeReply{ID: id, Result: result}
}

func (b *Bindings) dispatch(method string, params json.RawMessage) (json.RawMessage, error) {
	switch method {
	case process.MethodExecute:
		return b.execute(params)
	case process.MethodInput:
		return b.input(params)
	default:
		return nil, fmt.Errorf("unknown method %q", method)
	}
}

func (b *Bindings) execute(params json.RawMessage) (json.RawMessage, error) {
	var p process.ExecuteParams
	if err := decodeParams(params, &p); err != nil {
		return nil, err
	}

	c, err := b.ipc()
	if err != nil {
		return nil, err
	}

	result, err := c.Raw(b.ctx, p)
	if err != nil {
		b.dropClient()
		return nil, err
	}
	return result, nil
}

func (b *Bindings) input(params json.RawMessage) (json.RawMessage, error) {
	var p process.InputParams
	if err := decodeParams(params, &p); err != nil {
		return nil, err
	}
	if p.Keys == "" {
		return nil, errors.New("input requires keys")
	}

	c, err := b.ipc()
	if err != nil {
		return nil, err
	}

	if err := c.Input(p.Keys); err != nil {
		b.dropClient()
		return nil, err
	}

	return json.RawMessage(`{"ok":true}`), nil
}

func (b *Bindings) ipc() (*client.Client, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.client != nil {
		return b.client, nil
	}

	c, err := client.ClientConnect()
	if err != nil {
		return nil, err
	}
	b.client = c
	return c, nil
}

func (b *Bindings) dropClient() {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.client == nil {
		return
	}
	_ = b.client.Close()
	b.client = nil
}

func (b *Bindings) Close() {
	b.dropClient()
}

func decodeParams(raw json.RawMessage, dst any) error {
	if len(raw) == 0 || string(raw) == "null" {
		return nil
	}
	if err := json.Unmarshal(raw, dst); err != nil {
		return fmt.Errorf("decode params: %w", err)
	}
	return nil
}
