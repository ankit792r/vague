import { useFrameState } from "../hooks/useFrameState"

export function EditorArea() {
	const { fState } = useFrameState()

	return (
		<div class="editor-area" aria-label="editor">
			{fState?.root?.lines.map(line => <p> {line}</p>)}
		</div>
	)
}
