import { Fragment } from "preact/jsx-runtime"
import type { EditorCursor } from "../../types/editor"

type EditorLineProps = {
  line: string
  row: number
  cursor: EditorCursor
}

export function EditorLine({ line, row, cursor }: EditorLineProps) {
  const display = line === "" ? "\u00a0" : line

  if (!cursor?.visible || cursor.row !== row) {
    return <>{display}</>
  }

  const before = display.slice(0, cursor.column)
  const at = display[cursor.column] ?? "\u00a0"
  const after = display.slice(cursor.column + 1)

  return (
    <Fragment>
      {before}
      <span class="cursor-cell">{at}</span>
      {after}
    </Fragment>
  )
}
