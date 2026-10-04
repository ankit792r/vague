package runtime

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"runtime/debug"
	"sync"
)

var ErrShutdown = errors.New("editor runtime is shut down")

type DoFunc func(editor string, frames string, windows string) (any, error)

type reply struct {
	val any
	err error
}

type job struct {
	fn    DoFunc
	reply chan reply // buffered, cap 1 - a sender must never block
}

type Runtime struct {
	jobs chan job

	quit     chan struct{}
	quitOnce sync.Once
	done     chan struct{}

	// Editor *editor.Editor
	// Frames  map[uint64]*frame.Frame
	// Windows map[uint64]*window.Window

	Editor  string
	Frames  string
	Windows string

	CurrentFrameID uint64
	nextFrameID    uint64
	nextWindowID   uint64
}

func NewRuntime() *Runtime {
	return &Runtime{
		jobs: make(chan job, 64),
		quit: make(chan struct{}),
		done: make(chan struct{}),
	}
}

// Do runs fn on the runtime goroutine and waits for its result.
func (r *Runtime) Do(ctx context.Context, fn DoFunc) (any, error) {
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

func (r *Runtime) exec(fn DoFunc) (res reply) {
	defer func() {
		if p := recover(); p != nil {
			slog.Error("command panicked",
				"panic", p,
				"stack", string(debug.Stack()),
			)
			res = reply{err: fmt.Errorf("internal error: %v", p)}
		}
	}()

	val, err := fn(r.Editor, r.Frames, r.Windows)

	return reply{val: val, err: err}
}
