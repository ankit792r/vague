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
  echo,
  commandLine,
}: EmacsFrameProps) {
  const statusMode = commandLine.active
    ? promptStatusMode(commandLine.kind)
    : mode

  return (
    <div class="emacs-frame">
      <EditorArea
        editorRef={editorRef}
        lines={lines}
        mode={mode}
        cursor={cursor}
        selection={selection}
        hideCursor={commandLine.active}
      />
      <StatusLine
        bufferName={bufferName}
        modified={modified}
        mode={statusMode}
      />
      <CommandLine
        active={commandLine.active}
        kind={commandLine.kind}
        text={commandLine.text}
        error={commandLine.error}
        echo={echo}
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
