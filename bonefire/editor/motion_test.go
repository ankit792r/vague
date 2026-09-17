package editor

import (
	"testing"

	"vague/bonefire/text"
)

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

func TestCursorScreenPosWithWrap(t *testing.T) {
	t.Parallel()

	view := layoutView(text.New([]byte("abcdefghijklmn")), 13, 0, true)
	row, col, visible := cursorScreenPos(view.Meta, text.Point{Line: 0, Col: 13})

	if !visible || row != 1 || col != 0 {
		t.Fatalf("got row=%d col=%d visible=%v, want row=1 col=0", row, col, visible)
	}
}
