import { useEffect } from "preact/hooks";
import { VagueFrame } from "./comps/VagueFrame";
import { onKeyDown } from "./keymap/keydown";
import { HostInvoke } from "./bridge/bridge";

export function App() {
  useEffect(() => {
    window.onHostEvent = (state: string) => {
      console.log(state);
    }

    window.addEventListener("keydown", onKeyDown)

    HostInvoke("frame_attached", {
      args: [],
      bang: true,
      count: 10
    }).then((update) => {
      console.log("got ui attach update: ", update)
    })

    return () => window.removeEventListener("keydown", onKeyDown)
  }, [])

  return (
    <VagueFrame />
  )
}
