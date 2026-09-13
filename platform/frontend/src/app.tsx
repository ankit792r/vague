import { useEffect } from "preact/hooks"
import { callNative, encodeKey, nextCallId } from "./native.ts"

export function App() {
  useEffect(() => {
    const onKeyDown = (e: KeyboardEvent) => {
      const keys = encodeKey(e)
      console.log("keys", keys)
      if (!keys) {
        return
      }

      e.preventDefault()
      void callNative(nextCallId(), "input", { keys }).catch((err: unknown) => {
        console.error("callNative failed", err)
      })
    }

    window.addEventListener("keydown", onKeyDown)
    return () => window.removeEventListener("keydown", onKeyDown)
  }, [])

  return <div>Hello Vague</div>
}
