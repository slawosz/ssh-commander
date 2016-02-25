package communication

import (
	"encoding/json"
	"fmt"
	"github.com/elazarl/go-bindata-assetfs"
	"io/ioutil"
	"net/http"
	"worker"
)

type Http struct {
	scheduler worker.JobScheduler
}

func NewHttp(scheduler worker.JobScheduler) *Http {
	return &Http{
		scheduler: scheduler,
	}
}

func (h *Http) Start() {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "accept, content-type")
		w.Header().Set("Content-Type", "application/json")
		if r.Method != "POST" {
			w.WriteHeader(http.StatusOK)
			return
		}

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var payload []*worker.Host
		err = json.Unmarshal(body, &payload)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		resCh := make(chan []*worker.HostResult)
		h.scheduler.PushJob(&worker.SchedulerPayload{Hosts: payload, ResultsChan: resCh})
		res := <-resCh
		b, err := json.Marshal(res)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		fmt.Fprintf(w, "%v", string(b))

		w.WriteHeader(http.StatusOK)
	}
	http.HandleFunc("/run", handler)
	http.Handle("/",
		http.FileServer(
			&assetfs.AssetFS{Asset: Asset, AssetDir: AssetDir, AssetInfo: AssetInfo, Prefix: "src/assets"}))

	fmt.Printf("Http server started on %v\n", 7632)
	http.ListenAndServe(":3002", nil)
}
