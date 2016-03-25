# Forward SSH

Usage:
```bash
    go-forward-ssh (-L | -R) <host> <local-port> <remote-port> [--ssh-key ssh-key] [--ssh-user ssh-user]
    go-forward-ssh -h | --help
    go-forward-ssh --version
```

- Use `-L` for local forwarding
- Use `-R` for remote forwarding

Build:
```bash
go get github.com/segmentio/godep
make build
```

Run:
```bash
./bin/forward-ssh <options>
```