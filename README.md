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

Without the systemd service, the server starts automatically on first `vague`
launch. After install, manage it with `systemctl --user`.

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
| `~/.config/systemd/user/vague-server.service` | User systemd unit |

Ensure `~/.local/bin` is on your `PATH`.

### Systemd server

The install script enables a **user** service so the IPC server starts with your session:

```bash
systemctl --user status vague-server.service
systemctl --user restart vague-server.service
systemctl --user stop vague-server.service
```

Logs: `journalctl --user -u vague-server.service -f`

If the service fails to start during install, run `systemctl --user enable --now vague-server.service` after login. For the server to stay up without an active session, you may need `loginctl enable-linger "$USER"`.

## Uninstall

```bash
./scripts/install.sh uninstall
```

Removes the binary, desktop entry, icon, systemd unit, socket, server log, and
`~/.local/share/vague` if present.
