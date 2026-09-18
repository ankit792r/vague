package editor

import (
	"testing"

	"vague/bonefire/text"
)

func TestMoveToLineStartAndEnd(t *testing.T) {
	t.Parallel()

	tex := text.New([]byte("hello"))

	start := tex.OffsetOf(text.Point{Line: 0, Col: 2})
	if got := tex.PointOf(moveToLineStart(tex, start)); got.Col != 0 {
		t.Fatalf("line start: got col %d, want 0", got.Col)
	}

	if got := tex.PointOf(moveToLineEnd(tex, start, false)); got.Col != 4 {
		t.Fatalf("line end: got col %d, want 4", got.Col)
	}

	if got := moveToLineEnd(tex, start, true); got != 5 {
		t.Fatalf("line end past: got offset %d, want 5", got)
	}
}

func TestFirstNonBlank(t *testing.T) {
	t.Parallel()

	tex := text.New([]byte("  hello"))
	mid := tex.OffsetOf(text.Point{Line: 0, Col: 4})

	if got := tex.PointOf(firstNonBlank(tex, mid)); got.Col != 2 {
		t.Fatalf("first non-blank: got col %d, want 2", got.Col)
	}

	blank := text.New([]byte("   "))
	if got := blank.PointOf(firstNonBlank(blank, 0)); got.Col != 0 {
		t.Fatalf("blank line: got col %d, want 0", got.Col)
	}
}

func TestMoveLeftRight(t *testing.T) {
	t.Parallel()

	tex := text.New([]byte("hello"))

	got := moveRight(tex, 0, 1, false)
	if tex.PointOf(got).Col != 1 {
		t.Fatalf("move right: got col %d, want 1", tex.PointOf(got).Col)
	}

	got = moveLeft(tex, got, 1)
	if tex.PointOf(got).Col != 0 {
		t.Fatalf("move left: got col %d, want 0", tex.PointOf(got).Col)
	}
}

func TestMoveVerticalKeepsDesiredColumn(t *testing.T) {
	t.Parallel()

	tex := text.New([]byte("short\nmuchlonger"))
	start := tex.OffsetOf(text.Point{Line: 1, Col: 7})

	got := moveVertical(tex, start, 7, -1)
	point := tex.PointOf(got)

	if point.Line != 0 || point.Col != 4 {
		t.Fatalf("move up: got line=%d col=%d, want line 0 col 4", point.Line, point.Col)
	}
}

func TestMoveVerticalVisualWithinWrappedLine(t *testing.T) {
	t.Parallel()

	tex := text.New([]byte("abcdefgh"))
	view := layoutView(tex, 4, 0, true)

	start := tex.OffsetOf(text.Point{Line: 0, Col: 0})
	got := moveVerticalVisual(view, tex, start, 0, 1, false)
	point := tex.PointOf(got)

	if point.Line != 0 || point.Col != 4 {
		t.Fatalf("move down: got line=%d col=%d, want line 0 col 4", point.Line, point.Col)
	}

	got = moveVerticalVisual(view, tex, got, 0, -1, false)
	point = tex.PointOf(got)

	if point.Line != 0 || point.Col != 0 {
		t.Fatalf("move up: got line=%d col=%d, want line 0 col 0", point.Line, point.Col)
	}
}

func TestCursorScreenPosWithWrap(t *testing.T) {
	t.Parallel()

	view := layoutView(text.New([]byte("abcdefghijklmn")), 13, 0, true)
	row, col, visible := cursorScreenPos(view.Meta, text.Point{Line: 0, Col: 13})

	if !visible || row != 1 || col != 0 {
		t.Fatalf("got row=%d col=%d visible=%v, want row=1 col=0", row, col, visible)
	}
}
