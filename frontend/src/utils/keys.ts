export function encodeKey(e: KeyboardEvent): string {
  if (e.key === "Shift" || e.key === "Control" || e.key === "Alt" || e.key === "Meta") {
    return ""
  }

  const named: Record<string, string> = {
    Escape: "<Esc>",
    Enter: "<CR>",
    Tab: "<Tab>",
    Backspace: "<BS>",
    Delete: "<Del>",
    ArrowUp: "<Up>",
    ArrowDown: "<Down>",
    ArrowLeft: "<Left>",
    ArrowRight: "<Right>",
    Home: "<Home>",
    End: "<End>",
    PageUp: "<PageUp>",
    PageDown: "<PageDown>",
    Insert: "<Insert>",
    " ": "<Space>",
  }

  let base = named[e.key]
  if (base === undefined) {
    if (e.key.length === 1) {
      base = e.key
    } else if (e.code.startsWith("Key") && e.code.length === 4) {
      // Layout-independent A–Z when e.key is not a single character (some IME/locale paths).
      base = e.code.slice(3).toLowerCase()
    } else if (e.code.startsWith("Digit") && e.code.length === 6) {
      base = e.code.slice(5)
    } else {
      base = `<${e.key}>`
    }
  }

  if (e.ctrlKey) {
    return `<C-${stripBrackets(base)}>`
  }
  if (e.altKey) {
    return `<A-${stripBrackets(base)}>`
  }
  return base
}

function stripBrackets(s: string): string {
  if (s.startsWith("<") && s.endsWith(">")) {
    return s.slice(1, -1)
  }
  return s
}
