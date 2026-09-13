import { useEffect } from "preact/hooks"
import { Fragment } from "preact/jsx-runtime"
import { callIpc, nextCallId } from "./utils/ipc"

export function App() {
  useEffect(() => {
    const onKeyDown = (e: KeyboardEvent) => {
      e.preventDefault()
      void callIpc(nextCallId(), "input", { e })
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
      <h1>Hii there</h1>
      <p>this is an example buffer</p>
    </Fragment>
  )
}
