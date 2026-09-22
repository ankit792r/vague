package editor

import "errors"

// ErrConfirmPending means the client should prompt and send confirm_answer.
var ErrConfirmPending = errors.New("confirm pending")

// ErrNotSaved means quit was refused because the buffer has unsaved changes.
var ErrNotSaved = errors.New("No write since last change (add ! to override)")
