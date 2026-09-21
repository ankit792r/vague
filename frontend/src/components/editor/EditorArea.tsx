import type { RefObject } from "preact"
import type { EditorCursor, EditorSelection } from "../../types/editor"
import { EditorLine } from "./EditorLine"

type EditorAreaProps = {
  editorRef: RefObject<HTMLDivElement>
  lines: string[]
  cursor: EditorCursor
  selection: EditorSelection | null
  hideCursor?: boolean
}

export function EditorArea({ editorRef, lines, cursor, selection, hideCursor }: EditorAreaProps) {
  return (
    <div ref={editorRef} class="editor-area" aria-label="editor">
      {(lines ?? []).map((line, index) => (
        <div key={index} class="editor-line">
          <EditorLine
            line={line}
            row={index}
            cursor={hideCursor ? { ...cursor, visible: false } : cursor}
            selection={selection}
          />
        </div>
      ))}
    </div>
  )
}
