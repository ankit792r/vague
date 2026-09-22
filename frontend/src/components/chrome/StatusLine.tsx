type StatusLineProps = {
  bufferName: string
  modified: boolean
  mode: string
  line: number
  column: number
  lineMarks?: string
  statusLine?: string
  showMode?: boolean
  ruler?: boolean
}

export function StatusLine({
  bufferName,
  modified,
  mode,
  line,
  column,
  lineMarks,
  statusLine,
  showMode = true,
  ruler = true,
}: StatusLineProps) {
  const markSuffix = lineMarks ? ` '${lineMarks}'` : ""

  if (statusLine) {
    return (
      <div class="status-line" aria-label="status line">
        <span class="status-left">{statusLine}</span>
        {showMode ? (
          <span class="status-right">({mode})</span>
        ) : null}
      </div>
    )
  }

  const displayName = modified ? `${bufferName}*` : bufferName
  const rulerText = ruler ? `Ln ${line}, Col ${column}${markSuffix}` : markSuffix

  return (
    <div class="status-line" aria-label="status line">
      <span class="status-left">--**- {displayName}</span>
      <span class="status-right">
        {ruler ? rulerText : null}
        {showMode ? ` (${mode})` : null}
      </span>
    </div>
  )
}
