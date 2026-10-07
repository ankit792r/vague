import { useEffect } from "preact/hooks";
import { VagueFrame } from "./comps/VagueFrame";
import { onKeyDown } from "./keymap/keydown";
import { HostInvoke } from "./bridge/bridge";
import type { FrameState } from "./hooks/types";
import { FrameStateProvider } from "./hooks/useFrameState";

export function App() {
	useEffect(() => {
		window.onHostEvent = (state: FrameState) => {
			console.log(state);
		}

		window.addEventListener("keydown", onKeyDown)

		HostInvoke("frame_attached", { args: [], })
			.then((update) => {
				console.log("got ui attach update: ", update)
			})

		return () => window.removeEventListener("keydown", onKeyDown)
	}, [])

	return (
		<FrameStateProvider>
			<VagueFrame />
		</FrameStateProvider>
	)
}
