/** UI → host: method name → params */
export type HostRequestMap = {
  input: { keys: string }
  execute: {
    name: string
    args?: string[]
    bang?: boolean
    count?: number
  }
  ready: {
    height: number,
    width: number
  }
  attach: Record<string, never>
}


/** UI → host: method name → success result */
export type HostResultMap = {
  input: null
  execute: unknown
  attach: {
    session_id: number
    frame_id: number
  }
  ready: {
    session_id: number
    frame_id: number
  }
}


/** Host → UI: event name → payload */
export type RedrawPayload = {
  frame_id: number
  columns: number
  rows: number
  wrap: boolean
  full?: boolean
  buffer: {
    id: number
    name: string
    modified?: boolean
  }
  lines: string[]
  cursor: {
    row: number
    column: number
    visible: boolean
  }
  mode: string
}

export type HostEventMap = {
  attached: HostResultMap["attach"]
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
