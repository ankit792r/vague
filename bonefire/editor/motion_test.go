package editor

import "testing"

func TestMoveLeftRight(t *testing.T) {
	t.Parallel()

	lines := []string{"hello"}

	got := moveRight(lines, Point{Line: 0, Col: 0}, 1)
	if got.Col != 1 {
		t.Fatalf("move right: got col %d, want 1", got.Col)
	}

	got = moveLeft(lines, got, 1)
	if got.Col != 0 {
		t.Fatalf("move left: got col %d, want 0", got.Col)
	}
}

func TestMoveVerticalKeepsDesiredColumn(t *testing.T) {
	t.Parallel()

	lines := []string{"short", "muchlonger"}

	got := moveVertical(lines, Point{Line: 1, Col: 7}, 7, -1)
	if got.Line != 0 || got.Col != 4 {
		t.Fatalf("move up: got %+v, want line 0 col 4", got)
	}
}

func TestCursorScreenPosWithWrap(t *testing.T) {
	t.Parallel()

	view := layoutView([]string{"abcdefghijklmn"}, 13, 0, true)
	row, col, visible := cursorScreenPos(view.Meta, Point{Line: 0, Col: 13})

	if !visible || row != 1 || col != 0 {
		t.Fatalf("got row=%d col=%d visible=%v, want row=1 col=0", row, col, visible)
	}
}
