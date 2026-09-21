import type { StatusEcho } from "../../host/protocol"
import {
  type CommandLineState,
  type PromptKind,
} from "../../types/command"

type CommandLineProps = Pick<CommandLineState, "active" | "kind" | "text" | "error"> & {
  echo: StatusEcho | null
}

function promptPrefix(kind: PromptKind): string {
  switch (kind) {
    case "search-forward":
      return "/"
    case "search-backward":
      return "?"
    default:
      return ":"
  }
}

export function CommandLine({ active, kind, text, error, echo }: CommandLineProps) {
  const echoMessage = !active && echo?.message ? echo.message : null
  const echoIsError = echo?.kind === "error"

  return (
    <div
      class={`command-line${active ? " command-line-active" : ""}`}
      aria-label="command line"
    >
      {active ? (
        <>
          <span class="command-text">
            {promptPrefix(kind)}
            {text}
          </span>
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
