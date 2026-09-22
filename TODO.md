# Vague → Neovim parity

Checklist toward a proper Neovim-class editor (Vim core + Neovim UX). One line per item; order is rough priority within each section.

Legend: **done** items live under [Shipped](#shipped). Everything else is open.

---

## Shipped

- Normal / insert / visual char / visual line / visual block; numeric counts; `.` repeat (operators, text objects, indent, case, format, numbers)
- Motions: `h/j/k/l`, `w/b/e`, `0/$`/`^`/`_`/`|`, `gg`/`G`, `f`/`F`/`t`/`T`/`;`/`,`, line `j/k` with wrap, `o`/`O`, `i`/`a`/`A`/`I`
- Operators: `d`/`c`/`y` + linewise, text objects, `gu`/`gU`/`g~`, `gq`, `>>`/`<<`/`==`, visual `>`/`<`, `s`/`S`/`C`, `r`/`R`, `~`, `Ctrl-a`/`Ctrl-x`
- Unnamed register; `p`/`P`; `u` / `<C-r>`
- `/` `?` `n` `N`; current-match highlight on redraw
- Open / save / quit; `:e` `:e!` `:w` `:q` `:wq` `:x`; buffers `:b` `:bn` `:bp` `:buffers`
- `:goto` / `:123`; `:wrap` `:number` and `:set` aliases; `:only` placeholder echo
- Viewport scroll, soft wrap, optional line numbers, status line Ln/Col, echo area

---

## Phase 1 — Vim editing core

### Motions still missing

- [x] `f` `F` `t` `T` on line; `;` `,` repeat
- [x] `^` `_` first non-blank; `|` column
- [x] `{` `}` `(` `)` paragraph / sentence / block (simple heuristics OK first)
- [x] `%` matching bracket; `[` `]` prev/next section
- [x] `*` `#` search word forward/back; `g*` `g#`
- [x] `H` `M` `L` screen lines; `zt` `zz` `zb` scroll cursor line
- [x] `Ctrl-o` / `Ctrl-i` jump list
- [x] `ge` `gE` end of word variants
- [x] `Ctrl-f` `Ctrl-b` `Ctrl-d` `Ctrl-u` page scroll
- [x] `g0` `g$` `gm` screen line boundaries (when wrapped)

### Operators & changes

- [x] `s`/`S`/`cc`/`C` parity (char vs line change)
- [x] `r` replace char; `R` replace mode
- [x] `~` toggle case
- [x] `gu`/`gU`/`g~` operator-pending (minimal)
- [x] `>>` `<<` `==` indent lines (spaces/tabs, `:set shiftwidth`)
- [x] `>` `<` visual/operator indent
- [x] `!` external filter stub (echo; job runner later)
- [x] `gq` format text (wrap width)
- [x] `Ctrl-a` / `Ctrl-x` number increment/decrement
- [x] `Ctrl-v` visual block mode
- [x] `gv` reselect last visual
- [x] `Ctrl-r` register insert in insert mode (not just redo)
- [x] `Q`/`gQ` ex linewise (low priority)
- [x] `.` records more motion/operator combos; fix known gaps

### Text objects (operator + visual)

- [x] `iw` `aw` `iW` `aW`
- [x] `i"` `a"` `i'` `a'` backtick-quoted strings
- [x] `i(` `a(` `ib` `ab` and `i[` `a[` `i{` `a{` pairs
- [x] `it` `at` tag block (when HTML-ish)
- [x] `ip` `ap` paragraph; [x] `is` `as` sentence

---

## Phase 2 — Command-line & ex

### Command-line mode

- [x] **incsearch** — preview match while typing `/` `?`
- [x] Command-line history (`/` `?` `:`)
- [x] `Ctrl-c` / `<Esc>` clears highlight on cancel search
- [x] `Ctrl-r` insert register in cmdline
- [x] Wildmenu / tab completion (`:e` paths, `:b` names, `:set` keys)
- [ ] `:help` stub or generated key reference

### Ex commands

- [ ] `:substitute` `:s` ranges + flags `g` `c` `i` `I`
- [ ] `:%s` whole buffer; `:5,10s` range; `:'<,'>s` visual range
- [ ] `:global` `:g` / `:v`
- [ ] `:normal` `:norm` execute keys on range
- [ ] `:read` `:r` insert file or output
- [ ] `:write` ranges / `:w >>` append
- [ ] `:edit` no-arg reload; `:find` `:sf` path search
- [ ] `:buffer` `:bd` `:bdelete` `:bwipeout` with confirm
- [ ] `:cd` `:pwd` working directory per tab/frame
- [ ] `:set` / `:setlocal` full option parser (`set opt?` `set opt&` `opt+=`)
- [ ] `:map` `:nmap` `:imap` `:vmap` user key bindings
- [ ] `:command` user ex commands
- [ ] `:source` load config script
- [ ] `:reg` show registers
- [ ] `:marks` `:delm` marks
- [ ] `:jumps` `:clearjumps`
- [ ] `:undo` `:undolist` `:later` `:earlier` time travel (stretch)
- [ ] `:checktime` autoread prompt
- [ ] `:diffsplit` `:diffoff` `:diffget` `:diffput`
- [ ] `:terminal` `:term` (see Phase 14)

---

## Phase 3 — Search & pattern

- [ ] Vim regex engine (`\m` `\M` `\v` `\V`) or documented subset
- [ ] `:set ignorecase` `smartcase` `hlsearch` `incsearch` `wrapscan`
- [ ] Highlight **all** matches (`hlsearch`) + `:nohlsearch`
- [ ] Search offset `:ta` tags (later)
- [ ] `:vimgrep` / quickfix list (see Phase 12)
- [ ] `:sort` `:uniq` on ranges

---

## Phase 4 — Registers & clipboard

- [ ] Named registers `a`–`z` and `"A` append
- [ ] `"0` last yank; `"1`–`"9` delete ring
- [ ] `"-` small delete; `":` `.` last command; `"/` last search
- [ ] `"*` `"+` system clipboard (X11/Wayland/macOS/Windows)
- [ ] `"_` black hole; `"=` expression register (later)
- [ ] `:reg` display; getreg/setreg API for plugins
- [ ] OSC 52 for SSH/remote clipboard (optional)

---

## Phase 5 — Marks & digraphs

- [ ] Local marks `a`–`z`; file marks `A`–`Z` (persist in shada)
- [ ] `` ` `` / `'` jumps to mark; `` ' `` line jump; `` `. `` last change; `` `" `` last jump
- [ ] `:marks`; mark in status when set
- [ ] Digraphs `Ctrl-k` or `:digraphs`
- [ ] Keyboard locale / keycode normalization (fix gaps vs Neovim)

---

## Phase 6 — Windows & tab pages

- [ ] `:split` `:vsplit` + `:new` `:vnew`
- [ ] `:only` `:close` `:hide` `:wincmd`
- [ ] `<C-w>` + `h/j/k/l/H/J/K/L` focus; `w/W` cycle
- [ ] `<C-w>` + `+` `-` `<` `>` `=` resize; `_` `|` maximize
- [ ] Multi-window layout in webview (splits + active border)
- [ ] Per-window options (`:setlocal` wrap/number/scroll)
- [ ] Tab pages `:tabnew` `:tabclose` `:tabnext` tabline UI
- [ ] `:args` `:argadd` `:argdo` argument list

---

## Phase 7 — Buffers & files

- [ ] Modified-buffer confirm on `:b` `:bn` `:bp` `:e` `:q` (`:confirm` `:set confirm`)
- [ ] `:wa` `:wqall` `:qa` `:xall`
- [ ] Autoread when file changed on disk
- [ ] `:set autowrite` `autoread` `backup` `writebackup` `swapfile`
- [ ] Hidden buffers `:set hidden`
- [ ] `:lcd` local cwd; modeline `vim:` lines (optional)
- [ ] Binary `:e ++binary` `:xxd` (optional)
- [ ] Fuzzy find `:Telescope`-class UI (picker) — product decision

---

## Phase 8 — `:set` & display

- [ ] `relativenumber` / `n relativenumber`
- [ ] `scrolloff` `sidescrolloff` `smoothscroll` (wire existing ScrollOff)
- [ ] `list` `listchars` tabs/trail/eol
- [ ] `cursorline` `cursorcolumn`
- [ ] `colorcolumn` `textwidth` `wrapmargin`
- [ ] `signcolumn` yes/auto
- [ ] `statusline` `tabline` customizable (Lua-style or config DSL)
- [ ] `showmode` `showcmd` `ruler` (showcmd for pending operator)
- [ ] `termguicolors` / theme hooks; multiple built-in colorschemes
- [ ] `conceallevel` (syntax-dependent, later)
- [ ] `foldenable` folds (Phase 11)

---

## Phase 9 — Insert & replace mode

- [ ] `<C-w>` word back; `<C-u>` to start; `<C-k>` digraph
- [ ] `<C-a>` / `<C-e>` line start/end
- [ ] `<C-o>` normal one-shot in insert
- [ ] `<C-r>` insert register
- [ ] `<C-t>` / `<C-d>` indent in insert
- [ ] `<C-x>` completion submodes (`Ctrl-n/p`, omni — ties to LSP)
- [ ] Replace mode `R` and `gr` visual replace
- [ ] Insert paste `:set paste` `nopaste`

---

## Phase 10 — Undo & change history

- [ ] Undo branches visible `:undolist`
- [ ] Persistent undo across restarts
- [ ] `:earlier` `:later` time-based undo
- [ ] `:redo` ex command (if not same as `<C-r>`)
- [ ] `:changes` change list `` `. `` `` `" ``

---

## Phase 11 — Syntax, highlight, folds

- [ ] Filetype detection (`ft`, modeline, extension map)
- [ ] Basic syntax highlighting (regex rules per ft)
- [ ] Treesitter integration (Neovim-style `:TSInstall`) — long-term
- [ ] `:syntax` on/off; highlight groups linked to theme
- [ ] Conceal, spell (`:set spell`), spelllang
- [ ] Folds: manual/indent/syntax; `za` `zo` `zc`; `:foldopen` `:foldclose`

---

## Phase 12 — LSP, diagnostics, quickfix (Neovim core)

- [ ] LSP client: attach, notify, request/response JSON-RPC
- [ ] Diagnostics virtual text + sign column + `:lua vim.diagnostic` equivalent
- [ ] `gd` `gr` `K` hover, definition, references UI
- [ ] Completion menu (`completeopt` `pum` `complete` sources)
- [ ] Code actions lightbulb / `:lua vim.lsp.buf.code_action`
- [ ] Formatting `gq` / `:Format` / LSP format
- [ ] Quickfix / loclist windows; `:cnext` `:cp` `:cf` `:make`
- [ ] Outline / document symbol sidebar (optional)

---

## Phase 13 — Terminal & jobs

- [ ] `:terminal` buffer type; job control API
- [ ] Send text to terminal job; `:termopen`
- [ ] `Ctrl-\` `Ctrl-n` leave terminal mode
- [ ] Background jobs + `jobwait` equivalent for `:make` `!`

---

## Phase 14 — UI toolkit (Neovim floats/popups)

- [ ] Floating windows (centered popups, borders, `z-index`)
- [ ] `:popup` menu at cursor
- [ ] Notification system (`vim.notify` level + history)
- [ ] Message history (`:messages`)
- [ ] Wildmenu rendering; cmdline completion menu
- [ ] Multiline message area (avoid truncating `:buffers`)
- [ ] Mouse: click position, drag selection, scroll wheel
- [ ] Cursor blink; block/bar/shape per mode `:set guicursor`

---

## Phase 15 — Config, autocmd, API

- [ ] Config file (`init.vim` / `init.lua` or Vague-native `init.toml`)
- [ ] `:autocmd` FileType, BufWrite, VimEnter, etc.
- [ ] `:let` options variables (if scripting added)
- [ ] Stable embedding API (Go) for automation without forking
- [ ] RPC / socket protocol for external tools (Neovim msgpack RPC analogue)
- [ ] Export editor state: buffers, windows, mode, options

---

## Phase 16 — Plugins & extension

- [ ] Plugin loader (path, lazy load)
- [ ] User commands + keymaps registration API
- [ ] User autocommands; highlight namespace IDs
- [ ] Decoration provider (virtual text, extmarks on lines)
- [ ] Tree-sitter / LSP as optional built-in plugins
- [ ] Package manager story (lazy.nvim equivalent — far future)

---

## Phase 17 — Session, shada, remote

- [ ] ShaDa / viminfo: files, registers, history, marks persist
- [ ] Session `:mksession` `:source Session.vim`
- [ ] Remote editing single instance / server multiplex clients
- [ ] `--listen` headless server mode (already partial via backbone)

---

## Phase 18 — Quality, docs, release

- [ ] CI: `go test ./...`, frontend build, lint
- [ ] Fix `dispatch_test` socket lock for CI envs
- [ ] ROADMAP.md link to this file; changelog
- [ ] Manual test checklist for releases
- [ ] Large-file performance (incremental layout, partial redraw)
- [ ] Fuzz input / property tests on text layer
- [ ] Accessibility audit (ARIA, keyboard-only)
- [ ] Packaging: Flatpak/AppImage beyond install.sh

---

## Reference — Neovim docs map

Use official Neovim `:help` indices as the completeness bar: `motion.txt`, `change.txt`, `visual.txt`, `cmdline.txt`, `windows.txt`, `tabpage.txt`, `options.txt`, `lsp.txt`, `treesitter.txt`, `terminal.txt`, `api.txt`. Vague does not need every obsolete Vim feature, but should not surprise a Neovim user on everyday editing, ex, windows, search, registers, and `:set`.
