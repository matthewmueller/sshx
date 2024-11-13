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
	err = sshx.Exec(client, "ls -al")
	if err != nil {
		panic(err)
	}
}
