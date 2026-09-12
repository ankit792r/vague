package process

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
)

// MaxFrameSize caps a single message. Redraw batches for a large window are
// the biggest thing that crosses the wire, and they are nowhere near this.
// The cap exists so a corrupt length header cannot make the peer allocate
// unbounded memory.
const MaxFrameSize = 64 << 20

// Frames are a 4-byte big-endian payload length followed by that many bytes
// of JSON.
//
// Explicit lengths rather than newline delimiting, because bufio.Scanner caps
// lines at 64KB by default and a full-window redraw will exceed that. The
// cost is that the socket is no longer pokeable by hand with nc
const headerSize = 4

// WriteFrame encodes v as JSON and writes it as one frame.
//
// The header and payload go out in a single Write so that a frame cannot be
// torn apart by a concurrent writer on the same connection.
func WriteFrame(w io.Writer, v any) error {
	payload, err := json.Marshal(v)
	if err != nil {
		return fmt.Errorf("encode frame: %w", err)
	}

	if len(payload) > MaxFrameSize {
		return fmt.Errorf("frame too large: %d bytes", len(payload))
	}

	buf := make([]byte, headerSize+len(payload))
	binary.BigEndian.PutUint32(buf[:headerSize], uint32(len(payload)))
	copy(buf[headerSize:], payload)

	if _, err := w.Write(buf); err != nil {
		return fmt.Errorf("write frame: %w", err)
	}

	return nil
}

// ReadFrame reads one frame and returns its payload.
//
// It returns io.EOF when the peer closed cleanly between frames, and
// io.ErrUnexpectedEOF when it closed mid-frame.
func ReadFrame(r io.Reader) ([]byte, error) {
	var header [headerSize]byte
	if _, err := io.ReadFull(r, header[:]); err != nil {
		return nil, err
	}

	size := binary.BigEndian.Uint32(header[:])
	if size > MaxFrameSize {
		return nil, fmt.Errorf("frame too large: %d bytes", size)
	}

	payload := make([]byte, size)
	if _, err := io.ReadFull(r, payload); err != nil {
		return nil, err
	}

	return payload, nil
}
