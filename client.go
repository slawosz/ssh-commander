package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"time"

	"golang.org/x/crypto/ssh"
)

var cmdsfile = flag.String("cmdsfile", "cmds.json", "json file with commands to execute")

type Cmd struct {
	Host    string
	Command string
}

func ParseCommands(filename string) []*Cmd {
	jsonBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	cmds := make([]*Cmd, 0)
	err = json.Unmarshal(jsonBytes, &cmds)
	if err != nil {
		panic(err)
	}
	return cmds
}

type SignerContainer struct {
	signers []ssh.Signer
}

func (t *SignerContainer) Key(i int) (key ssh.PublicKey, err error) {
	if i >= len(t.signers) {
		return
	}
	key = t.signers[i].PublicKey()
	return
}

func main() {
	flag.Parse()
	results := make(chan string, 10)
	timeout := time.After(65 * time.Second)
	cmds := ParseCommands(*cmdsfile)

	for _, cmd := range cmds {
		go func(command *Cmd) {
			results <- executeCmd(command.Command, command.Host)
		}(cmd)
	}

	for i := 0; i < len(cmds); i++ {
		select {
		case res := <-results:
			fmt.Print(res)
		case <-timeout:
			fmt.Println("Timed out!")
			return
		}
	}
}

func executeCmd(cmd, hostname string) string {
	// return fmt.Sprintf("executing %v on %v\n", cmd, hostname)
	config := &ssh.ClientConfig{
		User:   "bob",
		Auth:   []ssh.AuthMethod{ssh.Password("bob2")},
		Config: ssh.Config{Ciphers: []string{"aes192-ctr"}},
	}

	fmt.Printf("executing %v on %v\n", cmd, hostname)
	conn, err := ssh.Dial("tcp", hostname, config)
	if err != nil {
		fmt.Println(err)
	}
	if conn == nil {
		fmt.Printf("Connection to %v failed\n", hostname)
		return ""
	}
	session, err := conn.NewSession()
	if err != err {
		panic(err)
	}
	defer session.Close()

	var stdoutBuf bytes.Buffer
	session.Stdout = &stdoutBuf
	err = session.Start(cmd)
	if err != nil {
		fmt.Printf("Session failed with error: %v\n", err)
	}
	session.Wait()
	fmt.Printf("DONE: %v on %v\n", cmd, hostname)

	return hostname + ": " + stdoutBuf.String()
}
