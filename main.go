package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"

	"github.com/tj/docopt"

	_ "net/http/pprof"

	"golang.org/x/crypto/ssh"
)

var usage = `go-forward-ssh.

  Usage:
    go-forward-ssh (-L | -R) <host> <local-port> <remote-port> [--ssh-key ssh-key] [--ssh-user ssh-user]
    go-forward-ssh -h | --help
    go-forward-ssh --version

  Options:
    --ssh-key ssh-key       The path to your ssh key [default: ~/.ssh/id_rsa]
    --ssh-user ssh-user     The ssh username to use [default: ubuntu]
    -h, --help              output help information
    -v, --version           output version
`

var (
	serverAddrString string
	localAddrString  string
	remoteAddrString string
	args             map[string]interface{}
)

func PublicKeyFile(file string) ssh.AuthMethod {
	buffer, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}

	key, err := ssh.ParsePrivateKey(buffer)
	if err != nil {
		log.Fatal(err)
	}

	return ssh.PublicKeys(key)
}

func ioCopy(dst io.Writer, src io.Reader) {
	if _, err := io.Copy(dst, src); err != nil {
		log.Fatalf("io.Copy failed: %v", err)
	}
}

func remoteForwarding() {
	// Setup localConn (type net.Conn)
	localConn, err := net.Dial("tcp", localAddrString)
	if err != nil {
		log.Fatalf("net.Listen failed: %v", err)
	}

	// Setup sshClientConn (type *ssh.ClientConn)
	sshClientConn, err := ssh.Dial("tcp", serverAddrString, &ssh.ClientConfig{
		User: args["--ssh-user"].(string),
		Auth: []ssh.AuthMethod{PublicKeyFile(args["--ssh-key"].(string))},
	})
	if err != nil {
		log.Fatalf("ssh.Dial failed: %s", err)
	}

	// Setup sshConn (type net.Conn)
	sshListen, err := sshClientConn.Listen("tcp", remoteAddrString)
	if err != nil {
		log.Fatalf("ssh.Listen failed: %s", err)
	}

	for {
		sshConn, err := sshListen.Accept()
		if err != nil {
			log.Fatalf("sshListen.Accept failed: %s", err)
		}

		go ioCopy(sshConn, localConn)
		go ioCopy(localConn, sshConn)
	}
}

func localForwarding() {
	// Setup localConn (type net.Conn)
	localListener, err := net.Listen("tcp", localAddrString)
	if err != nil {
		log.Fatalf("net.Listen failed: %v", err)
	}

	// Setup sshClientConn (type *ssh.ClientConn)
	sshClientConn, err := ssh.Dial("tcp", serverAddrString, &ssh.ClientConfig{
		User: args["--ssh-user"].(string),
		Auth: []ssh.AuthMethod{PublicKeyFile(args["--ssh-key"].(string))},
	})
	if err != nil {
		log.Fatalf("ssh.Dial failed: %s", err)
	}

	// Setup sshConn (type net.Conn)
	sshConn, err := sshClientConn.Dial("tcp", remoteAddrString)
	if err != nil {
		log.Fatalf("ssh.Listen failed: %s", err)
	}

	for {
		localConn, err := localListener.Accept()
		if err != nil {
			log.Fatalf("sshListen.Accept failed: %s", err)
		}

		go ioCopy(sshConn, localConn)
		go ioCopy(localConn, sshConn)
	}
}

func main() {
	var err error
	args, err = docopt.Parse(usage, nil, true, "dev", false)
	if err != nil {
		log.Fatal(err)
	}
	serverAddrString = args["<host>"].(string)
	localAddrString = fmt.Sprintf("localhost:%s", args["<local-port>"].(string))
	remoteAddrString = fmt.Sprintf("localhost:%s", args["<remote-port>"].(string))

	if args["-R"].(bool) {
		remoteForwarding()
	} else if args["-L"].(bool) {
		localForwarding()
	}

}
