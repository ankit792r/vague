#!/usr/bin/env bash
# Install or uninstall vague for the current user (Linux).
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"

PREFIX="${PREFIX:-$HOME/.local}"
BINDIR="${BINDIR:-$PREFIX/bin}"
DATADIR="${DATADIR:-$PREFIX/share}"
APPDIR="${APPDIR:-$DATADIR/applications}"
ICONDIR="${ICONDIR:-$DATADIR/icons/hicolor/256x256/apps}"
APPNAME="vague"
BIN="${ROOT}/bin/${APPNAME}"
DESKTOP_SRC="${ROOT}/vague.desktop"
ICON_SRC="${ROOT}/logo.png"

usage() {
  cat <<EOF
Usage: $(basename "$0") <install|uninstall>

  install    Build if needed, install binary, desktop entry, and icon
  uninstall  Remove binary, desktop entry, icon, and runtime data

Environment:
  PREFIX     Install root (default: ~/.local)
EOF
}

require_build() {
  if [ ! -x "$BIN" ]; then
    echo "Binary not found. Running build..."
    "$ROOT/scripts/build.sh"
  fi
}

install_app() {
  require_build

  if [ ! -f "$DESKTOP_SRC" ]; then
    echo "error: missing $DESKTOP_SRC" >&2
    exit 1
  fi

  if [ ! -f "$ICON_SRC" ]; then
    echo "error: missing $ICON_SRC" >&2
    exit 1
  fi

  install -d "$BINDIR" "$APPDIR" "$ICONDIR"
  install -m 755 "$BIN" "$BINDIR/$APPNAME"
  install -m 644 "$DESKTOP_SRC" "$APPDIR/$APPNAME.desktop"
  install -m 644 "$ICON_SRC" "$ICONDIR/$APPNAME.png"

  if command -v update-desktop-database >/dev/null 2>&1; then
    update-desktop-database "$APPDIR" 2>/dev/null || true
  fi

  if command -v gtk-update-icon-cache >/dev/null 2>&1; then
    gtk-update-icon-cache -f -t "$DATADIR/icons/hicolor" 2>/dev/null || true
  fi

  cat <<EOF
Installed vague:
  binary:  $BINDIR/$APPNAME
  desktop: $APPDIR/$APPNAME.desktop
  icon:    $ICONDIR/$APPNAME.png

Ensure $BINDIR is on your PATH.
Launch from the app menu or run: vague
EOF
}

remove_runtime_data() {
  local runtime_dir="${XDG_RUNTIME_DIR:-/run/user/$(id -u)}"

  rm -f "$runtime_dir/vague.sock"
  rm -f "$runtime_dir/vague.sock.lock"
  rm -f "$runtime_dir/vague-server.log"
  rm -rf "${XDG_DATA_HOME:-$HOME/.local/share}/vague"
}

uninstall_app() {
  rm -f "$BINDIR/$APPNAME"
  rm -f "$APPDIR/$APPNAME.desktop"
  rm -f "$ICONDIR/$APPNAME.png"

  remove_runtime_data

  if command -v update-desktop-database >/dev/null 2>&1; then
    update-desktop-database "$APPDIR" 2>/dev/null || true
  fi

  if command -v gtk-update-icon-cache >/dev/null 2>&1; then
    gtk-update-icon-cache -f -t "$DATADIR/icons/hicolor" 2>/dev/null || true
  fi

  echo "Removed vague from $PREFIX and cleared runtime data."
}

case "${1:-}" in
  install)
    install_app
    ;;
  uninstall)
    uninstall_app
    ;;
  *)
    usage
    exit 1
    ;;
esac
