import type { PromptKind } from "../types/command"

export type CompletionKind = "command" | "set" | "buffer" | "path"

export type CompletionContext = {
  kind: CompletionKind
  prefix: string
  replaceStart: number
  replaceEnd: number
}

function wordBeforeCursor(text: string, cursor: number): { start: number; end: number; word: string } {
  const end = Math.min(cursor, text.length)
  let start = end
  while (start > 0 && !/\s/.test(text[start - 1] ?? "")) {
    start--
  }
  return { start, end, word: text.slice(start, end) }
}

export function completionContext(
  promptKind: PromptKind,
  text: string,
  cursor: number,
): CompletionContext | null {
  if (promptKind !== "command") {
    return null
  }

  const before = text.slice(0, cursor)
  const trimmed = before.trimStart()

  if (!trimmed.includes(" ")) {
    const { start, end, word } = wordBeforeCursor(text, cursor)
    return {
      kind: "command",
      prefix: word,
      replaceStart: start,
      replaceEnd: end,
    }
  }

  const space = trimmed.indexOf(" ")
  const cmd = trimmed.slice(0, space).toLowerCase()
  const { start, end, word } = wordBeforeCursor(text, cursor)

  switch (cmd) {
    case "set":
      return {
        kind: "set",
        prefix: word,
        replaceStart: start,
        replaceEnd: end,
      }
    case "e":
    case "edit":
      return {
        kind: "path",
        prefix: word,
        replaceStart: start,
        replaceEnd: end,
      }
    case "b":
    case "buffer":
      return {
        kind: "buffer",
        prefix: word,
        replaceStart: start,
        replaceEnd: end,
      }
    default:
      return null
  }
}

export function applyCompletion(
  text: string,
  ctx: CompletionContext,
  candidate: string,
): { text: string; cursor: number } {
  const next =
    text.slice(0, ctx.replaceStart) + candidate + text.slice(ctx.replaceEnd)
  const cursor = ctx.replaceStart + candidate.length
  return { text: next, cursor }
}
