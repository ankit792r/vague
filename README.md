# Vague

A Vim-inspired text editor with an Emacs-style window model.

## Requirements (Linux)

- Go 1.27+
- [Bun](https://bun.sh)
- **Runtime:** `webkit2gtk-4.1`, `gtk3`

On Arch Linux:

```bash
sudo pacman -S go bun webkit2gtk-4.1 gtk3
```

## Build

```bash
./scripts/build.sh
./bin/vague
./bin/vague hello.py
```

The build script installs missing Arch packages (when `pacman` is available),
builds the frontend, and writes a standalone Go binary to `bin/vague`.

The server starts automatically on first launch. Logs:
`$XDG_RUNTIME_DIR/vague-server.log`.

## Install

```bash
./scripts/install.sh install
```

Installs to `~/.local`:

| Path | Purpose |
|------|---------|
| `~/.local/bin/vague` | Executable |
| `~/.local/share/applications/vague.desktop` | App menu entry |
| `~/.local/share/icons/hicolor/256x256/apps/vague.png` | Icon (`logo.png`) |

Ensure `~/.local/bin` is on your `PATH`.

## Uninstall

```bash
./scripts/install.sh uninstall
```

Removes the binary, desktop entry, icon, socket, server log, and
`~/.local/share/vague` if present.
