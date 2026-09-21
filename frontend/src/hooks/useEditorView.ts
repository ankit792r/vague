import { useEffect, useState } from "preact/hooks"
import type { RefObject } from "preact"
import { hostRequest, onHostEvent } from "../host/client"
import type { RedrawPayload } from "../host/protocol"
import {
  editorViewFromRedraw,
  initialEditorViewState,
  type EditorViewState,
} from "../types/editor"
import { measureEditor } from "../utils/measure"

export function useEditorView(
  editorRef: RefObject<HTMLDivElement>,
): EditorViewState {
  const [view, setView] = useState(initialEditorViewState)

  useEffect(() => {
    const unsubscribe = onHostEvent("redraw", (payload) => {
      setView(editorViewFromRedraw(payload as Partial<RedrawPayload>))
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
          void hostRequest("ui_ready", { height: rows, width: cols }).catch((err) => {
            console.error("ui_ready failed", err)
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
  }, [editorRef])

  return view
}
