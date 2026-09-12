package backbone

import (
	"encoding/json"
	"fmt"
	"vague/backbone/process"
)

// A Notification is a server-initiated message.
type Notification struct {
	Method string
	Params json.RawMessage
}

func (c *Client) readLoop() {
	for {
		payload, err := process.ReadFrame(c.reader)
		if err != nil {
			c.fail(err)
			return
		}

		var msg process.Message
		if err := json.Unmarshal(payload, &msg); err != nil {
			c.fail(fmt.Errorf("decode message: %w", err))
			return
		}

		switch msg.Kind {
		case process.KindResponse:
			c.deliver(msg)

		case process.KindNotify:
			c.queue(Notification{Method: msg.Method, Params: msg.Params})

		default:
			// The server does not call the client, so this is a protocol
			// violation rather than something to route.
			c.fail(fmt.Errorf("unexpected message kind %q from server", msg.Kind))
			return
		}

	}

}

// deliver hands a response to the call waiting on its ID.
func (c *Client) deliver(msg process.Message) {
	c.mu.Lock()
	waiter, ok := c.pending[msg.ID]
	delete(c.pending, msg.ID)
	c.mu.Unlock()

	if !ok {
		// A call that gave up, or a reply to a malformed frame, which
		// carries no usable ID.
		return
	}

	waiter <- msg
}

// queue offers a notification, discarding the oldest if nobody is reading.
func (c *Client) queue(note Notification) {
	// TODO: Handle this
}
