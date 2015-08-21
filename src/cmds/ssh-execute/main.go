package main

import (
	"flag"
	"fmt"
	"worker"
)

var host = flag.String("host", "localhost", "host to connect")
var port = flag.String("port", "2222", "port to connect")
var user = flag.String("user", "vagrant", "user")
var password = flag.String("password", "vagrant", "password")
var command = flag.String("command", "uname -a", "command to run")
var timeout = flag.Int("timeout", 5, "command timeout")

func main() {
	flag.Parse()

	w := worker.NewWorker()
	payload := prepareWorkerPayload()
	res := w.Work(payload)
	fmt.Printf("Result:\n%+v\n", res)
}

func prepareWorkerPayload() *worker.WorkerPayload {
	return &worker.WorkerPayload{
		Host: &worker.Host{
			Host:     *host,
			Port:     *port,
			User:     *user,
			Password: *password,
		},
		Command: *command,
	}
}
