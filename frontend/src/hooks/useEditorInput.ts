import { useEffect, useRef, useState } from "preact/hooks"
import { hostRequest } from "../host/client"
import {
  initialCommandLineState,
  isSearchPrompt,
  withText,
  type CommandLineState,
  type PromptKind,
} from "../types/command"
import { parseCommandLine } from "../utils/command"
import { historyEntry, pushHistory } from "../utils/commandHistory"
import {
  applyCompletion,
  completionContext,
  type CompletionContext,
} from "../utils/completion"
import { encodeKey } from "../utils/keys"

function commandCharFromEvent(e: KeyboardEvent): string | null {
  if (e.key === " ") {
    return " "
  }

  if (e.key.length === 1 && !e.ctrlKey && !e.altKey && !e.metaKey) {
    return e.key
  }

  return null
}

function executeErrorMessage(err: unknown): string {
  if (err instanceof Error && err.message) {
    return err.message
  }
  return "Command failed"
}

function openPrompt(kind: PromptKind): CommandLineState {
  return {
    ...initialCommandLineState(),
    active: true,
    kind,
  }
}

type CompleteResult = { candidates?: string[] }

async function fetchCompletions(ctx: CompletionContext): Promise<string[]> {
  const result = (await hostRequest("execute", {
    name: "complete",
    args: [ctx.kind, ctx.prefix],
  })) as CompleteResult
  return result.candidates ?? []
}

async function previewSearch(kind: PromptKind, text: string) {
  if (!isSearchPrompt(kind)) {
    return
  }
  await hostRequest("execute", {
    name: "search_preview",
    args: [text],
    bang: kind === "search-backward",
  })
}

async function beginSearch(kind: PromptKind) {
  await hostRequest("execute", {
    name: "search_begin",
    bang: kind === "search-backward",
  })
}

async function cancelSearchPreview() {
  await hostRequest("execute", { name: "search_cancel" })
}

async function insertRegisterAtCursor(cmd: CommandLineState): Promise<CommandLineState> {
  const result = (await hostRequest("execute", { name: "register_get" })) as {
    text?: string
  }
  const insert = result.text ?? ""
  if (!insert) {
    return { ...cmd, pendingRegister: false }
  }
  const before = cmd.text.slice(0, cmd.cursor)
  const after = cmd.text.slice(cmd.cursor)
  const text = before + insert + after
  const next = withText(cmd, text, cmd.cursor + insert.length)
  if (isSearchPrompt(cmd.kind)) {
    void previewSearch(cmd.kind, next.text)
  }
  return next
}

export function useEditorInput(editorMode: string) {
  const [commandLine, setCommandLine] = useState<CommandLineState>(
    initialCommandLineState(),
  )
  const modeRef = useRef(editorMode)
  const commandRef = useRef(commandLine)

  modeRef.current = editorMode
  commandRef.current = commandLine

  useEffect(() => {
    const syncCommand = (next: CommandLineState) => {
      commandRef.current = next
      setCommandLine(next)
    }

    const closeCommand = () => {
      syncCommand(initialCommandLineState())
    }

    const cancelCommand = () => {
      const cmd = commandRef.current
      if (isSearchPrompt(cmd.kind)) {
        void cancelSearchPreview().catch((err: unknown) => {
          console.error("search cancel failed", err)
        })
      }
      closeCommand()
    }

    const runTabCompletion = async (cmd: CommandLineState) => {
      const ctx = completionContext(cmd.kind, cmd.text, cmd.cursor)
      if (!ctx) {
        return
      }

      let candidates = cmd.completions
      let index = cmd.completionIndex
      if (candidates.length === 0) {
        candidates = await fetchCompletions(ctx)
        index = 0
      } else {
        index = (index + 1) % candidates.length
      }

      if (candidates.length === 0) {
        return
      }

      const pick = candidates[index] ?? candidates[0]!
      const applied = applyCompletion(cmd.text, ctx, pick)
      syncCommand({
        ...withText(cmd, applied.text, applied.cursor),
        completions: candidates,
        completionIndex: index,
      })
    }

    const onKeyDown = (e: KeyboardEvent) => {
      e.preventDefault()
      const keys = encodeKey(e)
      if (!keys) {
        return
      }

      const cmd = commandRef.current

      if (cmd.active) {
        if (keys === "<Esc>" || keys === "<C-c>" || keys === "<C-g>") {
          cancelCommand()
          return
        }

        if (keys === "<C-r>") {
          syncCommand({ ...cmd, pendingRegister: true })
          return
        }

        if (cmd.pendingRegister) {
          void insertRegisterAtCursor({ ...cmd, pendingRegister: false })
            .then(syncCommand)
            .catch((err: unknown) => {
              console.error("register insert failed", err)
              syncCommand({ ...cmd, pendingRegister: false })
            })
          return
        }

        if (keys === "<Tab>") {
          void runTabCompletion(cmd).catch((err: unknown) => {
            console.error("completion failed", err)
          })
          return
        }

        if (keys === "<Up>") {
          const index = cmd.historyIndex === null ? 0 : cmd.historyIndex + 1
          const entry = historyEntry(cmd.kind, index)
          if (entry === null) {
            return
          }
          const next = withText(cmd, entry, entry.length)
          syncCommand({ ...next, historyIndex: index })
          if (isSearchPrompt(cmd.kind)) {
            void previewSearch(cmd.kind, entry)
          }
          return
        }

        if (keys === "<Down>") {
          if (cmd.historyIndex === null || cmd.historyIndex === 0) {
            syncCommand(withText(cmd, "", 0))
            if (isSearchPrompt(cmd.kind)) {
              void previewSearch(cmd.kind, "")
            }
            return
          }
          const index = cmd.historyIndex - 1
          const entry = historyEntry(cmd.kind, index)
          if (entry === null) {
            return
          }
          const next = withText(cmd, entry, entry.length)
          syncCommand({ ...next, historyIndex: index })
          if (isSearchPrompt(cmd.kind)) {
            void previewSearch(cmd.kind, entry)
          }
          return
        }

        if (keys === "<Left>") {
          syncCommand({
            ...cmd,
            cursor: Math.max(0, cmd.cursor - 1),
            completions: [],
          })
          return
        }

        if (keys === "<Right>") {
          syncCommand({
            ...cmd,
            cursor: Math.min(cmd.text.length, cmd.cursor + 1),
            completions: [],
          })
          return
        }

        if (keys === "<Home>") {
          syncCommand({ ...cmd, cursor: 0, completions: [] })
          return
        }

        if (keys === "<End>") {
          syncCommand({ ...cmd, cursor: cmd.text.length, completions: [] })
          return
        }

        if (keys === "<CR>") {
          if (isSearchPrompt(cmd.kind)) {
            void hostRequest("execute", {
              name: "search",
              args: [cmd.text],
              bang: cmd.kind === "search-backward",
            })
              .then(() => {
                pushHistory(cmd.kind, cmd.text)
                closeCommand()
              })
              .catch((err: unknown) => {
                syncCommand({
                  ...cmd,
                  error: executeErrorMessage(err),
                })
              })
            return
          }

          const parsed = parseCommandLine(cmd.text)
          if (!parsed) {
            syncCommand({ ...cmd, error: "No command" })
            return
          }

          void hostRequest("execute", parsed)
            .then(() => {
              pushHistory(cmd.kind, cmd.text)
              closeCommand()
            })
            .catch((err: unknown) => {
              syncCommand({
                ...cmd,
                error: executeErrorMessage(err),
              })
            })
          return
        }

        if (keys === "<BS>") {
          if (cmd.cursor === 0) {
            return
          }
          const text =
            cmd.text.slice(0, cmd.cursor - 1) + cmd.text.slice(cmd.cursor)
          const next = withText(cmd, text, cmd.cursor - 1)
          syncCommand(next)
          if (isSearchPrompt(cmd.kind)) {
            void previewSearch(cmd.kind, next.text)
          }
          return
        }

        const ch = commandCharFromEvent(e)
        if (ch !== null) {
          const text =
            cmd.text.slice(0, cmd.cursor) + ch + cmd.text.slice(cmd.cursor)
          const next = withText(cmd, text, cmd.cursor + 1)
          syncCommand(next)
          if (isSearchPrompt(cmd.kind)) {
            void previewSearch(cmd.kind, next.text)
          }
        }

        return
      }

      if (modeRef.current === "normal") {
        if (keys === ":") {
          syncCommand(openPrompt("command"))
          return
        }
        if (keys === "/") {
          const next = openPrompt("search-forward")
          syncCommand(next)
          void beginSearch(next.kind).catch((err: unknown) => {
            console.error("search begin failed", err)
          })
          return
        }
        if (keys === "?") {
          const next = openPrompt("search-backward")
          syncCommand(next)
          void beginSearch(next.kind).catch((err: unknown) => {
            console.error("search begin failed", err)
          })
          return
        }
      }

      void hostRequest("input", { keys }).catch((err: unknown) => {
        console.error("input failed", err)
      })
    }

    window.addEventListener("keydown", onKeyDown)
    return () => window.removeEventListener("keydown", onKeyDown)
  }, [])

  return { commandLine }
}
