import { render } from 'preact'
import './index.css'
import { App } from './app.tsx'
import { FrameStateProvider } from './hooks/useFrameState.tsx'

render(
	<FrameStateProvider>
		<App />
	</FrameStateProvider>
	, document.getElementById('app')!)
