type StatusLineProps = {
  bufferName: string
  mode: string
}

export function StatusLine({ bufferName, mode }: StatusLineProps) {
  return (
    <div class="status-line" aria-label="status line">
      <span class="status-left">--**- {bufferName}</span>
      <span class="status-right">({mode})</span>
    </div>
  )
}
