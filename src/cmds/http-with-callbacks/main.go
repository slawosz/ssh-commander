package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// JID => url
var results map[string]string
var m *sync.Mutex

var scheduler *worker.Scheduler

type Payload struct {
	CallbackUrl string
	worker.JobPayload
}

func main() {
	w := worker.NewWorker()
	scheduler = worker.NewScheduler(200)

	go func() {
		for res := range scheduler.ResultsChan {
			go sendBackResults(res)
		}
	}()

	go func() {
		scheduler.Start()
	}()

	go startHttp()

	sleep := make(chan bool)
	sleep <- true
}

func startHttp() {
	handler := func(w http.ResponseWriter, r *http.Request) {
		body := []byte{}
		job := &JobPayload{}
		err := json.Unmarshal(job, body)
		if err != nil {
			// log
		}
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	}
	http.HandleFunc("/new_job", handler)

	http.ListenAndServe(":8080", nil)
}

func sendBackResults(res *worker.JobResult) {
	body, err := json.Marshal(res)
	if err != nil {
		// log error
	}
	sendRequest(results[res.JID], body)
}
