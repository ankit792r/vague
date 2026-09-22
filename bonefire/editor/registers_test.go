package editor

import "testing"

func TestNamedRegisterAppend(t *testing.T) {
	t.Parallel()

	ed := NewEditor()
	ed.activeReg = registerID{name: "a"}
	ed.recordYank([]byte("foo"), false)
	ed.activeReg = registerID{name: "a", append: true}
	ed.recordYank([]byte("bar"), false)

	text, _, ok := ed.readRegister(registerID{name: "a"})
	if !ok || string(text) != "foobar" {
		t.Fatalf("append = %q ok=%v", text, ok)
	}
}

func TestDeleteRing(t *testing.T) {
	t.Parallel()

	ed := NewEditor()
	ed.recordDelete([]byte("a"), false)
	ed.recordDelete([]byte("b"), false)

	text, _, ok := ed.readRegister(registerID{name: "1"})
	if !ok || string(text) != "b" {
		t.Fatalf("reg1 = %q", text)
	}
	text, _, ok = ed.readRegister(registerID{name: "2"})
	if !ok || string(text) != "a" {
		t.Fatalf("reg2 = %q", text)
	}
}

func TestYankZeroRegister(t *testing.T) {
	t.Parallel()

	ed := NewEditor()
	ed.recordYank([]byte("yank"), true)
	text, _, ok := ed.readRegister(registerID{name: "0"})
	if !ok || string(text) != "yank" {
		t.Fatalf("reg0 = %q", text)
	}
}

func TestBlackHoleRegister(t *testing.T) {
	t.Parallel()

	ed := NewEditor()
	ed.activeReg = registerID{name: "_"}
	ed.recordYank([]byte("gone"), false)
	if !ed.regs.unnamed.empty() {
		t.Fatal("black hole should not fill unnamed")
	}
}

func TestSlashAndColonRegisters(t *testing.T) {
	t.Parallel()

	ed := NewEditor()
	ed.SetLastSearch("pat")
	ed.SetLastCommand("set wrap")
	s, ok := ed.GetRegister("/")
	if !ok || s != "pat" {
		t.Fatalf("/ = %q", s)
	}
	s, ok = ed.GetRegister(":")
	if !ok || s != "set wrap" {
		t.Fatalf(": = %q", s)
	}
}

func TestGetSetRegisterAPI(t *testing.T) {
	t.Parallel()

	ed := NewEditor()
	if err := ed.SetRegister("b", "hello", false); err != nil {
		t.Fatal(err)
	}
	s, ok := ed.GetRegister("b")
	if !ok || s != "hello" {
		t.Fatalf("got %q", s)
	}
}
