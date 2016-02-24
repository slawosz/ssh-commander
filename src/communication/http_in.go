package communication

/*
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
	inchan chan []*worker.Host
}

func NewHttpInput(in chan []*worker.Host) *HttpInput {
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

		var payload []*worker.Host
		err = json.Unmarshal(body, &payload)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		go func() { c.inchan <- payload }()
		w.WriteHeader(http.StatusAccepted)
	}
	http.HandleFunc("/run", handler)

	fmt.Printf("Http server started on %v\n", 7632)
	http.ListenAndServe(":7632", nil)
}
*/
