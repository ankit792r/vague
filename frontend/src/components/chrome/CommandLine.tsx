import type { CommandLineState } from "../../types/command"

type CommandLineProps = Pick<CommandLineState, "active" | "text" | "error">

export function CommandLine({ active, text, error }: CommandLineProps) {
  return (
    <div
      class={`command-line${active ? " command-line-active" : ""}`}
      aria-label="command line"
      aria-hidden={!active}
    >
      {active ? (
        <>
          <span class="command-text">:{text}</span>
          {error ? <span class="command-error">{error}</span> : null}
        </>
      ) : null}
    </div>
  )
}
