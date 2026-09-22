export type PromptKind = "command" | "search-forward" | "search-backward"

export type CommandLineState = {
  active: boolean
  kind: PromptKind
  text: string
  error: string | null
  historyIndex: number | null
}

export function initialCommandLineState(): CommandLineState {
  return {
    active: false,
    kind: "command",
    text: "",
    error: null,
    historyIndex: null,
  }
}

export function isSearchPrompt(kind: PromptKind): boolean {
  return kind === "search-forward" || kind === "search-backward"
}
