type StatusLineProps = {
  bufferName: string
  modified: boolean
  mode: string
  line: number
  column: number
}

export function StatusLine({
  bufferName,
  modified,
  mode,
  line,
  column,
}: StatusLineProps) {
  const displayName = modified ? `${bufferName}*` : bufferName

  return (
    <div class="status-line" aria-label="status line">
      <span class="status-left">--**- {displayName}</span>
      <span class="status-right">
        Ln {line}, Col {column} ({mode})
      </span>
    </div>
  )
}
