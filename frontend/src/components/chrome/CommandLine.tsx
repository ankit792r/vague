const DUMMY_COMMAND = "open-file"

export function CommandLine() {
  return (
    <div class="command-line" aria-label="command line">
      <span class="command-prompt">:</span>
      <span class="command-text">{DUMMY_COMMAND}</span>
    </div>
  )
}
