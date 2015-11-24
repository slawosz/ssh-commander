package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"strings"

	"golang.org/x/crypto/ssh"
)

func main() {

	config := &ssh.ServerConfig{
		//Define a function to run when a client attempts a password login
		PasswordCallback: func(c ssh.ConnMetadata, pass []byte) (*ssh.Permissions, error) {
			// Should use constant-time compare (or better, salt+hash) in a production setting.
			if c.User() == "vagrant" && string(pass) == "vagrant" {
				return nil, nil
			}
			return nil, fmt.Errorf("password rejected for %q", c.User())
		},
		// NoClientAuth: true,
	}

	// You can generate a keypair with 'ssh-keygen -t rsa'
	privateBytes, err := ioutil.ReadFile("gossh_rsa")
	if err != nil {
		log.Fatal("Failed to load private key (./gossh_rsa)")
	}

	private, err := ssh.ParsePrivateKey(privateBytes)
	if err != nil {
		log.Fatal("Failed to parse private key")
	}

	config.AddHostKey(private)

	// Once a ServerConfig has been configured, connections can be accepted.
	socket, err := net.Listen("tcp", "0.0.0.0:2200")
	if err != nil {
		log.Fatalf("Failed to listen on 2200 (%s)", err)
	}

	for {
		conn, err := socket.Accept()
		if err != nil {
			panic(err)
		}

		// From a standard TCP connection to an encrypted SSH connection
		sshConn, newChans, _, err := ssh.NewServerConn(conn, config)
		if err != nil {
			panic(err)
		}
		defer sshConn.Close()

		log.Println("Connection from", sshConn.RemoteAddr())
		go func() {
			for chanReq := range newChans {
				go handleChanReq(chanReq)
			}
		}()
	}

}

func handleChanReq(chanReq ssh.NewChannel) {
	if chanReq.ChannelType() != "session" {
		chanReq.Reject(ssh.Prohibited, "channel type is not a session")
		return
	}

	ch, reqs, err := chanReq.Accept()
	if err != nil {
		log.Println("fail to accept channel request", err)
		return
	}

	req := <-reqs
	if req.Type != "exec" {
		ch.Write([]byte("request type '" + req.Type + "' is not 'exec'\r\n"))
		ch.Close()
		return
	}

	handleExec(ch, req)
}

// Payload: int: command size, string: command
func handleExec(ch ssh.Channel, req *ssh.Request) {
	command := string(req.Payload[4:])
	gitCmds := []string{"git-receive-pack", "git-upload-pack"}

	valid := false
	for _, cmd := range gitCmds {
		if strings.HasPrefix(command, cmd) {
			valid = true
		}
	}
	req.Reply(true, []byte("0\r\n"))
	if !valid {
		ch.Write([]byte("command is not a GIT command\r\n"))

		ch.Close()
		return
	}

	ch.Write([]byte("well done!\r\n"))
	ch.Close()
}
