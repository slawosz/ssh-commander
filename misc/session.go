package main

import (
	//"bytes"
	"fmt"
	//"io/ioutil"
	"log"
	"time"

	"golang.org/x/crypto/ssh"
)

func main() {
	config := &ssh.ClientConfig{
		User: "vagrant",
		Auth: []ssh.AuthMethod{
			ssh.Password("vagrant"),
		},
	}
	// Connect to ssh server
	conn, err := ssh.Dial("tcp", "localhost:2222", config)
	if err != nil {
		log.Fatalf("unable to connect: %s", err)
	}
	defer conn.Close()
	// Create a session
	session, err := conn.NewSession()
	if err != nil {
		log.Fatalf("unable to create session: %s", err)
	}
	defer session.Close()
	// Set up terminal modes
	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}
	// Request pseudo terminal
	if err := session.RequestPty("xterm", 80, 40, modes); err != nil {
		log.Fatalf("request for pseudo terminal failed: %s", err)
	}
	in, err := session.StdinPipe()
	if err != nil {
		panic(err)
	}
	o, _ := session.StdoutPipe()
	// Start remote shell
	if err := session.Shell(); err != nil {
		log.Fatalf("failed to start shell: %s", err)
	}
	in.Write([]byte("ls /\r\n"))
	in.Write([]byte("ls / | wc -l\r\n"))
	fmt.Println("writen")
	time.Sleep(1 * time.Second)
	fmt.Println("before read")
	buf := make([]byte, 10000)
	_, _ = o.Read(buf)
	fmt.Println("read")
	fmt.Println(string(buf))
	if err != nil {
		panic(err)
	}
}
