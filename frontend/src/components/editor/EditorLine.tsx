import { Fragment, type JSX } from "preact/jsx-runtime"
import type { EditorCursor, EditorSelection } from "../../types/editor"

export type CursorShape = "block" | "bar"

type EditorLineProps = {
  line: string
  row: number
  cursor: EditorCursor
  cursorShape: CursorShape
  selection: EditorSelection | null
}

function selectionSpanForRow(
  row: number,
  lineLen: number,
  selection: EditorSelection | null,
): { start: number; end: number } | null {
  if (!selection?.visible) {
    return null
  }

  if (row < selection.start.row || row > selection.end.row) {
    return null
  }

  if (selection.linewise) {
    return { start: 0, end: lineLen }
  }

  let start = 0
  let end = lineLen
  if (row === selection.start.row) {
    start = selection.start.column
  }
  if (row === selection.end.row) {
    end = Math.min(lineLen, selection.end.column + 1)
  }
  if (end <= start) {
    return null
  }
  return { start, end }
}

function classForSpan(
  colStart: number,
  colEnd: number,
  span: { start: number; end: number } | null,
): string | undefined {
  if (!span) {
    return undefined
  }
  if (colEnd <= span.start || colStart >= span.end) {
    return undefined
  }
  return "visual-selection"
}

function renderSegment(
  text: string,
  className?: string,
  key?: number,
): JSX.Element {
  if (!className) {
    return <Fragment key={key}>{text}</Fragment>
  }
  return (
    <span key={key} class={className}>
      {text}
    </span>
  )
}

export function EditorLine({
  line,
  row,
  cursor,
  cursorShape,
  selection,
}: EditorLineProps) {
  const display = line === "" ? "\u00a0" : line
  const span = selectionSpanForRow(row, display.length, selection)
  const showCursor = cursor?.visible && cursor.row === row
  const cursorCol = showCursor
    ? Math.min(Math.max(0, cursor.column), display.length)
    : -1

  if (!span && !showCursor) {
    return <>{display}</>
  }

  if (cursorShape === "bar" && showCursor) {
    const before = display.slice(0, cursorCol)
    const after = display.slice(cursorCol)
    const beforeClass = classForSpan(0, before.length, span)
    const afterClass = classForSpan(cursorCol, display.length, span)

    return (
      <>
        {renderSegment(before, beforeClass, 0)}
        <span class="cursor-bar" aria-hidden="true" />
        {renderSegment(after, afterClass, 1)}
      </>
    )
  }

  const nodes: JSX.Element[] = []
  let run = ""
  let runClass: string | undefined

  const flush = () => {
    if (run === "") {
      return
    }
    nodes.push(renderSegment(run, runClass, nodes.length))
    run = ""
    runClass = undefined
  }

  for (let col = 0; col < display.length; col++) {
    const ch = display[col]
    let className = classForSpan(col, col + 1, span)
    const atCursor = showCursor && cursor.column === col

    if (atCursor) {
      className = className ? "visual-selection cursor-cell" : "cursor-cell"
    }

    if (className !== runClass) {
      flush()
      runClass = className
    }
    run += ch
  }
  flush()

  if (showCursor && cursor.column >= display.length) {
    const trailingClass = classForSpan(
      display.length,
      display.length,
      span,
    )
    nodes.push(
      <span
        key="cursor-eol"
        class={trailingClass ? "visual-selection cursor-cell" : "cursor-cell"}
      >
        {"\u00a0"}
      </span>,
    )
  }

  return <>{nodes}</>
}
