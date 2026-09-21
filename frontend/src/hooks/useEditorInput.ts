import { useEffect, useRef, useState } from "preact/hooks"
import { hostRequest } from "../host/client"
import {
  initialCommandLineState,
  isSearchPrompt,
  type CommandLineState,
  type PromptKind,
} from "../types/command"
import { parseCommandLine } from "../utils/command"
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
  return { active: true, kind, text: "", error: null }
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

    const cancelCommand = () => {
      syncCommand(initialCommandLineState())
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

        if (keys === "<CR>") {
          if (isSearchPrompt(cmd.kind)) {
            void hostRequest("execute", {
              name: "search",
              args: [cmd.text],
              bang: cmd.kind === "search-backward",
            })
              .then(() => {
                cancelCommand()
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
              cancelCommand()
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
          syncCommand({
            ...cmd,
            text: cmd.text.slice(0, -1),
            error: null,
          })
          return
        }

        const ch = commandCharFromEvent(e)
        if (ch !== null) {
          syncCommand({
            ...cmd,
            text: cmd.text + ch,
            error: null,
          })
        }

        return
      }

      if (modeRef.current === "normal") {
        if (keys === ":") {
          syncCommand(openPrompt("command"))
          return
        }
        if (keys === "/") {
          syncCommand(openPrompt("search-forward"))
          return
        }
        if (keys === "?") {
          syncCommand(openPrompt("search-backward"))
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
