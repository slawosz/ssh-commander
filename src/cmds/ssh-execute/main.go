package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"strings"
	"worker"
)

var host = flag.String("host", "localhost", "host to connect")
var port = flag.String("port", "2222", "port to connect")
var user = flag.String("user", "vagrant", "user")
var password = flag.String("password", "vagrant", "password")
var commands = flag.String("commands", "ls", "commands to execute")
var prompt = flag.String("prompt", "$", "prompt sign")
var exit = flag.String("exit", "exit", "exit command")

//var timeout = flag.Int("timeout", 5, "command timeout")

func main() {
	flag.Parse()
	p := &worker.Host{
		User:     *user,
		Password: *password,
		Host:     *host,
		Port:     *port,
		Prompt:   *prompt,
		Exit:     *exit,
		Commands: strings.Split(*commands, ","),
	}
	b, _ := json.Marshal(p)
	fmt.Println(string(b))
	w := worker.NewExpectWorker()
	fmt.Println(w.Work(p).Payload)
}
