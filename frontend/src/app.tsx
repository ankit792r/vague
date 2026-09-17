import { useRef } from "preact/hooks"
import { EmacsFrame } from "./components/frame/EmacsFrame"
import { useEditorInput } from "./hooks/useEditorInput"
import { useEditorView } from "./hooks/useEditorView"

export function App() {
  const editorRef = useRef<HTMLDivElement>(null)
  const view = useEditorView(editorRef)
  const { commandLine } = useEditorInput(view.mode)

  return <EmacsFrame editorRef={editorRef} commandLine={commandLine} {...view} />
}
