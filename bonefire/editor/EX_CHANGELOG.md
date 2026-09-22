# Ex command rollout log

Each entry matches a Phase 2 Ex TODO commit.

1. `:substitute` with literal pattern and `g` `i` `I` flags (`c` accepted, no UI confirm).
2. Ranges: `%`, line numbers, `'<,'>` visual last selection.
3. `:global` / `:v` with pattern + subcommand per line.
4. `:normal` / `:norm` runs normal-mode keys on a range.
5. `:read` inserts a file at cursor or range start.
6. `:write >>` appends range lines to a file.
7. `:edit` no-arg reload; `:find` / `:sf` glob in frame workdir.
8. `:bdelete` / `:bd` / `:bwipeout` confirm via modified check and `!`.
9. `:cd` and `:pwd` per frame working directory.
10. `:set` with `?` `&` stubs and wrap/number toggles.
11. `:map` / `:nmap` / `:imap` / `:vmap` in-memory bindings.
12. `:command` user ex (stub error).
13. `:source` runs ex lines from a file.
14. `:reg` shows unnamed register.
15. `:marks` / `:delm` in-memory marks.
16. `:jumps` / `:clearjumps` jump list echo/clear.
17. `:undo` / `:redo` / `:undolist` stub.
18. `:checktime` stub echo.
19. Diff commands stub error.
20. `:terminal` stub error.
