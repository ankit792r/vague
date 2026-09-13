export type NativeReply = {
  id: number
  result?: unknown
  error?: string
}

type CallNativeFn = (
  id: number,
  method: string,
  params: unknown,
) => Promise<NativeReply>

declare global {
  interface Window {
    callNative?: CallNativeFn
  }
}

let nextId = 1

export function nextCallId(): number {
  return nextId++
}

export function callNative(
  id: number,
  method: string,
  params: unknown = null,
): Promise<NativeReply> {
  const native = window.callNative
  if (typeof native !== "function") {
    return Promise.reject(new Error("native bridge is not available"))
  }
  return native(id, method, params)
}

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
