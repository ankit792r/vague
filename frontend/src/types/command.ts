export type CommandLineState = {
  active: boolean
  text: string
  error: string | null
}

export function initialCommandLineState(): CommandLineState {
  return {
    active: false,
    text: "",
    error: null,
  }
}
