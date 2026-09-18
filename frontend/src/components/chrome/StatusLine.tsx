type StatusLineProps = {
  bufferName: string
  modified: boolean
  mode: string
}

export function StatusLine({ bufferName, modified, mode }: StatusLineProps) {
  const displayName = modified ? `${bufferName}*` : bufferName

  return (
    <div class="status-line" aria-label="status line">
      <span class="status-left">--**- {displayName}</span>
      <span class="status-right">({mode})</span>
    </div>
  )
}
