package editor

import "testing"

func TestParsePercentSubstitute(t *testing.T) {
	cmd, rng, err := parseExLine("%s/a/b/g")
	if err != nil {
		t.Fatal(err)
	}
	if !rng.wholeBuf || cmd.sub.pattern != "a" {
		t.Fatalf("cmd=%+v rng=%+v", cmd, rng)
	}
}

func TestParseLineRangeSubstitute(t *testing.T) {
	cmd, rng, err := parseExLine("2,3s/x/y/g")
	if err != nil {
		t.Fatal(err)
	}
	if rng.startLine != 1 || rng.endLine != 2 {
		t.Fatalf("rng=%+v", rng)
	}
	if cmd.sub.replacement != "y" {
		t.Fatalf("sub=%+v", cmd.sub)
	}
}

func TestParseVisualRangeSubstitute(t *testing.T) {
	_, rng, err := parseExLine("'<,'>s/a/b/g")
	if err != nil {
		t.Fatal(err)
	}
	if rng.startLine != -2 {
		t.Fatalf("rng=%+v", rng)
	}
}
