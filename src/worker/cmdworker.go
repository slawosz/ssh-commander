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

func (w *CmdWorker) Work(payload *WorkerPayload) *HostResult {
	cmd := exec.Command("/bin/sh", "ssh2.sh")
	fmt.Printf("%+v\n", cmd)
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

	return &HostResult{out.String(), payload.Host.Host, payload.Port, "0"}
}
