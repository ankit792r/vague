export type FrameState = {
	show_cmd: boolean;
	echo: string;
	width: number;
	height: number;
	work_dir: string;

	root: Window | null;
	active_window_id: number;
};


export type Window = {
	id: number;
	frame_id: number;
	buffer_id: number;
	lines: string[];
	children?: Window[];
};
