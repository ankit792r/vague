# Search patterns (subset)

Vague implements a **documented subset** of Vim search patterns, not full regex.

## Magic prefixes

| Prefix | Mode | Meaning |
|--------|------|---------|
| `\m` | magic (default) | `.` `^` `$` `[` `]` `\` `~` are special; other bytes are literal |
| `\M` | nomagic | only `^` and `$` are special |
| `\v` | very magic | `.` `^` `$` `[` `]` `\` and common regex metacharacters are special |
| `\V` | very nomagic | all bytes literal (except `\` escapes) |

Additional: `\c` force ignore case, `\C` force match case.

## Features

- Literal fast path when no special characters apply
- Character classes `[abc]` and `[^abc]` (no nested classes)
- Escapes `\n` `\t` `\.` etc.
- `.` does not match newline
- Search uses `:set ignorecase` / `smartcase` and `\c` / `\C`

## Not yet supported

- Captures, backreferences, `\{n,m}` counts
- `\|` alternation (except partial in `\v` translation)
- `\zs` `\ze`, lookaround, `\%[abc]`
- `\*` `\?` quantifiers in `\m` ( `*` is literal in `\m` )

Use `:help pattern` in Vim for the full reference; this file tracks Vague behavior.

## Search offsets

Trailing offsets on `/` and `?` patterns are supported: `+N`, `-N`, `e`, `s`, `b`, and forms like `e+1` (see `search_offset.go`).
