# Vague

A Vim-inspired text editor with an Emacs-style window model.

## Running

```bash
make build
./bin/vague              # starts the server automatically if needed
./bin/vague hello.py     # open a file
```

The server runs in the background on first launch. Logs go to
`$XDG_RUNTIME_DIR/vague-server.log`. To run it in the foreground manually:

```bash
vague server-start
```
