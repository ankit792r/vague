import { CommandLine } from "./CommandLine";
import { SplitEditor } from "./SplitEditor";
import { StatusLine } from "./StatusLine";

export function VagueFrame() {
  return (
    <div class="vague-frame">
      <SplitEditor />
      <StatusLine />
      <CommandLine />
    </div>
  )
}
