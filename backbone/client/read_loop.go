package client

import (
	"encoding/json"
	"fmt"
	"vague/backbone/process"
)

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
}
