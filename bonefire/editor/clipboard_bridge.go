package editor

import "vague/bonefire/clipboard"

func clipboardReadPrimary() (string, bool) {
	return clipboard.ReadPrimary()
}

func clipboardReadClipboard() (string, bool) {
	return clipboard.ReadClipboard()
}

func clipboardWritePrimary(text string) {
	clipboard.WritePrimary(text)
}

func clipboardWriteClipboard(text string) {
	clipboard.WriteClipboard(text)
}

func clipboardOSC52(data []byte) {
	_ = clipboard.WriteOSC52(data)
}
