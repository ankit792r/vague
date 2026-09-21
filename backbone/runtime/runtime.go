package runtime

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"runtime/debug"
	"sync"
	"vague/bonefire/workspace"
)

// ErrShutdown is returned by Do once the runtime has shut down.
var ErrShutdown = errors.New("editor runtime is shut down")

type reply struct {
	val any
	err error
}

type job struct {
	fn    func(*workspace.Workspace) (any, error)
	reply chan reply // buffered, cap 1 - a sender must never block
}

// Runtime serialises all access to editor state onto a single goroutine.
type Runtime struct {
	ws *workspace.Workspace // only the loop goroutine may touch this
	jobs   chan job

	quit     chan struct{}
	quitOnce sync.Once
	done     chan struct{}
}

func NewRuntime() *Runtime {
	return &Runtime{
		ws:   workspace.New(),
		jobs: make(chan job, 64),
		quit: make(chan struct{}),
		done: make(chan struct{}),
	}
}

// Do runs fn on the runtime goroutine and waits for its result.
func (r *Runtime) Do(ctx context.Context, fn func(*workspace.Workspace) (any, error)) (any, error) {
	j := job{fn: fn, reply: make(chan reply, 1)}

	select {
	case r.jobs <- j:
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-r.quit:
		return nil, ErrShutdown
	}

	select {
	case res := <-j.reply:
		return res.val, res.err
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-r.quit:
		return nil, ErrShutdown
	}
}

// Run drives the runtime until Shutdown is called.
func (r *Runtime) Run() {
	defer close(r.done)

	for {
		select {
		case j := <-r.jobs:
			j.reply <- r.exec(j.fn)
		case <-r.quit:
			return
		}
	}
}

// Shutdown stops the loop. It is safe to call more than once.
func (r *Runtime) Shutdown() {
	r.quitOnce.Do(func() { close(r.quit) })
}

// Wait blocks until the loop has stopped.
func (r *Runtime) Wait() { <-r.done }

func (r *Runtime) exec(fn func(*workspace.Workspace) (any, error)) (res reply) {
	defer func() {
		if p := recover(); p != nil {
			slog.Error("command panicked",
				"panic", p,
				"stack", string(debug.Stack()),
			)
			res = reply{err: fmt.Errorf("internal error: %v", p)}
		}
	}()

	val, err := fn(r.ws)

	return reply{val: val, err: err}
}
