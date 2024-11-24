# sshx

[![Go Reference](https://pkg.go.dev/badge/github.com/matthewmueller/sshx.svg)](https://pkg.go.dev/github.com/matthewmueller/sshx)

Tiny wrapper around [golang.org/x/crypto/ssh](https://pkg.go.dev/golang.org/x/crypto/ssh) that has better defaults. Returns an `*golang.org/x/crypto/ssh.Client` that can be used elsewhere

The goal being it works exactly like if you did `ssh user@host` on your machine.

## Features

- Handles `~/.ssh/known_hosts` on OSX thanks to [skeema/knownhosts](github.com/skeema/knownhosts).
- Uses the active SSH agent on your machine if there is one, allowing you to seamlessly connect without providing a private key (and often the password needed to decrypt that private key).
- Adds `Run(ssh, cmd) (stdout, error)` and `Exec(ssh, cmd) error` commands.
- Allocate an interactive shell with `sshx.Shell(sshClient, workDir)`

## Example

```go
// Dial a server
client, err := sshx.Dial("vagrant", "127.0.0.1:2222")
if err != nil {
  // handle error
}
defer client.Close()

// Run a command
stdout, err := sshx.Run(client, "ls -al")
if err != nil {
  // handle error
}
```

## Install

```sh
go get github.com/matthewmueller/sshx
```

## Development

First, clone the repo:

```sh
git clone https://github.com/matthewmueller/sshx
cd sshx
```

Next, install dependencies:

```sh
go mod tidy
```

Finally, try running the tests:

```sh
go test ./...
```

## License

MIT
