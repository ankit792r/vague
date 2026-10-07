import { createContext, type ComponentChildren } from 'preact';
import { useState, useContext } from 'preact/hooks';
import type { FrameState } from "./types";
import type { SetStateAction } from 'preact/compat';

interface FrameStateContextType {
	fState: FrameState | null;
	setFState: React.Dispatch<SetStateAction<FrameState | null>>;
	updateState: (state: Partial<FrameState>) => void;
}

const FrameStateContext = createContext<FrameStateContextType | undefined>(undefined);

export function FrameStateProvider({ children }: { children: ComponentChildren }) {
	const [fState, setFState] = useState<FrameState | null>(null);

	const updateState = (updates: Partial<FrameState>) => {
		setFState((prev) => {
			if (prev === null) {
				return null;
			}

			return {
				...prev,
				...updates,
			};
		});
	};

	return (
		<FrameStateContext.Provider value={{ fState, setFState, updateState }}>
			{children}
		</FrameStateContext.Provider>
	);
}

export function useFrameState() {
	const context = useContext(FrameStateContext);
	if (context === undefined) {
		throw new Error('useFrameState must be used within a FrameStateProvider');
	}
	return context;
}

