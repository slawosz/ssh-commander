package communication

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"worker"
)

type Payload struct {
	CallbackUrl string
	*worker.JobPayload
}

type HttpInput struct {
	inchan chan *worker.JobPayload
}

func NewHttpInput(in chan *worker.JobPayload) *HttpInput {
	return &HttpInput{in}
}

func (c *HttpInput) StartInput() {
	go c.serve()

	sleep := make(chan bool)
	sleep <- true
}

func (c *HttpInput) serve() {
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

		go func() { c.inchan <- payload.JobPayload }()
		w.WriteHeader(http.StatusAccepted)
	}
	http.HandleFunc("/new_job", handler)

	fmt.Printf("Http server started on %v\n", 7632)
	http.ListenAndServe(":7632", nil)
}
