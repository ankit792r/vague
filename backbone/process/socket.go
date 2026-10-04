package process

import (
	"errors"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"sync"
	"syscall"
)

// ErrAlreadyRunning means another server holds the lock for this socket.
var ErrAlreadyRunning = errors.New("vague server is already running")

// ErrNotRunning means nothing is listening on the socket.
var ErrNotRunning = errors.New("vague server is not running")

func SocketPath() (string, error) {
	runtimeDir := os.Getenv("XDG_RUNTIME_DIR")
	if runtimeDir == "" {
		return "", fmt.Errorf("XDG_RUNTIME_DIR is not set")

	}

	return filepath.Join(runtimeDir, "vague.sock"), nil
}

// Listener owns the server socket and the lock that guarantees it is unique.
type Listener struct {
	net.Listener

	lock *os.File
	path string

	closeOnce sync.Once
	closeErr  error
}

// Path is the filesystem path of the socket.
func (l *Listener) Path() string { return l.path }

// Listen binds the server socket, refusing to start if another server is
// already up.
//
// Exclusivity comes from an flock on a sidecar file rather than from probing
// the socket: the kernel drops the lock when the holder dies, so a server
// killed with SIGKILL leaves no stale state, and there is no window in which
// two servers can both decide the socket is theirs.
func Listen() (*Listener, error) {
	socketPath, err := SocketPath()
	if err != nil {
		return nil, err
	}

	lock, err := os.OpenFile(socketPath+".lock", os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		return nil, fmt.Errorf("open lock file: %w", err)
	}

	if err := syscall.Flock(int(lock.Fd()), syscall.LOCK_EX|syscall.LOCK_NB); err != nil {
		lock.Close()

		if errors.Is(err, syscall.EWOULDBLOCK) {
			return nil, ErrAlreadyRunning
		}
		return nil, fmt.Errorf("lock socket: %w", err)
	}

	// Holding the lock means any socket file present is left over from a
	// server that died, so removing it cannot disturb a live one.
	if err := os.Remove(socketPath); err != nil && !os.IsNotExist(err) {
		lock.Close()
		return nil, fmt.Errorf("remove stale socket: %w", err)
	}

	inner, err := net.Listen("unix", socketPath)
	if err != nil {
		lock.Close()
		return nil, fmt.Errorf("listen on %s: %w", socketPath, err)
	}

	// XDG_RUNTIME_DIR is already 0700, but the socket should not depend on
	// that for its access control.
	if err := os.Chmod(socketPath, 0600); err != nil {
		inner.Close()
		lock.Close()
		return nil, fmt.Errorf("secure socket: %w", err)
	}

	return &Listener{Listener: inner, lock: lock, path: socketPath}, nil
}

// Close stops listening, unlinks the socket, and releases the lock.
//
// Shutdown closes the listener from the goroutine watching for cancellation
// as well as from the serve loop unwinding, so this has to be safe to call
// concurrently and more than once.
func (l *Listener) Close() error {
	l.closeOnce.Do(func() {
		l.closeErr = l.Listener.Close() // also unlinks the socket file

		syscall.Flock(int(l.lock.Fd()), syscall.LOCK_UN)
		l.lock.Close()
	})

	return l.closeErr
}

// Dial connects to a running server.
func Dial() (net.Conn, error) {
	socketPath, err := SocketPath()
	if err != nil {
		return nil, err
	}

	conn, err := net.Dial("unix", socketPath)
	if err != nil {
		if errors.Is(err, syscall.ENOENT) || errors.Is(err, syscall.ECONNREFUSED) {
			return nil, ErrNotRunning
		}
		return nil, fmt.Errorf("connect to vague server: %w", err)
	}

	return conn, nil
}
