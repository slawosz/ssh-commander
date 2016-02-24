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
	chanin  chan *JobPayload
	chanout chan *HostsResult
	pool    int
	worker  Worker
}

func NewScheduler(w Worker, in chan *JobPayload, pool int) *Scheduler {
	return &Scheduler{
		chanin:  in,
		chanout: make(chan *HostsResult),
		worker:  w,
	}
}

func (s *Scheduler) Start(pool int) {
	limiter := make(chan bool, pool)

	for job := range s.chanin {

		go func(j *JobPayload) {
			wg := &sync.WaitGroup{}
			workerResults := make(chan *HostResult, len(j.Hosts))

			for _, hostPayload := range j.Hosts {
				wg.Add(1)
				limiter <- true // blocks when limiter chan buffer is full

				go func(hostPayload *Host) {
					w := &WorkerPayload{Host: hostPayload, Commands: j.Command, Script: j.Script}
					res := s.worker.Work(w)
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
				JID:         j.JID,
				HostsResult: hostsResult,
			}
			s.chanout <- res
		}(job)
	}
}

func (s *Scheduler) PushJob(j *JobPayload) {
	s.chanin <- j
}

func (s *Scheduler) ResultsChan() chan *HostsResult {
	return s.chanout
}
