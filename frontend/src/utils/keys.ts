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
    " ": "<Space>",
  }

  const base = named[e.key] ?? (e.key.length === 1 ? e.key : `<${e.key}>`)

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
