package editor

// consumeCountDigit accumulates 1-9 and 0 when a count is already started.
// Returns true if the key was consumed. Bare "0" is not a count digit (line start motion).
func (e *Editor) consumeCountDigit(keys string) bool {
	if len(keys) != 1 || keys[0] < '0' || keys[0] > '9' {
		return false
	}

	if keys == "0" && e.pendingCount == 0 {
		return false
	}

	d := int(keys[0] - '0')
	if e.pendingCount == 0 {
		e.pendingCount = d
	} else {
		e.pendingCount = e.pendingCount*10 + d
	}

	return true
}

func (e *Editor) takeCount() int {
	if e.pendingCount <= 0 {
		return 1
	}

	c := e.pendingCount
	e.pendingCount = 0
	return c
}

// takeCountOptional returns a count and whether digits were typed before this command.
func (e *Editor) takeCountOptional() (count int, explicit bool) {
	if e.pendingCount <= 0 {
		return 1, false
	}

	return e.takeCount(), true
}

func (e *Editor) clearPendingCount() {
	e.pendingCount = 0
}
