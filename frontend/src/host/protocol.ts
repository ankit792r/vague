/** UI → host: method name → params */
export type HostRequestMap = {
  input: { keys: string }
  execute: {
    name: string
    args?: string[]
    bang?: boolean
    count?: number
  }
  ui_ready: {
    height: number
    width: number
  }
  ui_attach: Record<string, never>
}

/** UI → host: method name → success result */
export type HostResultMap = {
  input: null
  execute: unknown
  ui_attach: {
    session_id: number
    frame_id: number
  }
  ui_ready: {
    session_id: number
    frame_id: number
  }
}

/** Host → UI: event name → payload */
export type StatusEcho = {
  message: string
  kind?: "info" | "error"
}

export type RedrawPayload = {
  frame_id: number
  columns: number
  rows: number
  wrap: boolean
  number?: boolean
  gutter_columns?: number
  full?: boolean
  buffer: {
    id: number
    name: string
    modified?: boolean
  }
  lines: string[]
  line_numbers?: number[]
  cursor: {
    row: number
    column: number
    visible: boolean
  }
  selection?: {
    visible: boolean
    linewise?: boolean
    start: { row: number; column: number }
    end: { row: number; column: number }
  }
  search_match?: {
    visible: boolean
    linewise?: boolean
    start: { row: number; column: number }
    end: { row: number; column: number }
  }
  mode: string
  position: {
    line: number
    column: number
  }
  echo?: StatusEcho
}

export type HostEventMap = {
  attached: HostResultMap["ui_attach"]
  redraw: RedrawPayload
  quit: { frame_id: number }
}
export type HostRequestMethod = keyof HostRequestMap
export type HostEventMethod = keyof HostEventMap
export type HostReply<M extends HostRequestMethod = HostRequestMethod> = {
  id: number
  result?: HostResultMap[M]
  error?: string
}
export type HostEvent<E extends HostEventMethod = HostEventMethod> = {
  event: E
  payload: HostEventMap[E]
}
