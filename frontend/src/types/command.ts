export type PromptKind = "command" | "search-forward" | "search-backward"

export type CommandLineState = {
  active: boolean
  kind: PromptKind
  text: string
  cursor: number
  error: string | null
  historyIndex: number | null
  completions: string[]
  completionIndex: number
  pendingRegister: boolean
}

export function initialCommandLineState(): CommandLineState {
  return {
    active: false,
    kind: "command",
    text: "",
    cursor: 0,
    error: null,
    historyIndex: null,
    completions: [],
    completionIndex: 0,
    pendingRegister: false,
  }
}

export function isSearchPrompt(kind: PromptKind): boolean {
  return kind === "search-forward" || kind === "search-backward"
}

export function withText(state: CommandLineState, text: string, cursor: number): CommandLineState {
  return {
    ...state,
    text,
    cursor: Math.min(Math.max(0, cursor), text.length),
    error: null,
    historyIndex: null,
    completions: [],
    completionIndex: 0,
  }
}
