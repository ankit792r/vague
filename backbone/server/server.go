package server

import (
	"context"
	"fmt"
	"sync"
	"vague/backbone/process"
	"vague/backbone/runtime"
	"vague/backbone/session"
)

type Server struct {
	wg sync.WaitGroup
	mu sync.Mutex

	sessions map[*session.Session]struct{}
	runtime  *runtime.Runtime

	nextSessionID uint64
}

func NewServer() *Server {
	return &Server{
		sessions: make(map[*session.Session]struct{}),
		runtime:  runtime.NewRuntime(),
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

func ServerStart() {
	ctx := context.TODO()

	listener, err := process.Listen()
	if err != nil {
		fmt.Errorf("server: %w", err)
	}

	if err := NewServer().Serve(ctx, listener); err != nil {
		fmt.Errorf("server: %w", err)
	}
}

func ServerStop() {

}
