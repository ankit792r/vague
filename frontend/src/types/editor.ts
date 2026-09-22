import type { RedrawPayload, StatusEcho } from "../host/protocol"

export type EditorCursor = RedrawPayload["cursor"]
export type EditorSelection = NonNullable<RedrawPayload["selection"]>
export type EditorSearchMatch = NonNullable<RedrawPayload["search_match"]>
export type EditorSearchHighlight = NonNullable<
  RedrawPayload["search_highlights"]
>[number]

export type EditorPaneState = {
  windowId: number
  x: number
  y: number
  columns: number
  rows: number
  active: boolean
  lines: string[]
  lineNumbers: number[]
  number: boolean
  relativeNumber: boolean
  gutterColumns: number
  signColumn: boolean
  list: boolean
  cursorLineRow: number
  cursorColumn: number
  cursorColumnOn: boolean
  colorColumns: number[]
  bufferName: string
  modified: boolean
  mode: string
  position: { line: number; column: number }
  cursor: EditorCursor
  selection: EditorSelection | null
  searchMatch: EditorSearchMatch | null
  searchHighlights: EditorSearchHighlight[]
  lineMarks: string
  statusLine: string
  showMode: boolean
  ruler: boolean
}

export type EditorTabState = {
  label: string
  active: boolean
}

export type EditorViewState = {
  columns: number
  rows: number
  panes: EditorPaneState[]
  tabs: EditorTabState[]
  activeTab: number
  lines: string[]
  lineNumbers: number[]
  number: boolean
  relativeNumber: boolean
  gutterColumns: number
  signColumn: boolean
  bufferName: string
  modified: boolean
  mode: string
  position: { line: number; column: number }
  cursor: EditorCursor
  selection: EditorSelection | null
  searchMatch: EditorSearchMatch | null
  searchHighlights: EditorSearchHighlight[]
  echo: StatusEcho | null
  lineMarks: string
  theme: string
  showCmd: string
  statusLine: string
  showMode: boolean
  ruler: boolean
}

export function initialEditorViewState(): EditorViewState {
  return {
    columns: 80,
    rows: 24,
    panes: [],
    tabs: [],
    activeTab: 0,
    lines: [],
    lineNumbers: [],
    number: false,
    relativeNumber: false,
    gutterColumns: 0,
    signColumn: false,
    bufferName: "*scratch*",
    modified: false,
    mode: "normal",
    position: { line: 1, column: 1 },
    cursor: { row: 0, column: 0, visible: true },
    selection: null,
    searchMatch: null,
    searchHighlights: [],
    echo: null,
    lineMarks: "",
    theme: "vague",
    showCmd: "",
    statusLine: "",
    showMode: true,
    ruler: true,
  }
}

export function editorViewFromRedraw(
  redraw: Partial<RedrawPayload>,
): EditorViewState {
  const columns = redraw.columns ?? 80
  const rows = redraw.rows ?? 24
  const panes =
    Array.isArray(redraw.panes) && redraw.panes.length > 0
      ? redraw.panes.map(paneFromPayload)
      : [singlePaneFromRedraw(redraw)]

  const tabs = Array.isArray(redraw.tabs)
    ? redraw.tabs.map((t) => ({
        label: t.label,
        active: t.active ?? false,
      }))
    : []

  const active =
    panes.find((p) => p.active) ?? panes[0] ?? singlePaneFromRedraw(redraw)

  return {
    columns,
    rows,
    panes,
    tabs,
    activeTab: redraw.active_tab ?? 0,
    bufferName: active.bufferName,
    modified: active.modified,
    lines: active.lines,
    lineNumbers: active.lineNumbers,
    number: active.number,
    relativeNumber: active.relativeNumber,
    gutterColumns: active.gutterColumns,
    signColumn: active.signColumn,
    mode: active.mode,
    position: active.position,
    cursor: active.cursor,
    selection: active.selection,
    searchMatch: active.searchMatch,
    searchHighlights: active.searchHighlights,
    echo: redraw.echo?.message ? redraw.echo : null,
    lineMarks: active.lineMarks,
    theme: redraw.theme ?? "vague",
    showCmd: redraw.showcmd ?? "",
    statusLine: active.statusLine,
    showMode: active.showMode,
    ruler: active.ruler,
  }
}

function paneFromPayload(
  p: NonNullable<RedrawPayload["panes"]>[number],
): EditorPaneState {
  return {
    windowId: p.window_id,
    x: p.x,
    y: p.y,
    columns: p.columns,
    rows: p.rows,
    active: p.active ?? false,
    lines: Array.isArray(p.lines) ? p.lines : [],
    lineNumbers: Array.isArray(p.line_numbers) ? p.line_numbers : [],
    number: p.number ?? false,
    relativeNumber: p.relative_number ?? false,
    gutterColumns: p.gutter_columns ?? 0,
    signColumn: p.sign_column ?? false,
    list: p.list ?? false,
    cursorLineRow: p.cursor_line_row ?? -1,
    cursorColumn: p.cursor_column ?? 0,
    cursorColumnOn: p.cursor_column_on ?? false,
    colorColumns: Array.isArray(p.color_columns) ? p.color_columns : [],
    bufferName: p.buffer?.name ?? "*scratch*",
    modified: p.buffer?.modified ?? false,
    mode: p.mode ?? "normal",
    position: p.position ?? { line: 1, column: 1 },
    cursor: p.cursor ?? { row: 0, column: 0, visible: false },
    selection: p.selection?.visible ? p.selection : null,
    searchMatch: p.search_match?.visible ? p.search_match : null,
    searchHighlights: Array.isArray(p.search_highlights)
      ? p.search_highlights.filter((h) => h.visible)
      : [],
    lineMarks: p.line_marks ?? "",
    statusLine: p.statusline ?? "",
    showMode: p.showmode ?? true,
    ruler: p.ruler ?? true,
  }
}

function singlePaneFromRedraw(redraw: Partial<RedrawPayload>): EditorPaneState {
  return {
    windowId: 0,
    x: 0,
    y: 0,
    columns: redraw.columns ?? 80,
    rows: redraw.rows ?? 24,
    active: true,
    lines: Array.isArray(redraw.lines) ? redraw.lines : [],
    lineNumbers: Array.isArray(redraw.line_numbers) ? redraw.line_numbers : [],
    number: redraw.number ?? false,
    relativeNumber: false,
    gutterColumns: redraw.gutter_columns ?? 0,
    signColumn: false,
    list: false,
    cursorLineRow: -1,
    cursorColumn: redraw.cursor?.column ?? 0,
    cursorColumnOn: false,
    colorColumns: [],
    bufferName: redraw.buffer?.name ?? "*scratch*",
    modified: redraw.buffer?.modified ?? false,
    mode: redraw.mode ?? "normal",
    position: redraw.position ?? { line: 1, column: 1 },
    cursor: redraw.cursor ?? { row: 0, column: 0, visible: false },
    selection: redraw.selection?.visible ? redraw.selection : null,
    searchMatch: redraw.search_match?.visible ? redraw.search_match : null,
    searchHighlights: Array.isArray(redraw.search_highlights)
      ? redraw.search_highlights.filter((h) => h.visible)
      : [],
    lineMarks: redraw.line_marks ?? "",
    statusLine: "",
    showMode: true,
    ruler: true,
  }
}
