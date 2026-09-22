package clipboard

import "testing"

func TestClipboardReadWrite(t *testing.T) {
	WritePrimary("star")
	WriteClipboard("plus")
	s, ok := ReadPrimary()
	if !ok || s != "star" {
		t.Fatalf("primary = %q", s)
	}
	s, ok = ReadClipboard()
	if !ok || s != "plus" {
		t.Fatalf("clipboard = %q", s)
	}
}

func TestOSC52Stub(t *testing.T) {
	if err := WriteOSC52([]byte("hi")); err != nil {
		t.Fatal(err)
	}
}
