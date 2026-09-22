import type { RefObject } from "preact"
import type { EditorCursor, EditorSearchMatch, EditorSelection } from "../../types/editor"
import { EditorLine } from "./EditorLine"

function cursorShapeForMode(mode: string): "block" | "bar" {
  return mode === "insert" ? "bar" : "block"
}

type EditorAreaProps = {
  editorRef: RefObject<HTMLDivElement>
  lines: string[]
  mode: string
  cursor: EditorCursor
  selection: EditorSelection | null
  searchMatch: EditorSearchMatch | null
  hideCursor?: boolean
}

export function EditorArea({
  editorRef,
  lines,
  mode,
  cursor,
  selection,
  searchMatch,
  hideCursor,
}: EditorAreaProps) {
  const cursorShape = cursorShapeForMode(mode)

  return (
    <div ref={editorRef} class="editor-area" aria-label="editor">
      {(lines ?? []).map((line, index) => (
        <div key={index} class="editor-line">
          <EditorLine
            line={line}
            row={index}
            cursor={hideCursor ? { ...cursor, visible: false } : cursor}
            cursorShape={cursorShape}
            selection={selection}
            searchMatch={searchMatch}
          />
        </div>
      ))}
    </div>
  )
}
