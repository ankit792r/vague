import type { RefObject } from "preact"
import type { EditorCursor, EditorSearchHighlight, EditorSearchMatch, EditorSelection } from "../../types/editor"
import { EditorLine } from "./EditorLine"

function cursorShapeForMode(mode: string): "block" | "bar" {
  return mode === "insert" ? "bar" : "block"
}

type EditorAreaProps = {
  editorRef: RefObject<HTMLDivElement>
  lines: string[]
  lineNumbers: number[]
  number: boolean
  gutterColumns: number
  mode: string
  cursor: EditorCursor
  selection: EditorSelection | null
  searchMatch: EditorSearchMatch | null
  searchHighlights: EditorSearchHighlight[]
  hideCursor?: boolean
}

export function EditorArea({
  editorRef,
  lines,
  lineNumbers,
  number,
  gutterColumns,
  mode,
  cursor,
  selection,
  searchMatch,
  searchHighlights,
  hideCursor,
}: EditorAreaProps) {
  const cursorShape = cursorShapeForMode(mode)
  const gutterStyle =
    number && gutterColumns > 0
      ? { width: `${gutterColumns}ch` }
      : undefined

  return (
    <div ref={editorRef} class="editor-area" aria-label="editor">
      {(lines ?? []).map((line, index) => (
        <div key={index} class="editor-line">
          {number ? (
            <span class="line-number" style={gutterStyle} aria-hidden="true">
              {lineNumbers[index] ?? ""}
            </span>
          ) : null}
          <div class="editor-line-content">
            <EditorLine
            line={line}
            row={index}
            cursor={hideCursor ? { ...cursor, visible: false } : cursor}
            cursorShape={cursorShape}
            selection={selection}
            searchMatch={searchMatch}
            searchHighlights={searchHighlights}
          />
          </div>
        </div>
      ))}
    </div>
  )
}
