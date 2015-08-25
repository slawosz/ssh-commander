package communication

import (
	"fmt"
	"time"
	"worker"
)

type StdoutOut struct {
	outchan chan *worker.HostsResult
}

func NewStdoutOut(out chan *worker.HostsResult) *StdoutOut {
	return &StdoutOut{out}
}

func (s *StdoutOut) StartOut() {
	start := time.Now().Unix()
	for r := range s.outchan {
		fmt.Printf("======= Result for %v after %vs ==========\n", r.JID, time.Now().Unix()-start)
		for _, h := range r.HostsResult {
			fmt.Printf("  Result on %v after %vs\n", h.Host, time.Now().Unix()-start)
			fmt.Printf("  %v\n", h.Payload)
		}
		fmt.Printf("----\n")
	}
}
