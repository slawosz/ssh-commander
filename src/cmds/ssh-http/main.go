package main

import (
	"communication"
	"worker"
)

func main() {

	// create worker and scheduler
	w := worker.NewExpectWorker()
	s := worker.NewScheduler(w, 20)
	http := communication.NewHttp(s)

	go http.Start()

	s.Start()
}
