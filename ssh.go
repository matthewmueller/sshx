package ssh

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"os"
	osuser "os/user"
	"path/filepath"
	"strings"

	"github.com/skeema/knownhosts"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

// Client is an alias for ssh.Client
type Client = ssh.Client
type ClientConfig = ssh.ClientConfig

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
func Configure(user, host string) (*ClientConfig, error) {
	// Create the client config
	config := &ClientConfig{
		User: user,
	}

	// Add the known hosts if available
	if knownHosts, err := loadKnownHosts(); nil == err {
		config.HostKeyCallback = knownHosts.HostKeyCallback()
		config.HostKeyAlgorithms = knownHosts.HostKeyAlgorithms(host)
	} else {
		config.HostKeyCallback = ssh.InsecureIgnoreHostKey()
	}

	// Add the agent auth method if available
	if agentAuth, err := loadAgent(); nil == err {
		config.Auth = []ssh.AuthMethod{agentAuth}
	}

	return config, nil
}

// Dial creates a new ssh.Client with sensible defaults
func Dial(user, host string) (*Client, error) {
	config, err := Configure(user, host)
	if err != nil {
		return nil, err
	}
	// Dial the ssh connection
	return ssh.Dial("tcp", host, config)
}

// Run a command on the remote host
func Run(ssh *Client, cmd string) (string, error) {
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
	return strings.TrimSpace(stdout.String()), nil
}

// Exec a command on the remote host
func Exec(ssh *Client, cmd string) error {
	session, err := ssh.NewSession()
	if err != nil {
		return fmt.Errorf("ssh: could not create session: %w", err)
	}
	defer session.Close()
	session.Stderr = os.Stderr
	return session.Run(cmd)
}

func loadKnownHosts() (knownhosts.HostKeyCallback, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	knownHostsPath := filepath.Join(home, ".ssh", "known_hosts")
	return knownhosts.New(knownHostsPath)
}

func loadAgent() (ssh.AuthMethod, error) {
	unixSocket := os.Getenv("SSH_AUTH_SOCK")
	if unixSocket == "" {
		return nil, errors.New("ssh: missing SSH_AUTH_SOCK")
	}
	sshAgent, err := net.Dial("unix", unixSocket)
	if err != nil {
		return nil, fmt.Errorf("could not find ssh agent: %w", err)
	}
	authMethod := ssh.PublicKeysCallback(agent.NewClient(sshAgent).Signers)
	return authMethod, nil
}
