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
  relativeNumber: boolean
  gutterColumns: number
  signColumn: boolean
  cursorLineRow: number
  cursorColumn: number
  cursorColumnOn: boolean
  colorColumns: number[]
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
  relativeNumber,
  gutterColumns,
  signColumn,
  cursorLineRow,
  cursorColumn,
  cursorColumnOn,
  colorColumns,
  mode,
  cursor,
  selection,
  searchMatch,
  searchHighlights,
  hideCursor,
}: EditorAreaProps) {
  const cursorShape = cursorShapeForMode(mode)
  const showGutter = number || relativeNumber
  const gutterStyle =
    showGutter && gutterColumns > 0
      ? { width: `${gutterColumns}ch` }
      : undefined

  return (
    <div ref={editorRef} class="editor-area" aria-label="editor">
      {(lines ?? []).map((line, index) => {
        const lineClass =
          cursorLineRow === index ? "editor-line cursor-line" : "editor-line"
        return (
          <div key={index} class={lineClass}>
            {signColumn ? (
              <span class="sign-column" aria-hidden="true">
                {" "}
              </span>
            ) : null}
            {showGutter ? (
              <span class="line-number" style={gutterStyle} aria-hidden="true">
                {lineNumbers[index] ?? ""}
              </span>
            ) : null}
            <div class="editor-line-content">
              {colorColumns.map((col) => (
                <span
                  key={`cc-${col}`}
                  class="color-column"
                  style={{ left: `${col}ch` }}
                  aria-hidden="true"
                />
              ))}
              {cursorColumnOn ? (
                <span
                  class="cursor-column"
                  style={{ left: `${cursorColumn}ch` }}
                  aria-hidden="true"
                />
              ) : null}
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
        )
      })}
    </div>
  )
}
