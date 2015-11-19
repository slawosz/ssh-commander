package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	exec "utils"
)

var commands = flag.String("commands", "ls", "commands to execute")

func main() {
	flag.Parse()
	cmd := exec.Command("/usr/bin/expect")
	cmd.Env = []string{
		fmt.Sprintf("PASS=%v", "vagrant"),
	}
	in, err := cmd.StdinPipe()
	//err = cmd.Start()
	if err != nil {
		fmt.Println("error in input pipe: ", err)
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println("error in output pipe: ", err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		fmt.Println("error in output pipe: ", err)
	}
	err = cmd.Start()
	if err != nil {
		switch err.(type) {
		case *exec.ExitError:
			fmt.Printf("exit error")
		default:
			fmt.Println("error in run:")
			fmt.Println(err)
		}
	}
	for _, line := range generateScript(*commands) {
		io.Copy(in, bytes.NewBuffer(line))
	}
	in.Close()
	if err != nil {
		fmt.Printf("error in output:")
		fmt.Println(err)
	}
	io.Copy(os.Stdout, stdout)
	io.Copy(os.Stdout, stderr)
	err = cmd.Wait()
	if err != nil {
		switch err.(type) {
		case *exec.ExitError:
			fmt.Println("exit error")
		default:
			fmt.Println("error in run:")
			fmt.Println(err)
		}
	}
}

// TODO: it should be possible to generate expect for debug
func generateScript(cstr string) [][]byte {
	cmdsStr := strings.Split(cstr, ",")
	cmds := [][]byte{
		[]byte("set timeout 10\n"),
		[]byte("spawn ssh -StrictHostKeyChecking=no -p 2222 vagrant@localhost\n"),
		[]byte("expect \"*?assword: \"\n"),
		[]byte("send -- \"$env(PASS)\r\"\n"),
		[]byte("expect \"*?$ \"\n"),
		[]byte("send -- \"ls\r\"\n"),
	}
	for _, cmd := range cmdsStr {
		cmds = append(cmds, []byte(fmt.Sprintf("send -- \"%v\r\"\n", cmd)))
		cmds = append(cmds, []byte("expect \"*?$ \"\n"))
	}
	end := []byte("expect \"*?$ \"\n")
	cmds = append(cmds, end)

	return cmds
}
