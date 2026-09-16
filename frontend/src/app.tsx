import { useEffect, useState } from "preact/hooks"
import { Fragment } from "preact/jsx-runtime"
import { hostRequest } from "./host/client"
import { encodeKey } from "./utils/keys"

export function App() {
  const [frameId, setFrameId] = useState<number | null>(null)

  useEffect(() => {
    void hostRequest("ready", { height: 100, width: 100 })
      .then(({ frame_id, session_id }) => {
        console.log("host ready", session_id, frame_id)
        setFrameId(frame_id)
      })
      .catch((err) => console.error("ready failed", err))
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
      <h1>Hii there {frameId}</h1>
      <p>this is an example buffer</p>
    </Fragment>
  )
}
