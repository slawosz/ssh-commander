package main

import (
	"communication"
	"worker"
)

func main() {
	w := worker.NewCmdWorker()
	input := make(chan *worker.JobPayload)
	scheduler := worker.NewScheduler(w, input, 2000)

	in := communication.NewHttpInput(input)
	go func() { in.StartInput() }()
	out := communication.NewStdoutOut(scheduler.ResultsChan())
	go func() { out.StartOut() }()
	scheduler.Start()
}
