package process

import (
	"encoding/json"
	"fmt"
)

// Kind distinguishes the three things that can cross the socket.
//
// JSON-RPC infers the kind from which fields are present, which is compact
// but makes a malformed frame look like a different valid message. Naming the
// kind costs a few bytes and means a decoder can reject nonsense instead of
// silently reinterpreting it.
type Kind string

const (
	// KindRequest expects exactly one KindResponse with a matching ID.
	KindRequest Kind = "request"

	// KindResponse answers a request.
	KindResponse Kind = "response"

	// KindNotify expects no reply and travels in either direction:
	// clients send input and resize, the server sends redraws.
	KindNotify Kind = "notify"
)

// A Message is the single envelope on the wire. Which fields are meaningful
// depends on Kind.
type Message struct {
	Kind Kind `json:"kind"`

	// ID correlates a request with its response. It is unset on
	// notifications.
	ID uint64 `json:"id,omitempty"`

	// Method names the call on requests and notifications.
	Method string `json:"method,omitempty"`

	Params json.RawMessage `json:"params,omitempty"`

	Result json.RawMessage `json:"result,omitempty"`
	Error  string          `json:"error,omitempty"`
}

// NewRequest builds a request, encoding its parameters.
func NewRequest(id uint64, method string, params any) (Message, error) {
	raw, err := Encode(params)
	if err != nil {
		return Message{}, err
	}

	return Message{Kind: KindRequest, ID: id, Method: method, Params: raw}, nil
}

// NewNotifn builds a notification.
func NewNotify(method string, params any) (Message, error) {
	raw, err := Encode(params)
	if err != nil {
		return Message{}, err
	}

	return Message{Kind: KindNotify, Method: method, Params: raw}, nil
}

// NewResponse builds a successful response.
func NewResponse(id uint64, result any) (Message, error) {
	raw, err := Encode(result)
	if err != nil {
		return Message{}, err
	}

	return Message{Kind: KindResponse, ID: id, Result: raw}, nil
}

// NewErrorResponse builds a failed response.
func NewErrorResponse(id uint64, cause error) Message {
	return Message{Kind: KindResponse, ID: id, Error: cause.Error()}
}

// Validate checks that a decoded message is internally consistent, so that
// handlers can assume the fields they care about are present.
func (m Message) Validate() error {
	switch m.Kind {
	case KindRequest:
		if m.ID == 0 {
			return fmt.Errorf("request %q has no id", m.Method)
		}
		if m.Method == "" {
			return fmt.Errorf("request %d has no method", m.ID)
		}

	case KindResponse:
		if m.ID == 0 {
			return fmt.Errorf("response has no id")
		}

	case KindNotify:
		if m.Method == "" {
			return fmt.Errorf("notify has no method")
		}

	default:
		return fmt.Errorf("unknown message kind %q", m.Kind)
	}

	return nil
}

// Encode marshals a value for a Params or Result field. Payloads are carried
// as raw JSON so that neither side has to marshal a value only to unmarshal
// it again into the type it wanted.
func Encode(v any) (json.RawMessage, error) {
	if v == nil {
		return nil, nil
	}

	raw, err := json.Marshal(v)
	if err != nil {
		return nil, fmt.Errorf("encode payload: %w", err)
	}

	return raw, nil
}

// DecodeParams unmarshals a request or notification's parameters into dst.
// Absent parameters leave dst untouched, so methods that take none can skip
// the call.
func (m Message) DecodeParams(dst any) error {
	if len(m.Params) == 0 {
		return nil
	}

	if err := json.Unmarshal(m.Params, dst); err != nil {
		return fmt.Errorf("decode params for %q: %w", m.Method, err)
	}

	return nil
}

// DecodeResult unmarshals a response's result into dst.
func (m Message) DecodeResult(dst any) error {
	if len(m.Result) == 0 {
		return nil
	}

	if err := json.Unmarshal(m.Result, dst); err != nil {
		return fmt.Errorf("decode result: %w", err)
	}

	return nil
}
