import type { FrameState } from "../hooks/types";

export type HostMethodParams = {
	args?: string[],
	bang?: boolean,
	count?: number,
}

declare global {
	interface Window {
		// host will only send state update, and client only have to update the UI
		// FIXME: state: State
		onHostEvent?: (state: FrameState) => void


		// FIXME: return of this will also give partial status update
		hostInvoke?: (method: string, params: HostMethodParams) => Promise<FrameState>
	}
}

export async function HostInvoke(
	method: string,
	params: HostMethodParams
): Promise<FrameState> { // FIXME: change unknow to proper state type
	const invoke = window.hostInvoke;
	if (!invoke)
		return Promise.reject(new Error("host bridge unavailable"))

	return await invoke(method, params)
}
