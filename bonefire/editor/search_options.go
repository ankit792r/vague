package editor

import "fmt"

// SearchOpts mirrors common Vim :set search options (per-editor defaults).
type SearchOpts struct {
	IgnoreCase bool
	SmartCase  bool
	HlSearch   bool
	IncSearch  bool
	WrapScan   bool
	Magic      MagicMode
}

func DefaultSearchOpts() SearchOpts {
	return SearchOpts{
		WrapScan:  true,
		IncSearch: true,
		Magic:     MagicMagic,
	}
}

func (e *Editor) ApplySetSearchOption(name string, enable bool) (string, error) {
	switch name {
	case "ignorecase", "ic":
		e.searchOpts.IgnoreCase = enable
		if enable {
			return "ignorecase", nil
		}
		return "noignorecase", nil
	case "smartcase", "scs":
		e.searchOpts.SmartCase = enable
		if enable {
			return "smartcase", nil
		}
		return "nosmartcase", nil
	case "hlsearch", "hls":
		e.searchOpts.HlSearch = enable
		if enable {
			return "hlsearch", nil
		}
		return "nohlsearch", nil
	case "incsearch":
		e.searchOpts.IncSearch = enable
		if enable {
			return "incsearch", nil
		}
		return "noincsearch", nil
	case "wrapscan", "ws":
		e.searchOpts.WrapScan = enable
		if enable {
			return "wrapscan", nil
		}
		return "nowrapscan", nil
	default:
		return "", fmt.Errorf("Unknown option: %s", name)
	}
}

func (e *Editor) SearchOptionQuery(name string) string {
	switch name {
	case "ignorecase", "ic":
		if e.searchOpts.IgnoreCase {
			return "ignorecase"
		}
		return "noignorecase"
	case "smartcase", "scs":
		if e.searchOpts.SmartCase {
			return "smartcase"
		}
		return "nosmartcase"
	case "hlsearch", "hls":
		if e.searchOpts.HlSearch {
			return "hlsearch"
		}
		return "nohlsearch"
	case "incsearch":
		if e.searchOpts.IncSearch {
			return "incsearch"
		}
		return "noincsearch"
	case "wrapscan", "ws":
		if e.searchOpts.WrapScan {
			return "wrapscan"
		}
		return "nowrapscan"
	default:
		return ""
	}
}
