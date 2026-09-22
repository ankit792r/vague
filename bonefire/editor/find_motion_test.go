package editor

import (
	"testing"

	"vague/bonefire/text"
)

func TestFindForwardOnLine(t *testing.T) {
	t.Parallel()

	tex := text.New([]byte("foobar"))
	line := tex.Line(0)

	off, ok := findForwardOnLine(tex, 0, line, 0, false, 'o', 1)
	if !ok {
		t.Fatal("expected match")
	}
	if tex.PointOf(off).Col != 1 {
		t.Fatalf("f col = %d, want 1", tex.PointOf(off).Col)
	}

	off, ok = findForwardOnLine(tex, 0, line, 0, true, 'o', 1)
	if !ok || tex.PointOf(off).Col != 0 {
		t.Fatalf("t col = %d, want 0", tex.PointOf(off).Col)
	}

	off, ok = findForwardOnLine(tex, 0, line, 0, false, 'o', 2)
	if !ok || tex.PointOf(off).Col != 2 {
		t.Fatalf("2nd o col = %d, want 2", tex.PointOf(off).Col)
	}

	_, ok = findForwardOnLine(tex, 0, line, 0, true, 'f', 1)
	if ok {
		t.Fatal("t before first char should fail")
	}

	triple := text.New([]byte("ooo"))
	off, ok = findCharOnLine(triple, 0, charFindF, 'o', 3)
	if !ok || triple.PointOf(off).Col != 2 {
		t.Fatalf("3rd o col = %d, want 2", triple.PointOf(off).Col)
	}
}

func TestFindBackwardOnLine(t *testing.T) {
	t.Parallel()

	tex := text.New([]byte("abab"))
	at := tex.OffsetOf(text.Point{Line: 0, Col: 3})

	off, ok := findCharOnLine(tex, at, charFindBigF, 'a', 1)
	if !ok || tex.PointOf(off).Col != 2 {
		t.Fatalf("F col = %d, want 2", tex.PointOf(off).Col)
	}

	off, ok = findCharOnLine(tex, at, charFindBigT, 'a', 1)
	if !ok || tex.PointOf(off).Col != 3 {
		t.Fatalf("T col = %d, want 3", tex.PointOf(off).Col)
	}
}

func TestOppositeCharFind(t *testing.T) {
	t.Parallel()

	if oppositeCharFind(charFindF) != charFindBigF {
		t.Fatal("F opposite")
	}
	if oppositeCharFind(charFindT) != charFindBigT {
		t.Fatal("T opposite")
	}
}
