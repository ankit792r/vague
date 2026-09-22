import type { EditorPaneState } from "../../types/editor"
import { EditorArea } from "../editor/EditorArea"

type SplitEditorProps = {
  panes: EditorPaneState[]
  columns: number
  rows: number
  editorRef: import("preact").RefObject<HTMLDivElement>
  hideCursor: boolean
  mode: string
}

export function SplitEditor({
  panes,
  columns,
  rows,
  editorRef,
  hideCursor,
  mode,
}: SplitEditorProps) {
  const unitW = columns > 0 ? 100 / columns : 100
  const unitH = rows > 0 ? 100 / rows : 100

  return (
    <div class="split-editor">
      {panes.map((pane) => (
        <div
          key={pane.windowId}
          class={`editor-pane${pane.active ? " editor-pane-active" : ""}`}
          style={{
            left: `${pane.x * unitW}%`,
            top: `${pane.y * unitH}%`,
            width: `${pane.columns * unitW}%`,
            height: `${pane.rows * unitH}%`,
          }}
        >
          <EditorArea
            editorRef={pane.active ? editorRef : { current: null }}
            lines={pane.lines}
            lineNumbers={pane.lineNumbers}
            number={pane.number}
            gutterColumns={pane.gutterColumns}
            mode={mode}
            cursor={pane.cursor}
            selection={pane.selection}
            searchMatch={pane.searchMatch}
            searchHighlights={pane.searchHighlights}
            hideCursor={hideCursor || !pane.active}
          />
        </div>
      ))}
    </div>
  )
}
