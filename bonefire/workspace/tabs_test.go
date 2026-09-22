package workspace

import "testing"

func TestTabNew(t *testing.T) {
	ws := New()
	frame, err := ws.NewFrame(80, 24, "/tmp")
	if err != nil {
		t.Fatal(err)
	}
	if err := ws.TabNew(frame.ID, ""); err != nil {
		t.Fatal(err)
	}
	if len(frame.Tabs) != 2 {
		t.Fatalf("want 2 tabs, got %d", len(frame.Tabs))
	}
}

func TestArgAddAndList(t *testing.T) {
	ws := New()
	if err := ws.ArgAdd("a.txt", "b.txt"); err != nil {
		t.Fatal(err)
	}
	list := ws.ArgsList()
	if list == "" || list == "argument list empty" {
		t.Fatalf("unexpected list %q", list)
	}
}
