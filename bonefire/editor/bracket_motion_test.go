package editor

import (
	"testing"

	"vague/bonefire/text"
)

func TestMatchingBracket(t *testing.T) {
	t.Parallel()

	tex := text.New([]byte("(a [b] c)\n"))
	at := tex.OffsetOf(text.Point{Line: 0, Col: 0})
	got := tex.PointOf(moveMatchingBracket(tex, at, 1))
	if got.Col != 8 {
		t.Fatalf("%% on (: col=%d, want 8", got.Col)
	}

	at = tex.OffsetOf(text.Point{Line: 0, Col: 4})
	got = tex.PointOf(moveMatchingBracket(tex, at, 1))
	if got.Col != 6 {
		t.Fatalf("%% on [: col=%d, want 6", got.Col)
	}
}
