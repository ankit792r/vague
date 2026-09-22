package text

import "testing"

func TestUndolistShowsBranches(t *testing.T) {
	t.Parallel()
	tx := New([]byte("a"))
	u := NewUndoTree()
	u.SetInitial(tx.Bytes())

	edit := func(ins string) {
		d := tx.Replace(tx.Len(), tx.Len(), []byte(ins))
		u.Record(Edit{At: d.Start, Removed: d.Removed, Inserted: []byte(ins)})
	}

	edit("b")
	edit("c")
	u.Undo()
	edit("d")

	lines := u.UndolistLines()
	if len(lines) < 2 {
		t.Fatalf("expected branches in undolist, got %v", lines)
	}
	foundBranch := false
	for _, line := range lines {
		if len(line) > 0 && line[len(line)-1] == '*' {
			foundBranch = true
		}
	}
	if !foundBranch {
		t.Fatalf("expected branch marker *, got %v", lines)
	}
}

func TestReplayToSeq(t *testing.T) {
	t.Parallel()
	tx := New([]byte("hi"))
	u := NewUndoTree()
	u.SetInitial(tx.Bytes())
	d := tx.Replace(2, 2, []byte("!"))
	u.Record(Edit{At: d.Start, Removed: d.Removed, Inserted: []byte("!")})

	out, _, ok := u.ReplayToSeq(1)
	if !ok {
		t.Fatal("ReplayToSeq failed")
	}
	if string(out) != "hi!" {
		t.Fatalf("got %q", out)
	}
}
