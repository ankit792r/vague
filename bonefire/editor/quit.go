package editor

import "errors"

// ErrNotSaved means quit was refused because the buffer has unsaved changes.
var ErrNotSaved = errors.New("No write since last change (add ! to override)")
