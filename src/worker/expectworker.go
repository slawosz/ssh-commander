package worker

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"strings"
	exec "utils"
)

type ExpectWorker struct{}

func NewExpectWorker() *ExpectWorker {
	return &ExpectWorker{}
}

func (w *ExpectWorker) Work(p *WorkerPayload) *HostResult {
	flag.Parse()
	cmd := exec.Command("/usr/bin/expect")
	in, err := cmd.StdinPipe()
	output := bytes.NewBuffer(make([]byte, 0))
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
	for _, line := range generateScript(p) {
		io.Copy(in, bytes.NewBuffer(line))
	}
	in.Close()
	if err != nil {
		fmt.Printf("error in output:")
		fmt.Println(err)
	}
	io.Copy(output, stdout)
	io.Copy(output, stderr)
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

	return &HostResult{Payload: output.String(), Port: p.Port, Host: p.Host.Host}
}

// TODO: it should be possible to generate expect for debug
func generateScript(p *WorkerPayload) [][]byte {
	cmdsStr := strings.Split(p.Commands, ",")
	cmds := [][]byte{
		[]byte("set timeout 10\n"),
		[]byte(fmt.Sprintf("spawn ssh -StrictHostKeyChecking=no -p %v %v@%v\n", p.Port, p.User, p.Host.Host)),
		[]byte("expect \"*?assword: \"\n"),
		// escape any $ in password for expect
		[]byte(fmt.Sprintf("send -- \"%v\r\"\n", strings.Replace(p.Password, "$", "\\$", -1))),
		[]byte(fmt.Sprintf("expect \"*?%v \"\n", p.Prompt)),
		[]byte("send -- \"ls\r\"\n"),
	}
	for _, cmd := range cmdsStr {
		cmds = append(cmds, []byte(fmt.Sprintf("send -- \"%v\r\"\n", cmd)))
		cmds = append(cmds, []byte("expect \"*?$ \"\n"))
	}
	end := []byte(fmt.Sprintf("expect \"*?%v \"\n", p.Prompt))
	cmds = append(cmds, end)

	return cmds
}
