package buffer

import (
	"bytes"

	"vague/bonefire/text"
)

// Replace substitutes ins for the bytes in [a, b) and records the change.
func (b *Buffer) Replace(a, end text.Offset, ins []byte) (text.Delta, error) {
	if b.ReadOnly {
		return text.Delta{}, ErrReadOnly
	}

	delta := b.Text.Replace(a, end, ins)

	b.History.Record(text.Edit{
		At:       delta.Start,
		Removed:  delta.Removed,
		Inserted: bytes.Clone(ins),
	})

	return delta, nil
}

// Insert adds text at an offset.
func (b *Buffer) Insert(at text.Offset, ins []byte) (text.Delta, error) {
	return b.Replace(at, at, ins)
}

// Delete removes the bytes in [a, b).
func (b *Buffer) Delete(a, end text.Offset) (text.Delta, error) {
	return b.Replace(a, end, nil)
}

// BeginEdit opens an undo group.
func (b *Buffer) BeginEdit(cursor text.Offset) {
	b.History.Begin(cursor)
}

// EndEdit closes the innermost undo group.
func (b *Buffer) EndEdit(cursor text.Offset) {
	b.History.Commit(cursor)
}

// Undo reverts the last change and returns where the cursor should go.
func (b *Buffer) Undo() (text.Offset, bool) {
	edits, cursor, ok := b.History.Undo()
	if !ok {
		return 0, false
	}

	b.applyRaw(edits)
	return cursor, true
}

// Redo reapplies the last undone change.
func (b *Buffer) Redo() (text.Offset, bool) {
	edits, cursor, ok := b.History.Redo()
	if !ok {
		return 0, false
	}

	b.applyRaw(edits)
	return cursor, true
}

func (b *Buffer) applyRaw(edits []text.Edit) {
	for _, edit := range edits {
		b.Text.Replace(edit.At, edit.OldEnd(), edit.Inserted)
	}
}
