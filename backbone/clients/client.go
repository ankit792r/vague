package backbone

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"vagues/backbone/process"
)

var ErrClosed = errors.New("connection closed")

type Client struct {
	conn   net.Conn
	reader *bufio.Reader

	mu      sync.Mutex
	writeMu sync.Mutex
	nextId  uint64

	pending map[uint64]chan process.Message

	err error

	closeOnce sync.Once
	done      chan struct{}
}

func ClientConnect() (*Client, error) {
	conn, err := process.Dial()
	if err != nil {
		return nil, err
	}

	client := &Client{
		conn:    conn,
		reader:  bufio.NewReader(conn),
		nextId:  1,
		pending: make(map[uint64]chan process.Message),
		done:    make(chan struct{}),
	}

	// Here start the reading loop
	return client, nil
}

func (c *Client) Close() error {
	c.closeOnce.Do(func() {
		close(c.done)
	})

	return c.conn.Close()
}

func (c *Client) call(ctx context.Context, method string, params any, result any) error {
	msg, err := func() (process.Message, error) {
		c.mu.Lock()
		defer c.mu.Unlock()

		// check error
		if c.err != nil {
			return process.Message{}, c.err
		}

		id := c.nextId
		c.nextId++

		msg, err := process.NewRequest(id, method, params)
		if err != nil {
			return process.Message{}, err
		}

		c.pending[id] = make(chan process.Message, 1)

		return msg, nil
	}()
	if err != nil {
		return err
	}

	waiter := c.waiterFor(msg.ID)

	if err := c.write(msg); err != nil {
		c.forget(msg.ID)
		return fmt.Errorf("send %s: %w", method, err)
	}

	select {
	case res, ok := <-waiter:
		if !ok {
			return c.reason()
		}

		if res.Error != "" {
			return errors.New(res.Error)
		}

		if result == nil {
			return nil
		}

		return res.DecodeResult(res)

	case <-ctx.Done():
		c.forget(msg.ID)
		return ctx.Err()

	case <-c.done:
		c.forget(msg.ID)
		return ctx.Err()
	}
}

func (c *Client) write(msg process.Message) error {
	c.writeMu.Lock()
	defer c.writeMu.Unlock()

	return process.WriteFrame(c.conn, msg)
}

func (c *Client) waiterFor(id uint64) chan process.Message {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.pending[id]
}

func (c *Client) forget(id uint64) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.pending, id)
}

func (c *Client) reason() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.err != nil {
		return c.err
	}

	return ErrClosed
}

// notify sends a message that expects no reply.
func (c *Client) notify(method string, params any) error {
	msg, err := process.NewNotify(method, params)
	if err != nil {
		return err
	}

	if err := c.write(msg); err != nil {
		return fmt.Errorf("send %s: %w", method, err)
	}

	return nil
}

// fail records why the connection ended and wakes every waiting call.
func (c *Client) fail(cause error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.err == nil {
		if errors.Is(cause, io.EOF) || errors.Is(cause, net.ErrClosed) {
			c.err = ErrClosed
		} else {
			c.err = cause
		}
	}

	for id, waiter := range c.pending {
		close(waiter)
		delete(c.pending, id)
	}
}
