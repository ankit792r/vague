import { useEffect, useState } from "preact/hooks"
import { Fragment } from "preact/jsx-runtime"
import { hostRequest, onHostEvent } from "./host/client"
import type { RedrawPayload } from "./host/protocol"
import { encodeKey } from "./utils/keys"

export function App() {
  const [frameId, setFrameId] = useState<number | null>(null)
  const [bufferName, setBufferName] = useState<string>("")
  const [lines, setLines] = useState<string[]>([])
  const [wrap, setWrap] = useState(true)
  const [columns, setColumns] = useState(0)

  useEffect(() => {
    const unsubscribe = onHostEvent("redraw", (payload) => {
      const redraw = payload as RedrawPayload
      setFrameId(redraw.frame_id)
      setBufferName(redraw.buffer.name)
      setLines(redraw.lines)
      setWrap(redraw.wrap)
      setColumns(redraw.columns)
    })

    void hostRequest("ready", { height: 24, width: 13 })
      .then(({ frame_id, session_id }) => {
        console.log("host ready", session_id, frame_id)
        setFrameId(frame_id)
      })
      .catch((err) => console.error("ready failed", err))

    return unsubscribe
  }, [])

  useEffect(() => {
    const onKeyDown = (e: KeyboardEvent) => {
      e.preventDefault()

      if (e.key === "F2") {
        void hostRequest("execute", { name: wrap ? "nowrap" : "wrap" })
        return
      }

      void hostRequest("input", { keys: encodeKey(e) })
        .then((reply) => {
          console.log(reply)
        })
        .catch((err: unknown) => {
          console.error("callIpc failed", err)
        })
    }

    window.addEventListener("keydown", onKeyDown)
    return () => window.removeEventListener("keydown", onKeyDown)
  }, [wrap])

  return (
    <Fragment>
      <div class="editor-view" aria-label="editor">
        {lines.length > 0
          ? lines.map((line, index) => (
              <div key={index} class="editor-line">
                {line === "" ? "\u00a0" : line}
              </div>
            ))
          : "waiting for redraw…"}
      </div>
    </Fragment>
  )
}
