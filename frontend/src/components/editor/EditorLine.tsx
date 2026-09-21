import { Fragment, type JSX } from "preact/jsx-runtime"
import type { EditorCursor, EditorSelection } from "../../types/editor"

type EditorLineProps = {
  line: string
  row: number
  cursor: EditorCursor
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

export function EditorLine({ line, row, cursor, selection }: EditorLineProps) {
  const display = line === "" ? "\u00a0" : line
  const span = selectionSpanForRow(row, display.length, selection)
  const showCursor = cursor?.visible && cursor.row === row

  if (!span && !showCursor) {
    return <>{display}</>
  }

  const nodes: JSX.Element[] = []
  let run = ""
  let runClass: string | undefined

  const flush = () => {
    if (run === "") {
      return
    }
    if (runClass) {
      nodes.push(
        <span key={nodes.length} class={runClass}>
          {run}
        </span>,
      )
    } else {
      nodes.push(<Fragment key={nodes.length}>{run}</Fragment>)
    }
    run = ""
    runClass = undefined
  }

  for (let col = 0; col < display.length; col++) {
    const ch = display[col]
    let className: string | undefined
    const inSelection = span && col >= span.start && col < span.end
    const atCursor = showCursor && cursor.column === col

    if (inSelection && atCursor) {
      className = "visual-selection cursor-cell"
    } else if (inSelection) {
      className = "visual-selection"
    } else if (atCursor) {
      className = "cursor-cell"
    }

    if (className !== runClass) {
      flush()
      runClass = className
    }
    run += ch
  }
  flush()

  return <>{nodes}</>
}
