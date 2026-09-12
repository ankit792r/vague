package session

import (
	"bufio"
	"errors"
	"log/slog"
	"net"
	"sync"
	"vagues/backbone/process"
)

type Session struct {
	Id       uint64
	Conn     net.Conn
	Reader   *bufio.Reader
	outbound chan process.Message

	done     chan struct{}
	doneOnce sync.Once
}

func NewSession(id uint64, conn net.Conn) *Session {
	return &Session{
		Id:       id,
		Conn:     conn,
		Reader:   bufio.NewReader(conn),
		outbound: make(chan process.Message),

		done: make(chan struct{}),
	}
}

func (s *Session) Close() {
	s.doneOnce.Do(func() {
		close(s.done)
	})
}


// writeLoop drains the outbound queue until the session closes. Frames are
// written by one goroutine only, so they cannot interleave.
func (s *Session) WriteLoop() {
	for {
		select {
		case msg := <-s.outbound:
			if !s.write(msg) {
				return
			}

		case <-s.done:
			// Flush what is already queued so a reply sent just before
			// shutdown still reaches the client.
			for {
				select {
				case msg := <-s.outbound:
					if !s.write(msg) {
						return
					}
				default:
					return
				}
			}
		}
	}
}

func (s *Session) write(msg process.Message) bool {
	if err := process.WriteFrame(s.Conn, msg); err != nil {
		if !errors.Is(err, net.ErrClosed) {
			slog.Debug("write failed", "error", err)
		}

		return false
	}

	return true
}
