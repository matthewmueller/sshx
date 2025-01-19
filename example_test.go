package sshx_test

import (
	"time"

	"github.com/matthewmueller/sshx"
)

func ExampleDial() {
	// Dial a server
	client, err := sshx.Dial("vagrant", "127.0.0.1:2222")
	if err != nil {
		panic(err)
	}
	defer client.Close()
	err = sshx.Exec(client, "echo 'sshx'")
	if err != nil {
		panic(err)
	}
	// Output:
	// sshx
}

func ExampleDialConfig() {
	cfg := sshx.Configure("vagrant", "127.0.0.1:2222")
	cfg.Timeout = time.Second
	// Dial a server
	client, err := sshx.DialConfig("127.0.0.1:2222", cfg)
	if err != nil {
		panic(err)
	}
	defer client.Close()
	err = sshx.Exec(client, "echo 'sshx'")
	if err != nil {
		panic(err)
	}
	// Output:
	// sshx
}

func ExampleTest() {
	// Dial a server
	signer, err := sshx.Test("vagrant", "127.0.0.1:2222")
	if err != nil {
		panic(err)
	}
	_ = signer.PublicKey()
	// Output:
}
