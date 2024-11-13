package ssh_test

import (
	"github.com/matthewmueller/ssh"
)

func ExampleDial() {
	// Dial a server
	client, err := ssh.Dial("vagrant@127.0.0.1:2222")
	if err != nil {
		panic(err)
	}
	defer client.Close()
	err = ssh.Exec(client, "ls -al")
	if err != nil {
		panic(err)
	}
}
