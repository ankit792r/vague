package text

// Gravity decides which way a marker moves when text is inserted at its position.
type Gravity int

const (
	GravityLeft Gravity = iota
	GravityRight
)

// A Marker is a position that survives edits elsewhere in the buffer.
type Marker struct {
	Off     Offset
	Gravity Gravity
}

// AddMarker creates a marker at an offset and registers it for adjustment.
func (t *Text) AddMarker(off Offset, g Gravity) *Marker {
	m := &Marker{Off: t.Clamp(off), Gravity: g}
	t.markers = append(t.markers, m)
	return m
}

// RemoveMarker unregisters a marker.
func (t *Text) RemoveMarker(m *Marker) {
	for i, cur := range t.markers {
		if cur == m {
			t.markers = append(t.markers[:i], t.markers[i+1:]...)
			return
		}
	}
}

func (t *Text) adjustMarkers(a, b, newEnd Offset) {
	shift := newEnd - b

	for _, m := range t.markers {
		switch {
		case m.Off < a:
		case m.Off > b:
			m.Off += shift
		case a == b:
			if m.Gravity == GravityRight {
				m.Off = newEnd
			}
		case m.Off == a && m.Gravity == GravityLeft:
		case m.Off == b && m.Gravity == GravityRight:
			m.Off = newEnd
		default:
			m.Off = a
		}
	}
}
