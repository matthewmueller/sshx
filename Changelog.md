# 0.0.12 / 2025-01-19

- add a DialConfig function

# 0.0.11 / 2024-12-28

- support connecting to unknown hosts and updating known_hosts by default

# 0.0.10 / 2024-12-24

- sshx.Exec should also pipe stdout to os.Stdout

# 0.0.9 / 2024-12-10

- imply host:22 when host doesn't include port

# 0.0.8 / 2024-11-29

- add a `DialEach` function that also returns the signer
- improved error message for `Test`

# 0.0.7 / 2024-11-28

- add `Test` function
- Allow `Dial` and `Configure` to take signers

# 0.0.6 / 2024-11-24

- support allocating an interactive shell

# 0.0.5 / 2024-11-13

- remove aliases

# 0.0.4 / 2024-11-13

- ssh -> sshx to make it easier to pull functions from underlying ssh

# 0.0.3 / 2024-11-13

- switch to splitting user and host as separate params

# 0.0.2 / 2024-11-12

- trim spaces

# 0.0.1 / 2024-11-12

- Initial commit
