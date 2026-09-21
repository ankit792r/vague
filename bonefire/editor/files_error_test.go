package editor

import (
	"errors"
	"testing"

	"vague/bonefire/buffer"
)

func TestWriteFileErrorNoFileName(t *testing.T) {
	t.Parallel()

	err := WriteFileError("", buffer.ErrNoFileName)
	if err.Error() != "No file name" {
		t.Fatalf("got %q", err.Error())
	}
}

func TestWriteFileErrorWrapsOther(t *testing.T) {
	t.Parallel()

	err := WriteFileError("", errors.New("disk full"))
	if err.Error() != "write: disk full" {
		t.Fatalf("got %q", err.Error())
	}
}
