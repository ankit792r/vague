import type { StatusEcho } from "../../host/protocol"
import type { CommandLineState } from "../../types/command"

type CommandLineProps = Pick<CommandLineState, "active" | "text" | "error"> & {
  echo: StatusEcho | null
}

export function CommandLine({ active, text, error, echo }: CommandLineProps) {
  const echoMessage = !active && echo?.message ? echo.message : null
  const echoIsError = echo?.kind === "error"

  return (
    <div
      class={`command-line${active ? " command-line-active" : ""}`}
      aria-label="command line"
    >
      {active ? (
        <>
          <span class="command-text">:{text}</span>
          {error ? (
            <span class="command-echo command-echo-error command-echo-after-input">
              {error}
            </span>
          ) : null}
        </>
      ) : echoMessage ? (
        <span
          class={`command-echo${echoIsError ? " command-echo-error" : ""}`}
          aria-live="polite"
        >
          {echoMessage}
        </span>
      ) : null}
    </div>
  )
}
