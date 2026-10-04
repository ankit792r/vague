import { useState } from "preact/hooks"

export function CommandLine() {
  const [error, setError] = useState("")
  const [active, setActive] = useState(true)
  const [prefix, setPrefix] = useState(":")

  if (active) {
    return (
      <div class="command-line command-line-active" aria-label="command line">
        <span class="command-text">
          <span class="command-input">
            {prefix}
          </span>
        </span>
        {error ? (
          <span class="command-echo command-echo-error command-echo-after-input">
            {error}
          </span>
        ) : null}
      </div>
    )
  }
}
