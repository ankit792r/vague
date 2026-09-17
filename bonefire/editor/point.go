package editor

// Point is a cursor position in buffer coordinates.
// Col is a rune index within the line, not a display column.
type Point struct {
	Line int
	Col  int
}
