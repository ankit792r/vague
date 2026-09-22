type StatusLineProps = {
  bufferName: string
  modified: boolean
  mode: string
  line: number
  column: number
  lineMarks?: string
}

export function StatusLine({
  bufferName,
  modified,
  mode,
  line,
  column,
  lineMarks,
}: StatusLineProps) {
  const displayName = modified ? `${bufferName}*` : bufferName
  const markSuffix = lineMarks ? ` '${lineMarks}'` : ""

  return (
    <div class="status-line" aria-label="status line">
      <span class="status-left">--**- {displayName}</span>
      <span class="status-right">
        Ln {line}, Col {column}
        {markSuffix} ({mode})
      </span>
    </div>
  )
}
