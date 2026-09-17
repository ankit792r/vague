import type { RedrawPayload } from "../host/protocol"

export type EditorCursor = RedrawPayload["cursor"]

export type EditorViewState = {
  lines: string[]
  bufferName: string
  mode: string
  cursor: EditorCursor
}

export function initialEditorViewState(): EditorViewState {
  return {
    lines: [],
    bufferName: "*scratch*",
    mode: "normal",
    cursor: { row: 0, column: 0, visible: true },
  }
}

export function editorViewFromRedraw(
  redraw: Partial<RedrawPayload>,
): EditorViewState {
  return {
    bufferName: redraw.buffer?.name ?? "*scratch*",
    lines: Array.isArray(redraw.lines) ? redraw.lines : [],
    mode: redraw.mode ?? "normal",
    cursor: redraw.cursor ?? { row: 0, column: 0, visible: false },
  }
}
