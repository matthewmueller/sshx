# SSH

[![Go Reference](https://pkg.go.dev/badge/github.com/matthewmueller/ssh.svg)](https://pkg.go.dev/github.com/matthewmueller/ssh)

Tiny wrapper around [golang.org/x/crypto/ssh](https://pkg.go.dev/golang.org/x/crypto/ssh) that has better defaults and a slightly higher-level API.

The goal being it works exactly like if you did `ssh user@host` on your machine.

## Features

- Handles `~/.ssh/known_hosts` on OSX thanks to [skeema/knownhosts](github.com/skeema/knownhosts).
- Uses the active SSH agent on your machine if there is one, allowing you to seamlessly connect without providing a private key (and often the password needed to decrypt that private key).
- Returns an `*golang.org/x/crypto/ssh.Client` that can be used elsewhere
- Adds `Run(ssh, cmd) (stdout, error)` and `Exec(ssh, cmd) error` commands.

## Example

```go
// Dial a server
client, err := ssh.Dial("vagrant@127.0.0.1:2222")
if err != nil {
  // handle error
}
defer client.Close()

// Run a command
stdout, err := ssh.Run(client, "ls -al")
if err != nil {
  // handle error
}
```

## Install

```sh
go get github.com/matthewmueller/ssh
```

## Development

First, clone the repo:

```sh
git clone https://github.com/matthewmueller/ssh
cd ssh
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
