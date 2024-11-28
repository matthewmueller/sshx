package sshx

import (
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"net"
	"os"
	osuser "os/user"
	"path/filepath"
	"strings"

	"github.com/skeema/knownhosts"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
	"golang.org/x/term"
)

// Split a user@host[:port] string into user and host.
func Split(userHost string) (user string, host string, err error) {
	parts := strings.Split(userHost, "@")
	if len(parts) != 2 {
		// SSH uses the current user by default
		current, err := osuser.Current()
		if err != nil {
			return "", "", fmt.Errorf("ssh: invalid user@host[:port] %q", userHost)
		}
		return current.Username, userHost, nil
	}
	user, host = parts[0], parts[1]
	parts = strings.Split(host, ":")
	if len(parts) == 1 {
		host += ":22"
	}
	return user, host, nil
}

// Configure creates a new *ClientConfig based on sensible defaults.
// This method is fairly error-resistent and intended for advanced use cases.
func Configure(user, host string, signers ...ssh.Signer) *ssh.ClientConfig {
	config := configure(user, host, signers...)

	// Add the agent auth method if available
	if agent, err := loadAgent(); nil == err {
		config.Auth = append(config.Auth, ssh.PublicKeysCallback(agent.Signers))
	}

	return config
}

func configure(user, host string, signers ...ssh.Signer) *ssh.ClientConfig {
	// Create the client config
	config := &ssh.ClientConfig{
		User: user,
	}

	// Add the known hosts if available
	if knownHosts, err := loadKnownHosts(); nil == err {
		config.HostKeyCallback = knownHosts.HostKeyCallback()
		config.HostKeyAlgorithms = knownHosts.HostKeyAlgorithms(host)
	} else {
		config.HostKeyCallback = ssh.InsecureIgnoreHostKey()
	}

	// Add the signers
	for _, signer := range signers {
		config.Auth = append(config.Auth, ssh.PublicKeys(signer))
	}

	return config
}

// Dial creates a new ssh.Client with sensible defaults
func Dial(user, host string, signers ...ssh.Signer) (*ssh.Client, error) {
	// Configure the ssh client
	config := Configure(user, host, signers...)
	// Dial the ssh connection
	return ssh.Dial("tcp", host, config)
}

// Test the remote host connection, returning the first signer that was
// successfully used to connect to the remote host.
func Test(user, host string, signers ...ssh.Signer) (ssh.Signer, error) {
	// Add the agent signers if available
	if agent, err := loadAgent(); nil == err {
		agentSigners, err := agent.Signers()
		if err != nil {
			return nil, err
		}
		signers = append(signers, agentSigners...)
	}

	// Try each signer until we find one that works
	for _, signer := range signers {
		config := configure(user, host, signer)
		if client, err := ssh.Dial("tcp", host, config); nil == err {
			client.Close()
			return signer, nil
		}
	}

	return nil, errors.New("ssh: no valid signers")
}

// Run a command on the remote host
func Run(ssh *ssh.Client, cmd string) (string, error) {
	session, err := ssh.NewSession()
	if err != nil {
		return "", fmt.Errorf("ssh: could not create session: %w", err)
	}
	defer session.Close()
	stdout := new(bytes.Buffer)
	session.Stdout = stdout
	session.Stderr = os.Stderr
	if err := session.Run(cmd); err != nil {
		return "", err
	}
	// Trim spacing before and after stdout by default
	return strings.TrimRight(stdout.String(), "\n"), nil
}

// Exec a command on the remote host
func Exec(ssh *ssh.Client, cmd string) error {
	session, err := ssh.NewSession()
	if err != nil {
		return fmt.Errorf("ssh: could not create session: %w", err)
	}
	defer session.Close()
	session.Stderr = os.Stderr
	return session.Run(cmd)
}

func Shell(sshc *ssh.Client, dir string, args ...string) error {
	if !fs.ValidPath(dir) {
		return fmt.Errorf("ssh: invalid directory %q", dir)
	}

	session, err := sshc.NewSession()
	if err != nil {
		return fmt.Errorf("ssh: could not create session: %w", err)
	}
	defer session.Close()

	// Change to the specified directory
	session.Stdin = os.Stdin
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr

	// If we have args, don't allocate a terminal. Just run the command and
	// return the result
	if len(args) > 0 {
		return session.Run(formatCommand(dir, args...))
	}

	fd := int(os.Stdin.Fd())

	state, err := term.MakeRaw(fd)
	if err != nil {
		return fmt.Errorf("ssh: could not make terminal raw: %w", err)
	}
	defer term.Restore(fd, state)

	// Get the terminal size
	termWidth, termHeight, err := term.GetSize(fd)
	if err != nil {
		return fmt.Errorf("ssh: could not get terminal size: %w", err)
	}

	// Default to xterm-256color
	termType := os.Getenv("TERM")
	if termType == "" {
		termType = "xterm-256color"
	}

	// request pty
	if err := session.RequestPty(termType, termHeight, termWidth, ssh.TerminalModes{}); err != nil {
		return fmt.Errorf("ssh: could not request pty: %w", err)
	}

	// Wait for the session to complete
	if err := session.Run(formatCommand(dir)); err != nil {
		switch e := err.(type) {
		case *ssh.ExitError:
			// Interrupted by the user (SIGINT)
			if e.ExitStatus() == 130 {
				return nil
			}
			return fmt.Errorf("ssh: exit status %d", e.ExitStatus())
		default:
			return fmt.Errorf("ssh: session ended unexpectedly: %w", err)
		}
	}

	return nil
}

func formatCommand(dir string, args ...string) string {
	if len(args) == 0 {
		return fmt.Sprintf("cd %s && exec $SHELL", dir)
	}
	return fmt.Sprintf("cd %s && exec $SHELL -c %q", dir, strings.Join(args, " "))
}

func loadKnownHosts() (knownhosts.HostKeyCallback, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	knownHostsPath := filepath.Join(home, ".ssh", "known_hosts")
	return knownhosts.New(knownHostsPath)
}

// loadAgent returns an SSH agent client if available.
func loadAgent() (agent.ExtendedAgent, error) {
	unixSocket := os.Getenv("SSH_AUTH_SOCK")
	if unixSocket == "" {
		return nil, errors.New("ssh: SSH_AUTH_SOCK is not set")
	}
	sshAgent, err := net.Dial("unix", unixSocket)
	if err != nil {
		return nil, fmt.Errorf("could not find ssh agent: %w", err)
	}
	return agent.NewClient(sshAgent), nil
}
