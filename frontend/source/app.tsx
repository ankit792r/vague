import { useEffect } from "preact/hooks";
import { VagueFrame } from "./comps/VagueFrame";
import { onKeyDown } from "./keymap/keydown";
import { HostInvoke } from "./bridge/bridge";
import { useFrameState } from "./hooks/useFrameState";
import useWindowSize from "./hooks/useWindowSize";

export function App() {
  const { setFState, updateState } = useFrameState()
  const { width, height } = useWindowSize();

  useEffect(()=> {
    HostInvoke("frame_updated", { args: { width: width, height: height }, })
      .then(setFState)
  }, [width, height])

  useEffect(() => {
    window.onHostEvent = updateState

    // window.onHostEvent = (state: FrameState) => {
    // 	console.log(state);
    // }

    window.addEventListener("keydown", onKeyDown)
    HostInvoke("frame_attached", { args: { width: width, height: height }, })
      .then(setFState)

    return () => window.removeEventListener("keydown", onKeyDown)
  }, [])

  return (
    <VagueFrame />
  )
}
