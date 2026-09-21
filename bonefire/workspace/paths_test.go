package workspace

import (
	"os"
	"path/filepath"
	"testing"
)

func TestInitialWorkDirUsesHomeWhenClientDirIsRoot(t *testing.T) {
	t.Parallel()

	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatal(err)
	}

	got := initialWorkDir("/")
	if got != home {
		t.Fatalf("got %q, want %q", got, home)
	}
}

func TestResolvePathAgainstWorkDir(t *testing.T) {
	t.Parallel()

	base := t.TempDir()
	got, err := resolvePathAgainst(base, "notes.txt")
	if err != nil {
		t.Fatal(err)
	}

	want := filepath.Join(base, "notes.txt")
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}
