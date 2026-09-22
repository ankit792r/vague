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
}

const SET_OPTIONS: Record<string, ParsedExecute["name"]> = {
  number: "number",
  nu: "number",
  nonumber: "nonumber",
  nonu: "nonumber",
  wrap: "wrap",
  nowrap: "nowrap",
}

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
    if (mapped) {
      return { name: mapped }
    }
    return { name: "set", args: rest ? [rest] : undefined }
  }

  name = ALIASES[name] ?? name
  if (!name) {
    return null
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
