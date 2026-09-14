import type {
  HostEvent,
  HostEventMap,
  HostEventMethod,
  HostReply,
  HostRequestMap,
  HostRequestMethod,
  HostResultMap,
} from "./protocol"

declare global {
  interface Window {
    hostRequest?: <M extends HostRequestMethod>(
      id: number,
      method: M,
      params: HostRequestMap[M],
    ) => Promise<HostReply<M>>

    hostEvent?: {
      subscribe: (listener: (msg: HostEvent) => void) => () => void
    }
  }
}

let nextId = 1

export async function hostRequest<M extends HostRequestMethod>(
  method: M,
  params: HostRequestMap[M],
): Promise<HostResultMap[M]> {
  const invoke = window.hostRequest
  if (!invoke) {
    return Promise.reject(new Error("host bridge unavailable"))
  }

  const id = nextId++
  const reply = await invoke(id, method, params)
  if (reply.error) {
    throw new Error(reply.error)
  }
  return reply.result as HostResultMap[M]
}

export function onHostEvent<E extends HostEventMethod>(
  event: E,
  handler: (payload: HostEventMap[E]) => void,
): () => void {
  const bus = window.hostEvent
  if (!bus) {
    throw new Error("host event bus unavailable")
  }

  return bus.subscribe((msg) => {
    if (msg.event === event) {
      handler(msg.payload as HostEventMap[E])
    }
  })
}
