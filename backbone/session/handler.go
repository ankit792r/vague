package session

import (
	"log/slog"
	"vague/backbone/process"
)

// send queues a message, reporting whether the session can still keep up.
func (s *Session) Send(msg process.Message) bool {
	select {
	case <-s.done:
		return false
	default:
	}

	select {
	case s.outbound <- msg:
		return true

	case <-s.done:
		return false

	default:
		slog.Warn("outbound queue full, dropping session")
		s.Conn.Close()

		return false
	}
}

// reply answers a request.
func (s *Session) Reply(id uint64, result any, cause error) {
	if cause != nil {
		s.Send(process.NewErrorResponse(id, cause))
		return
	}

	msg, err := process.NewResponse(id, result)
	if err != nil {
		s.Send(process.NewErrorResponse(id, err))
		return
	}

	s.Send(msg)
}
