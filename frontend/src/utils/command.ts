import type { HostRequestMap } from "../host/protocol"

export type ParsedExecute = HostRequestMap["execute"]

const ALIASES: Record<string, string> = {
  e: "edit",
  w: "write",
  q: "quit",
  wq: "wq",
  x: "x",
  bn: "bnext",
  bp: "bprev",
  b: "buffer",
  ls: "buffers",
  go: "goto",
  help: "help",
}

const SET_OPTIONS: Record<string, ParsedExecute["name"]> = {
  number: "number",
  nu: "number",
  nonumber: "nonumber",
  nonu: "nonumber",
  wrap: "wrap",
  nowrap: "nowrap",
  ignorecase: "set",
  ic: "set",
  noignorecase: "set",
  noic: "set",
  smartcase: "set",
  scs: "set",
  nosmartcase: "set",
  noscs: "set",
  hlsearch: "set",
  hls: "set",
  nohlsearch: "set",
  nohls: "set",
  incsearch: "set",
  noincsearch: "set",
  wrapscan: "set",
  ws: "set",
  nowrapscan: "set",
  nows: "set",
}

const SIMPLE_COMMANDS = new Set([
  "edit",
  "e",
  "write",
  "w",
  "quit",
  "q",
  "wq",
  "x",
  "bnext",
  "bn",
  "bprev",
  "bp",
  "buffer",
  "b",
  "buffers",
  "ls",
  "goto",
  "go",
  "help",
  "search",
  "search_begin",
  "search_preview",
  "search_cancel",
  "complete",
  "register_get",
  "wrap",
  "nowrap",
  "number",
  "nonumber",
  "set",
  "only",
])

export function parseCommandLine(input: string): ParsedExecute | null {
  let line = input.trim()
  if (line.startsWith(":")) {
    line = line.slice(1).trim()
  }

  if (!line) {
    return null
  }

  if (/^\d+$/.test(line)) {
    return { name: "goto", count: Number(line) }
  }

  const looksLikeEx =
    !/^set(\s|$)/i.test(line) &&
    (/^[%0-9.'$+\-,<>\s]*s[/\w]/.test(line) ||
      /^[gv]/.test(line) ||
      /^norm/.test(line) ||
      /^nohl/.test(line) ||
      /^sort\b/.test(line) ||
      /^uniq\b/.test(line) ||
      /^vimgrep\b/.test(line) ||
      /^[lt]?vi[m]?\b/.test(line) ||
      /^(ta|tag|tags)\b/.test(line))

  if (looksLikeEx) {
    return { name: "ex", args: [line] }
  }

  const space = line.indexOf(" ")
  let name = space === -1 ? line : line.slice(0, space)
  const rest = space === -1 ? "" : line.slice(space + 1).trim()

  let bang = false
  if (name.endsWith("!")) {
    bang = true
    name = name.slice(0, -1)
  }

  if (name === "set") {
    const opt = rest.toLowerCase()
    const mapped = SET_OPTIONS[opt]
    if (mapped === "set") {
      return { name: "set", args: rest ? [rest] : undefined }
    }
    if (mapped) {
      return { name: mapped }
    }
    return { name: "set", args: rest ? [rest] : undefined }
  }

  name = ALIASES[name] ?? name
  if (!name) {
    return null
  }

  if (!SIMPLE_COMMANDS.has(name)) {
    return { name: "ex", args: [line] }
  }

  const args = rest ? [rest] : undefined

  return {
    name,
    args,
    bang: bang || undefined,
    count:
      name === "goto" && rest && /^\d+$/.test(rest) ? Number(rest) : undefined,
  }
}
