package editor

import (
	"testing"
)

func TestPatternMagicModes(t *testing.T) {
	t.Parallel()

	ed := NewEditor()
	data := []byte("foo\nBar\nbaz")

	pp, err := ed.parseSearchPattern(`\Vfoo`)
	if err != nil {
		t.Fatal(err)
	}
	start, end, ok := findPatternForward(data, pp, -1, true)
	if !ok || start != 0 || end != 3 {
		t.Fatalf("literal forward = %d..%d ok=%v", start, end, ok)
	}

	pp, err = ed.parseSearchPattern(`\vB.r`)
	if err != nil {
		t.Fatal(err)
	}
	start, _, ok = findPatternForward(data, pp, -1, true)
	if !ok || start != 4 {
		t.Fatalf("very magic . = %d", start)
	}

	pp, err = ed.parseSearchPattern(`\M^Bar`)
	if err != nil {
		t.Fatal(err)
	}
	start, _, ok = findPatternForward(data, pp, -1, true)
	if !ok || start != 4 {
		t.Fatalf("^ anchor = %d", start)
	}
}

func TestPatternIgnoreCaseAndSmartCase(t *testing.T) {
	t.Parallel()

	ed := NewEditor()
	ed.searchOpts.IgnoreCase = true
	data := []byte("Foo bar")

	pp, err := ed.parseSearchPattern("foo")
	if err != nil {
		t.Fatal(err)
	}
	if !pp.ignoreCase {
		t.Fatal("expected ignorecase")
	}
	_, _, ok := findPatternForward(data, pp, -1, true)
	if !ok {
		t.Fatal("expected match")
	}

	ed.searchOpts.SmartCase = true
	pp, err = ed.parseSearchPattern("Foo")
	if err != nil {
		t.Fatal(err)
	}
	if pp.ignoreCase {
		t.Fatal("smartcase should disable ic for upper pattern")
	}
}

func TestFindAllPattern(t *testing.T) {
	t.Parallel()

	ed := NewEditor()
	data := []byte("aa aa aa")
	pp, err := ed.parseSearchPattern("aa")
	if err != nil {
		t.Fatal(err)
	}
	all := findAllPattern(data, pp)
	if len(all) != 3 {
		t.Fatalf("matches = %d, want 3", len(all))
	}
}
