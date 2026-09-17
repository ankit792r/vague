import { useEffect, useState } from "preact/hooks"
import { Fragment } from "preact/jsx-runtime"
import { hostRequest, onHostEvent } from "./host/client"
import type { RedrawPayload } from "./host/protocol"
import { encodeKey } from "./utils/keys"

export function App() {
  const [frameId, setFrameId] = useState<number | null>(null)
  const [bufferText, setBufferText] = useState<string>("")
  const [bufferName, setBufferName] = useState<string>("")

  useEffect(() => {
    const unsubscribe = onHostEvent("redraw", (payload) => {
      const redraw = payload as RedrawPayload
      setFrameId(redraw.frame_id)
      setBufferName(redraw.buffer.name)
      setBufferText(redraw.buffer.text)
    })

    void hostRequest("ready", { height: 100, width: 100 })
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
  }, [])

  return (
    <Fragment>
      <h1>Vague {frameId ?? "…"}</h1>
      <p>{bufferName || "loading buffer…"}</p>
      <pre>{bufferText || "waiting for scratch buffer…"}</pre>
    </Fragment>
  )
}
