package editor

import "testing"

func TestSplitSearchOffset(t *testing.T) {
	t.Parallel()

	pat, off := splitSearchOffset("foo+2")
	if pat != "foo" || off.delta != 2 || off.anchor != 0 {
		t.Fatalf("foo+2 = %q %+v", pat, off)
	}

	pat, off = splitSearchOffset("whole")
	if pat != "whole" || off != (searchOffset{}) {
		t.Fatalf("whole = %q %+v", pat, off)
	}
}

func TestApplySearchOffset(t *testing.T) {
	t.Parallel()

	pos := applySearchOffset(5, 8, 100, searchOffset{anchor: 'e'})
	if pos != 7 {
		t.Fatalf("e anchor = %d, want 7", pos)
	}

	pos = applySearchOffset(5, 8, 100, searchOffset{delta: 2})
	if pos != 7 {
		t.Fatalf("start+2 = %d, want 7", pos)
	}
}
