package communication

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"sync"
	"worker"
)

type Payload struct {
	CallbackUrl string
	*worker.JobPayload
}

type HttpCommunication struct {
	jobs      map[string]string // JID => callback url
	scheduler worker.IScheduler
	*sync.Mutex
}

func NewHttpCommunication(s worker.IScheduler) *HttpCommunication {
	return &HttpCommunication{
		jobs:      make(map[string]string),
		scheduler: s,
	}
}

func (c HttpCommunication) Start(s worker.IScheduler) {
	go func() {
		for res := range c.scheduler.ResultsChan() {
			go c.sendBackResults(res)
		}
	}()

	go func() {
		c.scheduler.Start()
	}()

	go c.serve()

	sleep := make(chan bool)
	sleep <- true
}

func (c *HttpCommunication) serve() {
	handler := func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}

		payload := &Payload{}
		err = json.Unmarshal(body, payload)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}

		c.scheduler.Push(payload.JobPayload)
		c.Lock()
		c.jobs[payload.JobPayload.JID] = payload.CallbackUrl
		c.Unlock()
		w.WriteHeader(http.StatusAccepted)
	}
	http.HandleFunc("/new_job", handler)

	http.ListenAndServe(":8080", nil)
}

func (c *HttpCommunication) sendBackResults(res *worker.HostsResult) {
	body, err := json.Marshal(res)
	bodybuf := bytes.NewBuffer(body)
	if err != nil {
		// log error
		return
	}

	url, ok := c.jobs[res.JID]
	if !ok {
		// log error, not such job
		return
	}
	http.Post(url, "application/json", bodybuf)
	if err != nil {
		// log error
	}
	c.Lock()
	delete(c.jobs, res.JID)
	c.Unlock()
}
