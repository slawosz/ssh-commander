package worker

import (
	"sync"
)

type JobScheduler interface {
	Start(pool int)
	PushJob(*JobPayload)
	ResultsChan() chan *HostsResult
}

type Scheduler struct {
	chanin  chan []*Host
	chanout chan *HostsResult
	worker  Worker
	pool    int
}

func NewScheduler(w Worker, in chan []*Host, pool int) *Scheduler {
	return &Scheduler{
		chanin:  in,
		chanout: make(chan *HostsResult),
		worker:  w,
		pool:    pool,
	}
}

func (s *Scheduler) Start() {
	limiter := make(chan bool, s.pool)

	for job := range s.chanin {

		go func(j []*Host) {
			// waitGroup is to wait until all hosts will do their job
			wg := &sync.WaitGroup{}
			workerResults := make(chan *HostResult, len(j))

			for _, hostPayload := range j {
				wg.Add(1)
				limiter <- true // blocks when limiter chan buffer is full

				go func(hostPayload *Host) {
					res := s.worker.Work(hostPayload)
					workerResults <- res

					<-limiter //release place in channel
					wg.Done()
				}(hostPayload)
			}
			wg.Wait()

			close(workerResults) // TODO: test if without it reduce blocks
			hostsResult := make([]*HostResult, 0)

			for r := range workerResults {
				hostsResult = append(hostsResult, r)
			}

			res := &HostsResult{
				HostsResult: hostsResult,
			}
			s.chanout <- res
		}(job)
	}
}

func (s *Scheduler) PushJob(j []*Host) {
	s.chanin <- j
}

func (s *Scheduler) ResultsChan() chan *HostsResult {
	return s.chanout
}
