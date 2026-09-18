package editor

import (
	"testing"

	"vague/bonefire/text"
)

func TestMoveWordForward(t *testing.T) {
	t.Parallel()

	tex := text.New([]byte("hello world"))
	start := tex.LineStart(0)

	got := tex.PointOf(moveWordForward(tex, start, 1))
	if got.Col != 6 {
		t.Fatalf("w from start: col = %d, want 6", got.Col)
	}

	got = tex.PointOf(moveWordForward(tex, tex.OffsetOf(text.Point{Line: 0, Col: 6}), 1))
	if got.Col != 6 {
		t.Fatalf("w on last word start: col = %d, want 6 (no move)", got.Col)
	}
}

func TestMoveWordBack(t *testing.T) {
	t.Parallel()

	tex := text.New([]byte("hello world"))
	from := tex.OffsetOf(text.Point{Line: 0, Col: 6})

	got := tex.PointOf(moveWordBack(tex, from, 1))
	if got.Col != 0 {
		t.Fatalf("b from world: col = %d, want 0", got.Col)
	}
}

func TestMoveWordEnd(t *testing.T) {
	t.Parallel()

	tex := text.New([]byte("hello world"))
	start := tex.LineStart(0)

	got := tex.PointOf(moveWordEnd(tex, start, 1))
	if got.Col != 4 {
		t.Fatalf("e from start: col = %d, want 4", got.Col)
	}

	got = tex.PointOf(moveWordEnd(tex, tex.OffsetOf(text.Point{Line: 0, Col: 4}), 1))
	if got.Col != 10 {
		t.Fatalf("e from end of hello: col = %d, want 10", got.Col)
	}
}

func TestMoveWordAcrossLines(t *testing.T) {
	t.Parallel()

	tex := text.New([]byte("one two\nthree"))
	at := tex.OffsetOf(text.Point{Line: 0, Col: 3})

	got := tex.PointOf(moveWordForward(tex, at, 1))
	if got.Line != 0 || got.Col != 4 {
		t.Fatalf("w mid-word: got line=%d col=%d, want 0,4", got.Line, got.Col)
	}

	at = tex.OffsetOf(text.Point{Line: 0, Col: 7})
	got = tex.PointOf(moveWordForward(tex, at, 1))
	if got.Line != 1 || got.Col != 0 {
		t.Fatalf("w across newline: got line=%d col=%d, want 1,0", got.Line, got.Col)
	}
}

func TestMoveWordPunctuation(t *testing.T) {
	t.Parallel()

	tex := text.New([]byte("hi, there"))
	start := tex.LineStart(0)

	got := tex.PointOf(moveWordForward(tex, start, 1))
	if got.Col != 2 {
		t.Fatalf("w to punct: col = %d, want 2", got.Col)
	}

	got = tex.PointOf(moveWordForward(tex, tex.OffsetOf(text.Point{Line: 0, Col: 2}), 1))
	if got.Col != 4 {
		t.Fatalf("w past punct: col = %d, want 4", got.Col)
	}
}

func TestWordMotionKeys(t *testing.T) {
	t.Parallel()

	ed := NewEditor()
	buf := ed.Scratch("*scratch*")
	buf.Text.SetBytes([]byte("hello world"))

	frame, err := ed.NewFrame(80, 10)
	if err != nil {
		t.Fatal(err)
	}

	win := ed.Windows[frame.ActiveWindowID]

	if err := ed.HandleInput(frame.ID, "w"); err != nil {
		t.Fatal(err)
	}

	point := buf.Text.PointOf(windowCursor(win))
	if point.Col != 6 {
		t.Fatalf("after w: col = %d, want 6", point.Col)
	}

	if err := ed.HandleInput(frame.ID, "e"); err != nil {
		t.Fatal(err)
	}

	point = buf.Text.PointOf(windowCursor(win))
	if point.Col != 10 {
		t.Fatalf("after e: col = %d, want 10", point.Col)
	}

	if err := ed.HandleInput(frame.ID, "b"); err != nil {
		t.Fatal(err)
	}

	point = buf.Text.PointOf(windowCursor(win))
	if point.Col != 6 {
		t.Fatalf("after b: col = %d, want 6", point.Col)
	}
}
