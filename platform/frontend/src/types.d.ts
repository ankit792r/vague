export type IpcReply = {
  id: number
  result?: unknown
  error?: string
}

type IpcBindingFn = (
  id: number,
  method: string,
  params: unknown,
) => Promise<IpcReply>


declare global {
  interface Window {
    ipcBinding?: IpcBindingFn
  }
}

