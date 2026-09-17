import type { CommandLineState } from "../../types/command"

type CommandLineProps = Pick<CommandLineState, "active" | "text" | "error">

export function CommandLine({ active, text, error }: CommandLineProps) {
  return (
    <div
      class={`command-line${active ? " command-line-active" : ""}`}
      aria-label="command line"
    >
      <span class="command-prompt">:</span>
      <span class="command-text">{active ? text : ""}</span>
      {error ? <span class="command-error">{error}</span> : null}
    </div>
  )
}
