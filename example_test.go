package sshx_test

import (
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

func ExampleTest() {
	// Dial a server
	signer, err := sshx.Test("vagrant", "127.0.0.1:2222")
	if err != nil {
		panic(err)
	}
	_ = signer.PublicKey()
	// Output:
}
