import { useEffect, useRef, useState } from "preact/hooks"
import { hostRequest, onHostEvent } from "./host/client"
import type { RedrawPayload } from "./host/protocol"
import { encodeKey } from "./utils/keys"

const DUMMY_MODE = "Normal"
const DUMMY_COMMAND = "open-file"

function measureEditor(el: HTMLElement) {
  const style = getComputedStyle(el)
  const lineHeight = parseFloat(style.lineHeight) || parseFloat(style.fontSize) * 1.4
  const fontSize = parseFloat(style.fontSize) || 16
  const charWidth = fontSize * 0.6

  return {
    rows: Math.max(1, Math.floor(el.clientHeight / lineHeight)),
    cols: Math.max(1, Math.floor(el.clientWidth / charWidth)),
  }
}

export function App() {
  const editorRef = useRef<HTMLDivElement>(null)
  const [lines, setLines] = useState<string[]>([])
  const [bufferName, setBufferName] = useState("*scratch*")

  useEffect(() => {
    const unsubscribe = onHostEvent("redraw", (payload) => {
      const redraw = payload as RedrawPayload
      setBufferName(redraw.buffer.name)
      setLines(redraw.lines)
    })

    const sendReady = () => {
      const el = editorRef.current
      if (!el) {
        return
      }

      const { rows, cols } = measureEditor(el)
      void hostRequest("ready", { height: rows, width: cols }).catch((err) => {
        console.error("ready failed", err)
      })
    }

    sendReady()

    const observer = new ResizeObserver(sendReady)
    if (editorRef.current) {
      observer.observe(editorRef.current)
    }

    return () => {
      unsubscribe()
      observer.disconnect()
    }
  }, [])

  useEffect(() => {
    const onKeyDown = (e: KeyboardEvent) => {
      e.preventDefault()
      void hostRequest("input", { keys: encodeKey(e) }).catch((err: unknown) => {
        console.error("input failed", err)
      })
    }

    window.addEventListener("keydown", onKeyDown)
    return () => window.removeEventListener("keydown", onKeyDown)
  }, [])

  return (
    <div class="emacs-frame">
      <div ref={editorRef} class="editor-area" aria-label="editor">
        {lines.length > 0
          ? lines.map((line, index) => (
              <div key={index} class="editor-line">
                {line === "" ? "\u00a0" : line}
              </div>
            ))
          : null}
      </div>

      <div class="status-line" aria-label="status line">
        <span class="status-left">--**- {bufferName}</span>
        <span class="status-right">({DUMMY_MODE})</span>
      </div>

      <div class="command-line" aria-label="command line">
        <span class="command-prompt">:</span>
        <span class="command-text">{DUMMY_COMMAND}</span>
      </div>
    </div>
  )
}
