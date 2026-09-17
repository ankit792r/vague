import { useEffect } from "preact/hooks"
import { hostRequest } from "../host/client"
import { encodeKey } from "../utils/keys"

export function useEditorInput() {
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
}
