package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	//	"io"
	"io/ioutil"
	//"os"
	"time"

	// "code.google.com/p/go.crypto/ssh"
	"golang.org/x/crypto/ssh"
)

var cmdsfile = flag.String("cmdsfile", "cmds.json", "json file with commands to execute")
var pubkey = flag.String("pubkey", "", "pub key")
var privkey = flag.String("privkey", "", "priv key")

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

/*
func (t *SignerContainer) Sign(i int, rand io.Reader, data []byte) (sig []byte, err error) {
	if i >= len(t.signers) {
		return
	}
	sig, err = t.signers[i].Sign(rand, data)
	return
}
*/

func main() {
	flag.Parse()
	results := make(chan string, 10)
	timeout := time.After(65 * time.Second)
	config := &ssh.ClientConfig{
		User: "vagrant",
		//Auth: []ssh.AuthMethod{makeKeyring()},
		Auth: []ssh.AuthMethod{ssh.Password("vagrant")},
	}

	cmds := ParseCommands(*cmdsfile)

	for _, cmd := range cmds {
		go func(command *Cmd) {
			results <- executeCmd(command.Command, command.Host, config)
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

func executeCmd(cmd, hostname string, config *ssh.ClientConfig) string {
	// return fmt.Sprintf("executing %v on %v\n", cmd, hostname)
	fmt.Printf("executing %v on %v\n", cmd, hostname)
	conn, err := ssh.Dial("tcp", hostname, config)
	if err != err {
		panic(err)
	}
	session, err := conn.NewSession()
	if err != err {
		panic(err)
	}
	defer session.Close()

	var stdoutBuf bytes.Buffer
	session.Stdout = &stdoutBuf
	err = session.Run(cmd)
	if err != nil {
		panic(err)
	}
	fmt.Printf("DONE: %v on %v\n", cmd, hostname)

	return hostname + ": " + stdoutBuf.String()
}

// for keys
/*
func makeKeyring() []ssh.Signer {
	signers := []ssh.Signer{}
	keys := []string{*privkey, *privkey}

	for _, keyname := range keys {
		signer, err := makeSigner(keyname)
		if err != nil {
			panic(err)
		}
		signers = append(signers, signer)
	}

	return signers
}

func makeSigner(keyname string) (signer ssh.Signer, err error) {
	fp, err := os.Open(keyname)
	if err != nil {
		panic(err)
	}
	defer fp.Close()

	buf, err := ioutil.ReadAll(fp)
	if err != nil {
		panic(err)
	}
	signer, err = ssh.ParsePrivateKey(buf)
	if err != nil {
		fmt.Println(keyname)
		panic(err)
	}
	return
}
*/
