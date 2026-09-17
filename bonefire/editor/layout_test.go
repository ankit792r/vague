package editor

import (
	"testing"

	"vague/bonefire/text"
)

func TestVisualLinesWrapNoTrailingPadding(t *testing.T) {
	t.Parallel()

	got := visualLines(text.New([]byte("abcdefghijklmn")), 13, 0, true)
	want := []string{"abcdefghijklm", "n"}

	if len(got) != len(want) {
		t.Fatalf("got %d lines, want %d: %v", len(got), len(want), got)
	}

	for i := range want {
		if got[i] != want[i] {
			t.Errorf("line %d = %q, want %q", i, got[i], want[i])
		}
	}
}

func TestVisualLinesNoWrapTruncates(t *testing.T) {
	t.Parallel()

	got := visualLines(text.New([]byte("abcdefghijklmn")), 13, 0, false)
	want := []string{"abcdefghijklm"}

	if len(got) != 1 || got[0] != want[0] {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestVisualLinesRespectsMaxRows(t *testing.T) {
	t.Parallel()

	got := visualLines(text.New([]byte("abcdefghijklmnop")), 4, 2, true)
	want := []string{"abcd", "efgh"}

	if len(got) != len(want) {
		t.Fatalf("got %v, want %v", got, want)
	}

	for i := range want {
		if got[i] != want[i] {
			t.Errorf("line %d = %q, want %q", i, got[i], want[i])
		}
	}
}

func TestVisualLinesMultipleLogicalLines(t *testing.T) {
	t.Parallel()

	got := visualLines(text.New([]byte("ab\ncdef")), 3, 0, true)
	want := []string{"ab", "cde", "f"}

	if len(got) != len(want) {
		t.Fatalf("got %v, want %v", got, want)
	}

	for i := range want {
		if got[i] != want[i] {
			t.Errorf("line %d = %q, want %q", i, got[i], want[i])
		}
	}
}
