package worker

import (
	"fmt"
	"os/exec"
)

type CmdWorker struct{}

func NewCmdWorker() *CmdWorker {
	return &CmdWorker{}
}

func (w *CmdWorker) Work(p *WorkerPayload) *HostResult {
	cmd := exec.Command("/bin/sh", p.Script)
	cmd.Env = []string{
		fmt.Sprintf("PORT=%v", p.Port),
		fmt.Sprintf("USER=%v", p.User),
		fmt.Sprintf("HOST=%v", p.Host.Host),
		fmt.Sprintf("COMMAND=%v", p.Command),
		fmt.Sprintf("PASSWORD=%v", p.Password),
	}
	//var out bytes.Buffer
	//var errbuf bytes.Buffer
	//cmd.Stdout = &out
	//cmd.Stderr = &errbuf
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("error in output:")
		fmt.Println(err)
	}
	err = cmd.Run()
	if err != nil {
		switch err.(type) {
		case *exec.ExitError:
			//fmt.Println("exit error")
		default:
			//fmt.Println("error in run:")
			//fmt.Println(err)
		}
	}
	//fmt.Println("=====")
	//fmt.Println(out.String())
	//fmt.Println("----")
	//fmt.Println(errbuf.String())
	//fmt.Println("=====")
	//fmt.Println("=====")

	return &HostResult{string(output), p.Host.Host, p.Port, "0"}
}
