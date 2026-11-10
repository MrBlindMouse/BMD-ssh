module github.com/MrBlindMouse/BMD-ssh

go 1.27.3

require github.com/MrBlindMouse/Portico v0.1.2

replace github.com/MrBlindMouse/Portico => ../Portico

require (
	github.com/anmitsu/go-shlex v0.0.0-20200514113438-38f4b401e2be // indirect
	github.com/gliderlabs/ssh v0.3.8 // indirect
	golang.org/x/crypto v0.58.0 // indirect
	golang.org/x/sys v0.50.0 // indirect
)
