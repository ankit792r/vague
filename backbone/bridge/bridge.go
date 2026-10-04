package bridge

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/abemedia/go-webview"
)

type MethodParams struct {
	args  []string
	bang  bool
	count uint
}
type Handler func(params MethodParams) (any, error)

type Bridge struct {
	ctx      context.Context
	webview  webview.WebView
	mu       sync.RWMutex
	handlers map[string]Handler
}

func NewBridge(w webview.WebView, ctx context.Context) *Bridge {
	return &Bridge{
		ctx:      ctx,
		webview:  w,
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
func (b *Bridge) invoke(methodName string, params MethodParams) any {
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

	b.webview.Eval(js)
}

// This will attach bridge in webview
func AttachBridgeToWebView(w webview.WebView) {
	b := NewBridge(w)

	// Here will register other handler
	// b.RegisterFrame()
	// b.RegisterInput()
	// b.RegisterCommand()
	// handlers.RegisterFrameHandlers(b)
	RegisterFrameHandler(b)

	err := b.webview.Bind("hostInvoke", b.invoke)
	if err != nil {
		log.Fatalf("Failed to attach UI bridge: %v", err)
	}
}
