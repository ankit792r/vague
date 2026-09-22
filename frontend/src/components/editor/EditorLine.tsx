import { Fragment, type JSX } from "preact/jsx-runtime"
import type { EditorCursor, EditorSearchHighlight, EditorSelection } from "../../types/editor"

export type CursorShape = "block" | "bar"
export type CursorFill = "solid" | "hollow"

type EditorLineProps = {
  line: string
  row: number
  cursor: EditorCursor
  cursorShape: CursorShape
  cursorFill: CursorFill
  selection: EditorSelection | null
  searchMatch: EditorSelection | null
  searchHighlights: EditorSearchHighlight[]
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

function highlightSpansForRow(
  row: number,
  lineLen: number,
  highlights: EditorSearchHighlight[],
): { start: number; end: number }[] {
  const spans: { start: number; end: number }[] = []
  for (const highlight of highlights) {
    const span = highlightSpanForRow(row, lineLen, highlight)
    if (span) {
      spans.push(span)
    }
  }
  return spans
}

function overlapsSpan(
  colStart: number,
  colEnd: number,
  span: { start: number; end: number },
): boolean {
  return colEnd > span.start && colStart < span.end
}

function cellClassName(
  colStart: number,
  colEnd: number,
  selection: { start: number; end: number } | null,
  searchMatch: { start: number; end: number } | null,
  searchHighlights: { start: number; end: number }[],
  atCursor: boolean,
  cursorFill: CursorFill,
): string | undefined {
  const parts: string[] = []

  for (const span of searchHighlights) {
    if (overlapsSpan(colStart, colEnd, span)) {
      parts.push("search-hl")
      break
    }
  }
  if (searchMatch && overlapsSpan(colStart, colEnd, searchMatch)) {
    parts.push("search-match")
  }
  if (selection && colEnd > selection.start && colStart < selection.end) {
    parts.push("visual-selection")
  }
  if (atCursor) {
    parts.push(cursorFill === "hollow" ? "cursor-cell-hollow" : "cursor-cell")
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
  searchHighlights: { start: number; end: number }[],
  showBlockCursor: boolean,
  cursorColumn: number,
  cursorFill: CursorFill,
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
      searchHighlights,
      showBlockCursor && cursorColumn === col,
      cursorFill,
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
      searchHighlights,
      true,
      cursorFill,
    )
    nodes.push(
      <span
        key="cursor-eol"
        class={
          trailingClass ??
          (cursorFill === "hollow" ? "cursor-cell-hollow" : "cursor-cell")
        }
      >
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
  cursorFill,
  selection,
  searchMatch,
  searchHighlights,
}: EditorLineProps) {
  const display = line === "" ? "\u00a0" : line
  const selectionSpan = highlightSpanForRow(row, display.length, selection)
  const searchSpan = highlightSpanForRow(row, display.length, searchMatch)
  const hlSpans = highlightSpansForRow(row, display.length, searchHighlights)
  const showCursor = cursor?.visible && cursor.row === row
  const cursorCol = showCursor
    ? Math.min(Math.max(0, cursor.column), display.length)
    : -1

  const useBarOverlay = cursorShape === "bar" && showCursor
  const useBlockCursor = showCursor && !useBarOverlay

  if (!selectionSpan && !searchSpan && hlSpans.length === 0 && !showCursor) {
    return <>{display}</>
  }

  const text = renderLineText(
    display,
    selectionSpan,
    searchSpan,
    hlSpans,
    useBlockCursor,
    cursorCol,
    cursorFill,
  )

  if (!useBarOverlay) {
    return text
  }

  const barClass =
    cursorFill === "hollow" ? "cursor-bar-hollow" : "cursor-bar"

  return (
    <>
      {text}
      <span
        class={barClass}
        style={{ left: `${cursorCol}ch` }}
        aria-hidden="true"
      />
    </>
  )
}
