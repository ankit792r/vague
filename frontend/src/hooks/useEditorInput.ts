import { useEffect, useRef, useState } from "preact/hooks"
import { hostRequest } from "../host/client"
import {
  initialCommandLineState,
  isSearchPrompt,
  type CommandLineState,
  type PromptKind,
} from "../types/command"
import { parseCommandLine } from "../utils/command"
import { historyEntry, pushHistory } from "../utils/commandHistory"
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
  return { active: true, kind, text: "", error: null, historyIndex: null }
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

        if (keys === "<Up>") {
          const index = cmd.historyIndex === null ? 0 : cmd.historyIndex + 1
          const entry = historyEntry(cmd.kind, index)
          if (entry === null) {
            return
          }
          syncCommand({
            ...cmd,
            text: entry,
            error: null,
            historyIndex: index,
          })
          if (isSearchPrompt(cmd.kind)) {
            void previewSearch(cmd.kind, entry)
          }
          return
        }

        if (keys === "<Down>") {
          if (cmd.historyIndex === null || cmd.historyIndex === 0) {
            syncCommand({ ...cmd, text: "", error: null, historyIndex: null })
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
          syncCommand({
            ...cmd,
            text: entry,
            error: null,
            historyIndex: index,
          })
          if (isSearchPrompt(cmd.kind)) {
            void previewSearch(cmd.kind, entry)
          }
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
          const text = cmd.text.slice(0, -1)
          syncCommand({ ...cmd, text, error: null })
          if (isSearchPrompt(cmd.kind)) {
            void previewSearch(cmd.kind, text)
          }
          return
        }

        const ch = commandCharFromEvent(e)
        if (ch !== null) {
          const text = cmd.text + ch
          syncCommand({ ...cmd, text, error: null })
          if (isSearchPrompt(cmd.kind)) {
            void previewSearch(cmd.kind, text)
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
