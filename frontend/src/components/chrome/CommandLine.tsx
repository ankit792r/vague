import { useLayoutEffect, useRef } from "preact/hooks"
import type { StatusEcho } from "../../host/protocol"
import {
  type CommandLineState,
  type PromptKind,
} from "../../types/command"
import { syncWildmenuScroll } from "../../utils/wildmenuScroll"

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
  const wildmenuRef = useRef<HTMLDivElement>(null)
  const wildmenuItemRefs = useRef<(HTMLSpanElement | null)[]>([])

  useLayoutEffect(() => {
    if (!active || completions.length === 0) {
      wildmenuItemRefs.current = []
      if (wildmenuRef.current) {
        wildmenuRef.current.scrollLeft = 0
      }
      return
    }
    const menu = wildmenuRef.current
    if (!menu) {
      return
    }
    syncWildmenuScroll(menu, wildmenuItemRefs.current, completionIndex)
  }, [active, completions, completionIndex])

  useLayoutEffect(() => {
    const menu = wildmenuRef.current
    if (!menu || !active || completions.length === 0) {
      return
    }

    const ro = new ResizeObserver(() => {
      syncWildmenuScroll(menu, wildmenuItemRefs.current, completionIndex)
    })
    ro.observe(menu)
    return () => ro.disconnect()
  }, [active, completions, completionIndex])

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
          <div
            ref={wildmenuRef}
            class="command-wildmenu"
            role="listbox"
            aria-label="command completion"
          >
            {completions.map((item, i) => (
              <span
                ref={(el) => {
                  wildmenuItemRefs.current[i] = el
                }}
                key={`${item}-${i}`}
                class={`command-wildmenu-item${i === completionIndex ? " command-wildmenu-selected" : ""}`}
                role="option"
                aria-selected={i === completionIndex}
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

  const multiline = echoMessage.includes("\n")

  if (multiline) {
    return (
      <div class="command-line-wrap">
        <div
          class={`command-minibuffer command-minibuffer-expanded${echoIsError ? " command-minibuffer-error" : ""}`}
          aria-live="polite"
        >
          <pre class="command-minibuffer-body">{echoMessage}</pre>
        </div>
      </div>
    )
  }

  return (
    <div class="command-line-wrap">
      <div class="command-line">
        <span
          class={`command-echo${echoIsError ? " command-echo-error" : ""}`}
          aria-live="polite"
        >
          {echoMessage}
        </span>
      </div>
    </div>
  )
}
