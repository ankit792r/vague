import { useEffect, useRef, useState } from "preact/hooks"
import { Fragment } from "preact/jsx-runtime"
import { hostRequest, onHostEvent } from "./host/client"
import type { RedrawPayload } from "./host/protocol"
import { encodeKey } from "./utils/keys"
import { measureEditor } from "./utils/measure"

const DUMMY_COMMAND = "open-file"

function renderLine(
  line: string,
  row: number,
  cursor: RedrawPayload["cursor"],
) {
  const display = line === "" ? "\u00a0" : line

  if (!cursor?.visible || cursor.row !== row) {
    return display
  }

  const before = display.slice(0, cursor.column)
  const at = display[cursor.column] ?? "\u00a0"
  const after = display.slice(cursor.column + 1)

  return (
    <Fragment>
      {before}
      <span class="cursor-cell">{at}</span>
      {after}
    </Fragment>
  )
}

export function App() {
  const editorRef = useRef<HTMLDivElement>(null)
  const [lines, setLines] = useState<string[]>([])
  const [bufferName, setBufferName] = useState("*scratch*")
  const [mode, setMode] = useState("normal")
  const [cursor, setCursor] = useState<RedrawPayload["cursor"]>({
    row: 0,
    column: 0,
    visible: true,
  })

  useEffect(() => {
    const unsubscribe = onHostEvent("redraw", (payload) => {
      const redraw = payload as Partial<RedrawPayload>
      setBufferName(redraw.buffer?.name ?? "*scratch*")
      setLines(Array.isArray(redraw.lines) ? redraw.lines : [])
      setMode(redraw.mode ?? "normal")
      setCursor(
        redraw.cursor ?? { row: 0, column: 0, visible: false },
      )
    })

    let readyTimer: ReturnType<typeof setTimeout> | undefined

    const sendReady = () => {
      const el = editorRef.current
      if (!el) {
        return
      }

      clearTimeout(readyTimer)
      readyTimer = setTimeout(() => {
        void document.fonts.ready.then(() => {
          const el = editorRef.current
          if (!el) {
            return
          }

          const { rows, cols } = measureEditor(el)
          void hostRequest("ready", { height: rows, width: cols }).catch((err) => {
            console.error("ready failed", err)
          })
        })
      }, 50)
    }

    sendReady()

    const observer = new ResizeObserver(sendReady)
    if (editorRef.current) {
      observer.observe(editorRef.current)
    }

    return () => {
      clearTimeout(readyTimer)
      unsubscribe()
      observer.disconnect()
    }
  }, [])

  useEffect(() => {
    const onKeyDown = (e: KeyboardEvent) => {
      e.preventDefault()
      const keys = encodeKey(e)
      if (!keys) {
        return
      }

      void hostRequest("input", { keys }).catch((err: unknown) => {
        console.error("input failed", err)
      })
    }

    window.addEventListener("keydown", onKeyDown)
    return () => window.removeEventListener("keydown", onKeyDown)
  }, [])

  return (
    <div class="emacs-frame">
      <div ref={editorRef} class="editor-area" aria-label="editor">
        {(lines ?? []).map((line, index) => (
          <div key={index} class="editor-line">
            {renderLine(line, index, cursor)}
          </div>
        ))}
      </div>

      <div class="status-line" aria-label="status line">
        <span class="status-left">--**- {bufferName}</span>
        <span class="status-right">({mode})</span>
      </div>

      <div class="command-line" aria-label="command line">
        <span class="command-prompt">:</span>
        <span class="command-text">{DUMMY_COMMAND}</span>
      </div>
    </div>
  )
}
