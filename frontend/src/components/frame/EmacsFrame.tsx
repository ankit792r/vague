import type { RefObject } from "preact"
import { CommandLine } from "../chrome/CommandLine"
import { StatusLine } from "../chrome/StatusLine"
import { EditorArea } from "../editor/EditorArea"
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
  lineNumbers,
  number,
  gutterColumns,
  echo,
  commandLine,
  position,
}: EmacsFrameProps) {
  const statusMode = commandLine.active
    ? promptStatusMode(commandLine.kind)
    : mode

  return (
    <div class="emacs-frame">
      <EditorArea
        editorRef={editorRef}
        lines={lines}
        lineNumbers={lineNumbers}
        number={number}
        gutterColumns={gutterColumns}
        mode={mode}
        cursor={cursor}
        selection={selection}
        searchMatch={searchMatch}
        hideCursor={commandLine.active}
      />
      <StatusLine
        bufferName={bufferName}
        modified={modified}
        mode={statusMode}
        line={position.line}
        column={position.column}
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
