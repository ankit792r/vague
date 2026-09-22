package workspace_test

import (
	"testing"

	"vague/bonefire/workspace"
)

func TestRunExLineSetHlsearch(t *testing.T) {
	t.Parallel()

	ws := workspace.New()
	ws.Editor.Scratch("*scratch*")
	frame, err := ws.NewFrame(80, 10, "")
	if err != nil {
		t.Fatal(err)
	}

	msg, err := ws.RunExLine(frame.ID, "set hlsearch")
	if err != nil {
		t.Fatal(err)
	}
	if msg != "hlsearch" {
		t.Fatalf("msg = %q, want hlsearch", msg)
	}
	if !ws.Editor.SearchOptsForTest().HlSearch {
		t.Fatal("expected hlsearch on")
	}

	msg, err = ws.RunExLine(frame.ID, "set nohls")
	if err != nil {
		t.Fatal(err)
	}
	if msg != "nohlsearch" {
		t.Fatalf("msg = %q, want nohlsearch", msg)
	}
}
