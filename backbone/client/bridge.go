package client

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/abemedia/go-webview"
)

type MethodParams struct {
	Args  []string
	Bang  bool
	Count uint
}
type Handler func(params MethodParams) (any, error)

type Bridge struct {
	ctx      context.Context
	cli      *Client
	Webview  webview.WebView
	mu       sync.RWMutex
	handlers map[string]Handler
}

func NewBridge(ctx context.Context, w webview.WebView, cli *Client) *Bridge {
	return &Bridge{
		ctx:      ctx,
		cli:      cli,
		Webview:  w,
		handlers: make(map[string]Handler),
	}
}

// Register adds a named RPC method callable from JavaScript as directorInvoke(method, payloadJSON).
func (b *Bridge) Register(method string, h Handler) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.handlers[method] = h
}

// TODO: we have to fix the return type of **invoke()**
func (b *Bridge) Invoke(methodName string, params MethodParams) any {
	b.mu.RLock()
	h, ok := b.handlers[methodName]
	b.mu.RUnlock()
	if !ok {
		return fmt.Sprintf("got unknown method: %s \n", methodName)
	}

	data, err := h(params)

	if err != nil {
		return err.Error()
	}

	return data
}

// Emit to client
func (b *Bridge) Emit(event string, data any) {
	encoded, err := json.Marshal(data)
	if err != nil {
		return
	}

	js := fmt.Sprintf(
		`(function(){var fn=window.onHostEvent;if(typeof fn==="function"){fn(%q,%s);}})();`,
		event,
		string(encoded),
	)

	b.Webview.Eval(js)
}
