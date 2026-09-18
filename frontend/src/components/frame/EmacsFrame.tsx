import type { RefObject } from "preact"
import { CommandLine } from "../chrome/CommandLine"
import { StatusLine } from "../chrome/StatusLine"
import { EditorArea } from "../editor/EditorArea"
import type { CommandLineState } from "../../types/command"
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
  commandLine,
}: EmacsFrameProps) {
  const statusMode = commandLine.active ? "command" : mode

  return (
    <div class="emacs-frame">
      <EditorArea
        editorRef={editorRef}
        lines={lines}
        cursor={cursor}
        hideCursor={commandLine.active}
      />
      <StatusLine
        bufferName={bufferName}
        modified={modified}
        mode={statusMode}
      />
      <CommandLine
        active={commandLine.active}
        text={commandLine.text}
        error={commandLine.error}
      />
    </div>
  )
}
