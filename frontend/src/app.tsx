import { useRef } from "preact/hooks"
import { EmacsFrame } from "./components/frame/EmacsFrame"
import { useEditorInput } from "./hooks/useEditorInput"
import { useEditorView } from "./hooks/useEditorView"

export function App() {
  const editorRef = useRef<HTMLDivElement>(null)
  const view = useEditorView(editorRef)

  useEditorInput()

  return <EmacsFrame editorRef={editorRef} {...view} />
}
