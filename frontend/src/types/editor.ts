import type { RedrawPayload, StatusEcho } from "../host/protocol"

export type EditorCursor = RedrawPayload["cursor"]
export type EditorSelection = NonNullable<RedrawPayload["selection"]>
export type EditorSearchMatch = NonNullable<RedrawPayload["search_match"]>

export type EditorViewState = {
  lines: string[]
  bufferName: string
  modified: boolean
  mode: string
  position: { line: number; column: number }
  cursor: EditorCursor
  selection: EditorSelection | null
  searchMatch: EditorSearchMatch | null
  echo: StatusEcho | null
}

export function initialEditorViewState(): EditorViewState {
  return {
    lines: [],
    bufferName: "*scratch*",
    modified: false,
    mode: "normal",
    position: { line: 1, column: 1 },
    cursor: { row: 0, column: 0, visible: true },
    selection: null,
    searchMatch: null,
    echo: null,
  }
}

export function editorViewFromRedraw(
  redraw: Partial<RedrawPayload>,
): EditorViewState {
  return {
    bufferName: redraw.buffer?.name ?? "*scratch*",
    modified: redraw.buffer?.modified ?? false,
    lines: Array.isArray(redraw.lines) ? redraw.lines : [],
    mode: redraw.mode ?? "normal",
    position: redraw.position ?? { line: 1, column: 1 },
    cursor: redraw.cursor ?? { row: 0, column: 0, visible: false },
    selection: redraw.selection?.visible ? redraw.selection : null,
    searchMatch: redraw.search_match?.visible ? redraw.search_match : null,
    echo: redraw.echo?.message ? redraw.echo : null,
  }
}
