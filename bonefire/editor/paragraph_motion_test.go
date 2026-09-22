package editor

import (
	"testing"

	"vague/bonefire/text"
)

func TestMoveParagraphForward(t *testing.T) {
	t.Parallel()

	tex := text.New([]byte("aa\n\nbb\ncc\n\ndd\n"))
	at := tex.OffsetOf(text.Point{Line: 0, Col: 0})

	got := tex.PointOf(moveParagraphForward(tex, at, 1))
	if got.Line != 2 || got.Col != 0 {
		t.Fatalf("} once: line=%d col=%d, want 2,0", got.Line, got.Col)
	}

	at = tex.OffsetOf(text.Point{Line: 2, Col: 0})
	got = tex.PointOf(moveParagraphForward(tex, at, 1))
	if got.Line != 5 {
		t.Fatalf("} to last para: line=%d, want 5", got.Line)
	}
}

func TestMoveParagraphBackward(t *testing.T) {
	t.Parallel()

	tex := text.New([]byte("aa\n\nbb\ncc\n\ndd\n"))
	at := tex.OffsetOf(text.Point{Line: 3, Col: 1 })

	got := tex.PointOf(moveParagraphBackward(tex, at, 1))
	if got.Line != 2 {
		t.Fatalf("{ to para start: line=%d, want 2", got.Line)
	}

	got = tex.PointOf(moveParagraphBackward(tex, at, 2))
	if got.Line != 0 {
		t.Fatalf("{ twice: line=%d, want 0", got.Line)
	}
}

func TestMoveSentenceForward(t *testing.T) {
	t.Parallel()

	tex := text.New([]byte("One. Two! Three?\n"))
	at := tex.OffsetOf(text.Point{Line: 0, Col: 0})

	got := tex.PointOf(moveSentenceForward(tex, at, 1))
	if got.Col != 5 {
		t.Fatalf(") once: col=%d, want 5", got.Col)
	}

	got = tex.PointOf(moveSentenceForward(tex, at, 2))
	if got.Col != 10 {
		t.Fatalf(") twice: col=%d, want 10", got.Col)
	}
}

func TestMoveSentenceBackward(t *testing.T) {
	t.Parallel()

	tex := text.New([]byte("One. Two! Three?\n"))
	at := tex.OffsetOf(text.Point{Line: 0, Col: 15 })

	got := tex.PointOf(moveSentenceBackward(tex, at, 1))
	if got.Col != 10 {
		t.Fatalf("( once: col=%d, want 10", got.Col)
	}

	got = tex.PointOf(moveSentenceBackward(tex, at, 2))
	if got.Col != 5 {
		t.Fatalf("( twice: col=%d, want 5", got.Col)
	}

	got = tex.PointOf(moveSentenceBackward(tex, at, 3))
	if got.Col != 0 {
		t.Fatalf("( thrice: col=%d, want 0", got.Col)
	}
}
