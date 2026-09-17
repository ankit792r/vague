import { useEffect, useRef, useState } from "preact/hooks"
import { hostRequest } from "../host/client"
import {
  initialCommandLineState,
  type CommandLineState,
} from "../types/command"
import { parseCommandLine } from "../utils/command"
import { encodeKey } from "../utils/keys"

function appendCommandChar(keys: string): string | null {
  if (keys === "<Space>") {
    return " "
  }

  if (keys.length === 1) {
    return keys
  }

  return null
}

export function useEditorInput(editorMode: string) {
  const [commandLine, setCommandLine] = useState<CommandLineState>(
    initialCommandLineState,
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
              const message = err instanceof Error ? err.message : String(err)
              syncCommand({ ...cmd, error: message })
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

        const ch = appendCommandChar(keys)
        if (ch !== null) {
          syncCommand({
            ...cmd,
            text: cmd.text + ch,
            error: null,
          })
        }

        return
      }

      if (modeRef.current === "normal" && keys === ":") {
        syncCommand({ active: true, text: "", error: null })
        return
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
