import type { RefObject } from "preact"
import type { EditorCursor } from "../../types/editor"
import { EditorLine } from "./EditorLine"

type EditorAreaProps = {
  editorRef: RefObject<HTMLDivElement>
  lines: string[]
  cursor: EditorCursor
}

export function EditorArea({ editorRef, lines, cursor }: EditorAreaProps) {
  return (
    <div ref={editorRef} class="editor-area" aria-label="editor">
      {(lines ?? []).map((line, index) => (
        <div key={index} class="editor-line">
          <EditorLine line={line} row={index} cursor={cursor} />
        </div>
      ))}
    </div>
  )
}
