import type { RefObject } from "preact"
import { CommandLine } from "../chrome/CommandLine"
import { StatusLine } from "../chrome/StatusLine"
import { TabLine } from "../chrome/TabLine"
import { SplitEditor } from "../editor/SplitEditor"
import type { CommandLineState, PromptKind } from "../../types/command"
import { isSearchPrompt } from "../../types/command"
import type { EditorViewState } from "../../types/editor"

type EmacsFrameProps = EditorViewState & {
  editorRef: RefObject<HTMLDivElement>
  commandLine: CommandLineState
}

export function EmacsFrame({
  editorRef,
  lines,
  bufferName,
  modified,
  mode,
  cursor,
  selection,
  searchMatch,
  searchHighlights,
  lineNumbers,
  number,
  gutterColumns,
  echo,
  commandLine,
  position,
  lineMarks,
  panes,
  tabs,
  columns,
  rows,
}: EmacsFrameProps) {
  const statusMode = commandLine.active
    ? promptStatusMode(commandLine.kind)
    : mode

  return (
    <div class="emacs-frame">
      <TabLine tabs={tabs} />
      <SplitEditor
        editorRef={editorRef}
        panes={panes.length > 0 ? panes : [{
          windowId: 0,
          x: 0,
          y: 0,
          columns,
          rows,
          active: true,
          lines,
          lineNumbers,
          number,
          gutterColumns,
          bufferName,
          modified,
          mode,
          position,
          cursor,
          selection,
          searchMatch,
          searchHighlights,
          lineMarks,
        }]}
        columns={columns}
        rows={rows}
        hideCursor={commandLine.active}
        mode={mode}
      />
      <StatusLine
        bufferName={bufferName}
        modified={modified}
        mode={statusMode}
        line={position.line}
        column={position.column}
        lineMarks={lineMarks}
      />
      <CommandLine
        active={commandLine.active}
        kind={commandLine.kind}
        text={commandLine.text}
        cursor={commandLine.cursor}
        error={commandLine.error}
        echo={echo}
        completions={commandLine.completions}
        completionIndex={commandLine.completionIndex}
        pendingRegister={commandLine.pendingRegister}
      />
    </div>
  )
}

function promptStatusMode(kind: PromptKind): string {
  if (isSearchPrompt(kind)) {
    return "search"
  }
  return "command"
}
