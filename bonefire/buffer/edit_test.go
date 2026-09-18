package buffer

import "testing"

func TestUndoClearsModifiedAtSavedState(t *testing.T) {
	t.Parallel()

	buf := NewScratch(1, "test")
	buf.Text.SetBytes([]byte("hi"))

	if _, err := buf.Insert(2, []byte("!")); err != nil {
		t.Fatal(err)
	}

	buf.savedSeq = buf.History.Seq()

	if _, err := buf.Insert(3, []byte("x")); err != nil {
		t.Fatal(err)
	}

	if !buf.Modified() {
		t.Fatal("expected modified after second edit")
	}

	if _, ok := buf.Undo(); !ok {
		t.Fatal("undo failed")
	}

	if buf.Modified() {
		t.Fatal("undo to saved state should clear modified")
	}
}
