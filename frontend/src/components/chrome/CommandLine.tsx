import type { StatusEcho } from "../../host/protocol"
import {
  type CommandLineState,
  type PromptKind,
} from "../../types/command"

type CommandLineProps = Pick<
  CommandLineState,
  | "active"
  | "kind"
  | "text"
  | "cursor"
  | "error"
  | "completions"
  | "completionIndex"
  | "pendingRegister"
> & {
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

export function CommandLine({
  active,
  kind,
  text,
  cursor,
  error,
  echo,
  completions,
  completionIndex,
  pendingRegister,
}: CommandLineProps) {
  const echoMessage = !active && echo?.message ? echo.message : null
  const echoIsError = echo?.kind === "error"
  const prefix = promptPrefix(kind)
  const before = text.slice(0, cursor)
  const after = text.slice(cursor)
  const cursorCol = prefix.length + cursor

  if (active) {
    return (
      <div class="command-line-wrap">
        <div class="command-line command-line-active" aria-label="command line">
          <span class="command-text">
            <span class="command-input">
              {prefix}
              {before}
              {after}
            </span>
            <span
              class="cursor-bar command-cursor"
              style={{ left: `${cursorCol}ch` }}
              aria-hidden="true"
            />
          </span>
          {pendingRegister ? (
            <span class="command-hint">Insert register</span>
          ) : null}
          {error ? (
            <span class="command-echo command-echo-error command-echo-after-input">
              {error}
            </span>
          ) : null}
        </div>
        {completions.length > 0 ? (
          <div class="command-wildmenu" role="listbox">
            {completions.map((item, i) => (
              <span
                key={`${item}-${i}`}
                class={`command-wildmenu-item${i === completionIndex ? " command-wildmenu-selected" : ""}`}
              >
                {item}
              </span>
            ))}
          </div>
        ) : null}
      </div>
    )
  }

  if (!echoMessage) {
    return <div class="command-line-wrap" />
  }

  return (
    <div class="command-line-wrap">
      <div class="command-line">
        <span
          class={`command-echo${echoIsError ? " command-echo-error" : ""}${echoMessage.includes("\n") ? " command-echo-multiline" : ""}`}
          aria-live="polite"
        >
          {echoMessage}
        </span>
      </div>
    </div>
  )
}
