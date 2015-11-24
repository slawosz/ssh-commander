package main

import (
	"flag"
	"fmt"
	"worker"
)

var commands = flag.String("commands", "ls", "commands to execute")
var host = flag.String("host", "localhost", "host to connect")
var port = flag.String("port", "2222", "port to connect")
var user = flag.String("user", "vagrant", "user")
var password = flag.String("password", "$vagrant", "password")
var prompt = flag.String("prompt", "$", "prompt sign")

//var timeout = flag.Int("timeout", 5, "command timeout")

func main() {
	flag.Parse()
	p := &worker.WorkerPayload{
		Host: &worker.Host{
			User:     *user,
			Password: *password,
			Host:     *host,
			Port:     *port,
		},
		Commands: *commands,
		Prompt:   *prompt,
	}
	w := worker.NewExpectWorker()
	fmt.Println(w.Work(p).Payload)
}
