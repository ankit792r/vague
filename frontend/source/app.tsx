import { useEffect } from "preact/hooks";
import { VagueFrame } from "./comps/VagueFrame";
import { onKeyDown } from "./keymap/keydown";
import { HostInvoke } from "./bridge/bridge";
import type { FrameState } from "./hooks/types";
import { useFrameState } from "./hooks/useFrameState";

export function App() {
	const { setFState } = useFrameState()

	useEffect(() => {
		window.onHostEvent = (state: FrameState) => {
			console.log(state);
		}

		window.addEventListener("keydown", onKeyDown)

		HostInvoke("frame_attached", { args: [], })
			.then(setFState)

		return () => window.removeEventListener("keydown", onKeyDown)
	}, [])

	return (
		<VagueFrame />
	)
}
