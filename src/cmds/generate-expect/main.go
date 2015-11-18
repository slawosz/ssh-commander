package main

import (
	"flag"
	"fmt"
	"os/exec"
)

var commands = flag.String("commands", "ls", "commands to execute")

func main() {
	flag.Parse()
	cmd := exec.Command("/usr/bin/wc", "-l")
	cmd.Env = []string{
		fmt.Sprintf("PASS=%v", "vagrant"),
	}
	in, err := cmd.StdinPipe()
	//err = cmd.Start()
	if err != nil {
		fmt.Println("error in input pipe: ", err)
	}
	for _, line := range generateScript(*commands) {
		in.Write(line)
	}
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("error in output:")
		fmt.Println(err)
	}
	if err != nil {
		switch err.(type) {
		case *exec.ExitError:
			fmt.Printf("exit error")
		default:
			fmt.Println("error in run:")
			fmt.Println(err)
		}
	}
	in.Close()
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
	fmt.Println(output)
}

func generateScript(cstr string) [][]byte {
	cmds := [][]byte{
		[]byte("set timeout 10"),
		[]byte("spawn ssh -StrictHostKeyChecking=no -p 2222 vagrant@localhost"),
		[]byte("expect \"*?assword: \""),
		[]byte("send -- \"$env(PASS)\r\""),
		[]byte("expect \"*?$ \""),
		[]byte("send -- \"ls\r\""),
		[]byte("expect \"*?$ \""),
	}

	return cmds
}
