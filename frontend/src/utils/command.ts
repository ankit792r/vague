import type { HostRequestMap } from "../host/protocol"

export type ParsedExecute = HostRequestMap["execute"]

const ALIASES: Record<string, string> = {
  e: "edit",
  w: "write",
  q: "quit",
  wq: "wq",
  x: "x",
}

export function parseCommandLine(input: string): ParsedExecute | null {
  let line = input.trim()
  if (line.startsWith(":")) {
    line = line.slice(1).trim()
  }

  if (!line) {
    return null
  }

  const space = line.indexOf(" ")
  let name = space === -1 ? line : line.slice(0, space)
  const rest = space === -1 ? "" : line.slice(space + 1).trim()

  let bang = false
  if (name.endsWith("!")) {
    bang = true
    name = name.slice(0, -1)
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
  }
}
