package sshx_test

import (
	"testing"

	"github.com/matryer/is"
	"github.com/matthewmueller/sshx"
)

func TestSplitNoPort(t *testing.T) {
	is := is.New(t)
	user, host, err := sshx.Split("user@host")
	is.NoErr(err)
	is.Equal(user, "user")
	is.Equal(host, "host:22")
}

func TestSplitWithPort(t *testing.T) {
	is := is.New(t)
	user, host, err := sshx.Split("user@host:1234")
	is.NoErr(err)
	is.Equal(user, "user")
	is.Equal(host, "host:1234")
}
