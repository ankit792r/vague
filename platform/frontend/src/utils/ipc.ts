import type { IpcReply } from "../types"

let nextId = 1

export function nextCallId(): number {
  return nextId++
}

export function callIpc(
  id: number,
  method: string,
  params: unknown = null,
): Promise<IpcReply> {
  const ipc = window.ipcBinding
  if (typeof ipc !== "function") {
    return Promise.reject(new Error("ipc bridge is not available"))
  }
  return ipc(id, method, params)
}

