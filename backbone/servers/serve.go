package backbone

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"time"
	"vague/backbone/process"
	"vague/backbone/session"
)

func (s *Server) serveSession(ctx context.Context, sess *session.Session) {
	// The request context dies with the connection, so a client that
	// disconnects mid-call does not leave a goroutine parked in Do.
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	sessionWriteComplete := make(chan struct{})

	go func() {
		defer close(sessionWriteComplete)
		sess.WriteLoop()
	}()

	defer func() {
		// here we may release the window or frame in client
		sess.Close()

		// Let the writer flush before the connection goes away, but do not
		// hang on a client that has stopped reading.
		select {
		case <-sessionWriteComplete:
		case <-time.After(time.Second):

		}

		sess.Conn.Close()
	}()

	for {
		payload, err := process.ReadFrame(sess.Reader)
		if err != nil {
			if !errors.Is(err, io.EOF) && !errors.Is(err, net.ErrClosed) {
				slog.Debug("session read ended", "error", err)
			}
			break
		}

		var msg process.Message
		if err := json.Unmarshal(payload, &msg); err != nil {
			// The ID is unknown, so the client cannot correlate this. It is
			// still worth sending so a malformed call is not a silent hang.
			sess.Reply(0, nil, fmt.Errorf("malformed message: %w", err))
			continue
		}

		if err := msg.Validate(); err != nil {
			sess.Reply(msg.ID, nil, err)
			continue
		}

		// Dispatch here
		s.dispatch(ctx, sess, msg)
	}
}
