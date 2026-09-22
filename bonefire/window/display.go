package window

// ListChars controls :set listchars rendering.
type ListChars struct {
	Tab     string
	Trail   string
	Eol     string
	Extends string
}

func DefaultListChars() ListChars {
	return ListChars{Tab: ">", Trail: "-", Eol: "$", Extends: "-"}
}

// DisplayOpts are per-window :setlocal display options.
type DisplayOpts struct {
	RelativeNumber bool
	SideScrollOff  int
	SmoothScroll   bool

	List      bool
	ListChars ListChars

	CursorLine   bool
	CursorColumn bool

	ColorColumn []int
	TextWidth   int
	WrapMargin  int

	SignColumn bool

	ShowMode bool
	ShowCmd  bool
	Ruler    bool

	StatusLine string
}

func DefaultDisplayOpts() DisplayOpts {
	return DisplayOpts{
		ListChars:  DefaultListChars(),
		ShowMode:   true,
		ShowCmd:    true,
		Ruler:      true,
		StatusLine: "%f %m %l,%c",
	}
}
