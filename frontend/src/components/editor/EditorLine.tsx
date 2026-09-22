import { Fragment, type JSX } from "preact/jsx-runtime"
import type { EditorCursor, EditorSelection } from "../../types/editor"

export type CursorShape = "block" | "bar"

type EditorLineProps = {
  line: string
  row: number
  cursor: EditorCursor
  cursorShape: CursorShape
  selection: EditorSelection | null
  searchMatch: EditorSelection | null
}

function highlightSpanForRow(
  row: number,
  lineLen: number,
  highlight: EditorSelection | null,
): { start: number; end: number } | null {
  if (!highlight?.visible) {
    return null
  }

  if (row < highlight.start.row || row > highlight.end.row) {
    return null
  }

  if (highlight.linewise) {
    return { start: 0, end: lineLen }
  }

  let start = 0
  let end = lineLen
  if (row === highlight.start.row) {
    start = highlight.start.column
  }
  if (row === highlight.end.row) {
    end = Math.min(lineLen, highlight.end.column + 1)
  }
  if (end <= start) {
    return null
  }
  return { start, end }
}

function cellClassName(
  colStart: number,
  colEnd: number,
  selection: { start: number; end: number } | null,
  searchMatch: { start: number; end: number } | null,
  atCursor: boolean,
): string | undefined {
  const parts: string[] = []

  if (searchMatch && colEnd > searchMatch.start && colStart < searchMatch.end) {
    parts.push("search-match")
  }
  if (selection && colEnd > selection.start && colStart < selection.end) {
    parts.push("visual-selection")
  }
  if (atCursor) {
    parts.push("cursor-cell")
  }

  if (parts.length === 0) {
    return undefined
  }

  return parts.join(" ")
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

function renderLineText(
  display: string,
  selection: { start: number; end: number } | null,
  searchMatch: { start: number; end: number } | null,
  showBlockCursor: boolean,
  cursorColumn: number,
): JSX.Element {
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
    const className = cellClassName(
      col,
      col + 1,
      selection,
      searchMatch,
      showBlockCursor && cursorColumn === col,
    )

    if (className !== runClass) {
      flush()
      runClass = className
    }
    run += ch
  }
  flush()

  if (showBlockCursor && cursorColumn >= display.length) {
    const trailingClass = cellClassName(
      display.length,
      display.length,
      selection,
      searchMatch,
      true,
    )
    nodes.push(
      <span key="cursor-eol" class={trailingClass ?? "cursor-cell"}>
        {"\u00a0"}
      </span>,
    )
  }

  return <>{nodes}</>
}

export function EditorLine({
  line,
  row,
  cursor,
  cursorShape,
  selection,
  searchMatch,
}: EditorLineProps) {
  const display = line === "" ? "\u00a0" : line
  const selectionSpan = highlightSpanForRow(row, display.length, selection)
  const searchSpan = highlightSpanForRow(row, display.length, searchMatch)
  const showCursor = cursor?.visible && cursor.row === row
  const cursorCol = showCursor
    ? Math.min(Math.max(0, cursor.column), display.length)
    : -1

  const useBarOverlay = cursorShape === "bar" && showCursor
  const useBlockCursor = showCursor && !useBarOverlay

  if (!selectionSpan && !searchSpan && !showCursor) {
    return <>{display}</>
  }

  const text = renderLineText(
    display,
    selectionSpan,
    searchSpan,
    useBlockCursor,
    cursorCol,
  )

  if (!useBarOverlay) {
    return text
  }

  return (
    <>
      {text}
      <span
        class="cursor-bar"
        style={{ left: `${cursorCol}ch` }}
        aria-hidden="true"
      />
    </>
  )
}
