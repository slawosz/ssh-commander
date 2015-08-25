package worker

import (
	"bytes"
	"fmt"
	"os/exec"
)

type CmdWorker struct{}

func NewCmdWorker() *CmdWorker {
	return &CmdWorker{}
}

func (w *CmdWorker) Work(p *WorkerPayload) *HostResult {
	//script := fmt.Sprintf("%v %v %v %v %v %v", p.Script, p.Port, p.User, p.Host.Host, p.Command, p.Password)
	cmd := exec.Command("/bin/sh", p.Script)
	var out bytes.Buffer
	var errbuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errbuf
	err := cmd.Run()
	if err != nil {
		fmt.Println("error:")
		fmt.Println(err)
	}
	//fmt.Println("=====")
	//fmt.Println(out.String())
	//fmt.Println("----")
	//fmt.Println(errbuf.String())
	//fmt.Println("=====")
	//fmt.Println("=====")

	return &HostResult{out.String(), p.Host.Host, p.Port, "0"}
}
