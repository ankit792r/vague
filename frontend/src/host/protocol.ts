/** UI → host: method name → params */
export type HostRequestMap = {
  input: { keys: string }
  execute: {
    name: string
    args?: string[]
    bang?: boolean
    count?: number
  }
  ready: Record<string, never>
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
export type HostEventMap = {
  attached: HostResultMap["attach"]
  redraw: unknown // replace with RedrawPayload when defined in Go
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
