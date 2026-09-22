import type { PromptKind } from "../types/command"

export type HistoryBucket = "command" | "search"

const MAX_HISTORY = 100

const store: Record<HistoryBucket, string[]> = {
  command: [],
  search: [],
}

export function historyBucket(kind: PromptKind): HistoryBucket {
  return kind === "command" ? "command" : "search"
}

export function historyFor(kind: PromptKind): string[] {
  return store[historyBucket(kind)]
}

export function pushHistory(kind: PromptKind, line: string) {
  const trimmed = line.trim()
  if (!trimmed) {
    return
  }
  const bucket = historyBucket(kind)
  const list = store[bucket]
  if (list[list.length - 1] === trimmed) {
    return
  }
  const filtered = list.filter((entry) => entry !== trimmed)
  filtered.push(trimmed)
  if (filtered.length > MAX_HISTORY) {
    filtered.splice(0, filtered.length - MAX_HISTORY)
  }
  store[bucket] = filtered
}

export function historyEntry(kind: PromptKind, index: number): string | null {
  const list = historyFor(kind)
  if (index < 0 || index >= list.length) {
    return null
  }
  return list[list.length - 1 - index]
}
