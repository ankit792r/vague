package editor

import "fmt"

// FileOpts are global buffer/file behavior flags.
type FileOpts struct {
	Confirm     bool
	Autowrite   bool
	Autoread    bool
	Hidden      bool
	Backup      bool
	WriteBackup bool
	Swapfile    bool
	Paste       bool
}

func defaultFileOpts() FileOpts {
	return FileOpts{
		Backup:      true,
		WriteBackup: true,
	}
}

func (e *Editor) FileOptionQuery(name string) string {
	switch name {
	case "confirm":
		return boolOpt(e.fileOpts.Confirm)
	case "autowrite", "aw":
		return boolOpt(e.fileOpts.Autowrite)
	case "autoread", "ar":
		return boolOpt(e.fileOpts.Autoread)
	case "hidden", "hid":
		return boolOpt(e.fileOpts.Hidden)
	case "backup", "bk":
		return boolOpt(e.fileOpts.Backup)
	case "writebackup", "wb":
		return boolOpt(e.fileOpts.WriteBackup)
	case "swapfile", "swf":
		return boolOpt(e.fileOpts.Swapfile)
	case "paste":
		return boolOpt(e.fileOpts.Paste)
	default:
		return ""
	}
}

func (e *Editor) ApplySetFileOption(name string, enable bool) (string, error) {
	switch name {
	case "confirm":
		e.fileOpts.Confirm = enable
		return "confirm", nil
	case "autowrite", "aw":
		e.fileOpts.Autowrite = enable
		return "autowrite", nil
	case "autoread", "ar":
		e.fileOpts.Autoread = enable
		return "autoread", nil
	case "hidden", "hid":
		e.fileOpts.Hidden = enable
		return "hidden", nil
	case "backup", "bk":
		e.fileOpts.Backup = enable
		return "backup", nil
	case "writebackup", "wb":
		e.fileOpts.WriteBackup = enable
		return "writebackup", nil
	case "swapfile", "swf":
		e.fileOpts.Swapfile = enable
		if enable {
			return "", fmt.Errorf("swapfile not implemented")
		}
		return "swapfile", nil
	case "paste":
		e.fileOpts.Paste = enable
		return "paste", nil
	default:
		return "", fmt.Errorf("Unknown option: %s", name)
	}
}

func boolOpt(v bool) string {
	if v {
		return "set"
	}
	return "noset"
}
