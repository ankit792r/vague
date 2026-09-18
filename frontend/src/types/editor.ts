import type { RedrawPayload } from "../host/protocol"

export type EditorCursor = RedrawPayload["cursor"]

export type EditorViewState = {
  lines: string[]
  bufferName: string
  modified: boolean
  mode: string
  cursor: EditorCursor
}

export function initialEditorViewState(): EditorViewState {
  return {
    lines: [],
    bufferName: "*scratch*",
    modified: false,
    mode: "normal",
    cursor: { row: 0, column: 0, visible: true },
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
  }
}
