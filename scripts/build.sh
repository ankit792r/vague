#!/usr/bin/env bash

set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

BIN="${ROOT}/bin/vague"
GO="${GO:-go}"
BUN="${BUN:-bun}"

install_deps() {
  local missing=()
  local runtime_missing=()

  command -v "$GO" >/dev/null 2>&1 || missing+=(go)
  command -v "$BUN" >/dev/null 2>&1 || missing+=(bun)

  if command -v pacman >/dev/null 2>&1; then
    pacman -Q webkit2gtk-4.1 >/dev/null 2>&1 || runtime_missing+=(webkit2gtk-4.1)
    pacman -Q gtk3 >/dev/null 2>&1 || runtime_missing+=(gtk3)
  fi

  if [ ${#missing[@]} -eq 0 ] && [ ${#runtime_missing[@]} -eq 0 ]; then
    return 0
  fi

  if ! command -v pacman >/dev/null 2>&1; then
    [ ${#missing[@]} -gt 0 ] && {
      echo "error: missing build tools: ${missing[*]}" >&2
      echo "Install Go and Bun, then re-run scripts/build.sh" >&2
      exit 1
    }
    [ ${#runtime_missing[@]} -gt 0 ] && {
      echo "warning: missing runtime libraries: ${runtime_missing[*]}" >&2
      echo "The binary will build, but vague needs WebKit/GTK to run the GUI." >&2
    }
    return 0
  fi

  local to_install=()
  to_install+=("${missing[@]}")
  to_install+=("${runtime_missing[@]}")

  if [ ${#to_install[@]} -eq 0 ]; then
    return 0
  fi

  echo "Installing packages: ${to_install[*]}"
  sudo pacman -S --needed --noconfirm "${to_install[@]}"
}

build_frontend() {
  echo "Building frontend..."
  cd "$ROOT/frontend"
  "$BUN" install
  "$BUN" run build
  cd "$ROOT"
}

build_binary() {
  echo "Building vague binary..."
  mkdir -p "$ROOT/bin"

  # Pure Go binary with embedded UI + libwebview.so (WebKit still required at runtime).
  CGO_ENABLED=0 "$GO" build \
    -ldflags "-s -w" \
    -trimpath \
    -o "$BIN" \
    .

  echo "Built $BIN"
  file "$BIN" || true
}

main() {
  install_deps
  build_frontend
  build_binary
  echo "Done. Run: ./bin/vague"
}

main "$@"
