import type { RefObject } from "preact"
import { CommandLine } from "../chrome/CommandLine"
import { StatusLine } from "../chrome/StatusLine"
import { EditorArea } from "../editor/EditorArea"
import type { EditorViewState } from "../../types/editor"

type EmacsFrameProps = EditorViewState & {
  editorRef: RefObject<HTMLDivElement>
}

export function EmacsFrame({
  editorRef,
  lines,
  bufferName,
  mode,
  cursor,
}: EmacsFrameProps) {
  return (
    <div class="emacs-frame">
      <EditorArea editorRef={editorRef} lines={lines} cursor={cursor} />
      <StatusLine bufferName={bufferName} mode={mode} />
      <CommandLine />
    </div>
  )
}
