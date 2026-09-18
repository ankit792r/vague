package text

import "testing"

func applyAll(tx *Text, edits []Edit) {
	for _, e := range edits {
		tx.Replace(e.At, e.OldEnd(), e.Inserted)
	}
}

func edit(tx *Text, u *UndoTree, a, b Offset, ins string) {
	d := tx.Replace(a, b, []byte(ins))
	u.Record(Edit{At: d.Start, Removed: d.Removed, Inserted: []byte(ins)})
}

func TestUndoRedoSingleEdit(t *testing.T) {
	t.Parallel()

	tx := New([]byte("hello"))
	u := NewUndoTree()

	edit(tx, u, 5, 5, " world")
	if got := string(tx.Bytes()); got != "hello world" {
		t.Fatalf("after edit: %q", got)
	}

	edits, _, ok := u.Undo()
	if !ok {
		t.Fatal("Undo returned not ok")
	}
	applyAll(tx, edits)

	if got := string(tx.Bytes()); got != "hello" {
		t.Fatalf("after undo: %q", got)
	}

	edits, _, ok = u.Redo()
	if !ok {
		t.Fatal("Redo returned not ok")
	}
	applyAll(tx, edits)

	if got := string(tx.Bytes()); got != "hello world" {
		t.Fatalf("after redo: %q", got)
	}
}

func TestUndoGroupIsOneStep(t *testing.T) {
	t.Parallel()

	tx := New(nil)
	u := NewUndoTree()

	u.Begin(0)
	edit(tx, u, 0, 0, "a")
	edit(tx, u, 1, 1, "b")
	edit(tx, u, 2, 2, "c")
	u.Commit(3)

	edits, cursor, ok := u.Undo()
	if !ok {
		t.Fatal("Undo failed")
	}
	applyAll(tx, edits)

	if got := string(tx.Bytes()); got != "" {
		t.Fatalf("one undo should remove the whole group, got %q", got)
	}
	if cursor != 0 {
		t.Errorf("cursor = %d, want 0", cursor)
	}
}
