import type { EditorSelection } from "../types/editor"

/** Whether this viewport row uses a full-width selection background. */
export function rowSelectionHighlight(
  row: number,
  selection: EditorSelection | null,
): "full" | "partial" | null {
  if (!selection?.visible) {
    return null
  }

  const lo = Math.min(selection.start.row, selection.end.row)
  const hi = Math.max(selection.start.row, selection.end.row)
  if (row < lo || row > hi) {
    return null
  }

  if (selection.linewise) {
    return "full"
  }

  if (lo === hi) {
    return "partial"
  }

  if (row > lo && row < hi) {
    return "full"
  }

  return "partial"
}

export function selectionForLineSpans(
  row: number,
  selection: EditorSelection | null,
): EditorSelection | null {
  if (!selection?.visible) {
    return null
  }
  if (rowSelectionHighlight(row, selection) === "full") {
    return null
  }
  return selection
}
