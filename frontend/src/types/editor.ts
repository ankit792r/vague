import type { RedrawPayload, StatusEcho } from "../host/protocol"

export type EditorCursor = RedrawPayload["cursor"]

export type EditorViewState = {
  lines: string[]
  bufferName: string
  modified: boolean
  mode: string
  cursor: EditorCursor
  echo: StatusEcho | null
}

export function initialEditorViewState(): EditorViewState {
  return {
    lines: [],
    bufferName: "*scratch*",
    modified: false,
    mode: "normal",
    cursor: { row: 0, column: 0, visible: true },
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
    cursor: redraw.cursor ?? { row: 0, column: 0, visible: false },
    echo: redraw.echo?.message ? redraw.echo : null,
  }
}
