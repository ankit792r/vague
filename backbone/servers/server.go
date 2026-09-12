package backbone

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"sync"
	"time"
	"vague/backbone/process"
	"vague/backbone/session"
)

type Server struct {
	wg sync.WaitGroup
	mu sync.Mutex

	sessions map[*session.Session]struct{}

	nextSessionID uint64
}

func NewServer() *Server {
	return &Server{
		sessions: make(map[*session.Session]struct{}),
	}
}

func (s *Server) Serve(ctx context.Context, listener *process.Listener) error {
	defer listener.Close()

	// Start the server runtime HERE
	slog.Info("Server listening")

	go func() {
		<-ctx.Done()
		listener.Close()
		s.closeAllConns()
	}()

	acceptErr := s.accept(ctx, listener)

	s.wg.Wait()

	return acceptErr
}

func (s *Server) accept(ctx context.Context, listener net.Listener) error {
	for {
		conn, err := listener.Accept()
		if err != nil {
			// A closed listener is the normal shutdown path.
			if errors.Is(err, net.ErrClosed) || ctx.Err() != nil {
				return nil
			}

			// Back off on transient failures instead of spinning a core.
			var netErr net.Error
			if errors.As(err, &netErr) && netErr.Timeout() {
				slog.Warn("accept failed, retrying", "error", err)
				time.Sleep(10 * time.Millisecond)
				continue
			}

			return fmt.Errorf("accept: %w", err)
		}

		s.nextSessionID++
		sess := session.NewSession(s.nextSessionID, conn)

		s.addConn(sess)
		s.wg.Add(1)

		go func() {
			defer s.wg.Done()
			defer s.removeConn(sess)

			s.serveSession(ctx, sess)
		}()
	}
}

// Session bookkeeping
func (s *Server) addConn(sess *session.Session) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.sessions[sess] = struct{}{}
}

func (s *Server) removeConn(sess *session.Session) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.sessions, sess)
}

func (s *Server) closeAllConns() {
	s.mu.Lock()
	defer s.mu.Unlock()

	for sess := range s.sessions {
		sess.Conn.Close()
	}
}
