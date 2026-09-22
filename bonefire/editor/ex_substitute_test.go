package editor

import (
	"testing"

	frame "vague/bonefire/frame"
	"vague/bonefire/text"
	"vague/bonefire/window"
)

func TestParseSubstituteLine(t *testing.T) {
	cmd, rng, err := parseExLine("s/foo/bar/g")
	if err != nil {
		t.Fatal(err)
	}
	if cmd.kind != "substitute" || cmd.sub == nil {
		t.Fatalf("kind = %q", cmd.kind)
	}
	if cmd.sub.pattern != "foo" || cmd.sub.replacement != "bar" || !cmd.sub.flags.global {
		t.Fatalf("sub = %+v", cmd.sub)
	}
	if rng.wholeBuf || rng.startLine != -1 {
		t.Fatalf("range = %+v", rng)
	}
}

func TestSubstituteOnLine(t *testing.T) {
	ed := NewEditor()
	buf := ed.Scratch("ex-sub-test")
	buf.Text.SetBytes([]byte("foo foo\n"))
	win := &window.Window{Cursor: buf.Text.AddMarker(0, text.GravityRight)}
	fm := &frame.Frame{Width: 80, Height: 24, Dirty: true}

	sub, err := parseSubstitute("s/foo/bar/g")
	if err != nil {
		t.Fatal(err)
	}
	if err := ed.exSubstitute(fm, win, buf, sub, exRange{}); err != nil {
		t.Fatal(err)
	}
	if string(buf.Text.Bytes()) != "bar bar\n" {
		t.Fatalf("got %q", buf.Text.Bytes())
	}
}

func TestSubstitutePercentRange(t *testing.T) {
	ed := NewEditor()
	buf := ed.Scratch("ex-sub-pct")
	buf.Text.SetBytes([]byte("a\nb\n"))
	win := &window.Window{Cursor: buf.Text.AddMarker(0, text.GravityRight)}
	fm := &frame.Frame{Width: 80, Height: 24, Dirty: true}

	sub, _ := parseSubstitute("s/a/x/g")
	rng := exRange{wholeBuf: true}
	if err := ed.exSubstitute(fm, win, buf, sub, rng); err != nil {
		t.Fatal(err)
	}
	if string(buf.Text.Bytes()) != "x\nb\n" {
		t.Fatalf("got %q", buf.Text.Bytes())
	}
}
